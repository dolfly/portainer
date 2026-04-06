package stackbuilders

import (
	"context"
	"errors"
	"time"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	httperrors "github.com/portainer/portainer/api/http/errors"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"
	"github.com/portainer/portainer/pkg/libhttp/request"

	"github.com/rs/zerolog/log"
)

type PostDeployHook func(context.Context, *portainer.Stack) error

type StackBuilderDirector struct {
	builder   any
	dataStore dataservices.DataStore
}

func NewStackBuilderDirector(dataStore dataservices.DataStore, b any) *StackBuilderDirector {
	return &StackBuilderDirector{
		builder:   b,
		dataStore: dataStore,
	}
}

// Build executes the stack build process based on the builder type. It returns the
// created stack and any error encountered during the process.
// The returned error is of type *httperror.HandlerError, which could be a BadRequest
// or InternalServerError depending on the error encountered during the stack build process.
//
// For all stack types, the stack is saved to DB with Status=Deploying and returned
// immediately. Deployment runs in a background goroutine. The caller must poll
// GET /stacks/{id} to track completion.
func (d *StackBuilderDirector) Build(ctx context.Context, payload *StackPayload, endpoint *portainer.Endpoint) (*portainer.Stack, *httperror.HandlerError) {
	// To align with the flow of the actual service deployment tools, we save
	// the stack before the deployment. This allows us to track the stack
	// metadata and partially created resources.
	switch builder := d.builder.(type) {
	case GitMethodStackBuildProcess:
		stack, err := builder.SetGeneralInfo(payload, endpoint).
			SetUniqueInfo(payload).
			SetGitRepository(ctx, payload).
			SaveStack()
		if err != nil {
			if errors.Is(err, httperrors.ErrUnauthorized) {
				return nil, httperror.Forbidden("User not authorized to use git credential", err)
			}
			return nil, httperror.InternalServerError("Failed to save stack via Git repository method", err)
		}

		d.spawnAsyncDeployment(ctx, stack.ID, func() error {
			builder.Deploy(ctx, payload, endpoint)
			return builder.Error()
		}, builder.EnableAutoUpdate)

		return stack, nil

	case FileUploadMethodStackBuildProcess:
		stack, err := builder.SetGeneralInfo(payload, endpoint).
			SetUniqueInfo(payload).
			SetUploadedFile(payload).
			SaveStack()
		if err != nil {
			return nil, httperror.InternalServerError("Failed to save stack via File Upload method", err)
		}

		d.spawnAsyncDeployment(ctx, stack.ID, func() error {
			builder.Deploy(ctx, payload, endpoint)
			return builder.Error()
		})

		return stack, nil

	case FileContentMethodStackBuildProcess:
		stack, err := builder.SetGeneralInfo(payload, endpoint).
			SetUniqueInfo(payload).
			SetFileContent(payload).
			SaveStack()
		if err != nil {
			return nil, httperror.InternalServerError("Failed to save stack via File Content method", err)
		}

		d.spawnAsyncDeployment(ctx, stack.ID, func() error {
			builder.Deploy(ctx, payload, endpoint)
			return builder.Error()
		})

		return stack, nil

	case UrlMethodStackBuildProcess:
		stack, err := builder.SetGeneralInfo(payload, endpoint).
			SetUniqueInfo(payload).
			SetURL(payload).
			SaveStack()
		if err != nil {
			return nil, httperror.InternalServerError("Failed to save stack via URL method", err)
		}

		d.spawnAsyncDeployment(ctx, stack.ID, func() error {
			builder.Deploy(ctx, payload, endpoint)
			return builder.Error()
		})

		return stack, nil

	default:
		return nil, httperror.BadRequest("Invalid value for query parameter: method. Value must be one of: string or repository or url or file", errors.New(request.ErrInvalidQueryParameter))
	}
}

// spawnAsyncDeployment runs the provided deploy function in a background goroutine
// and updates the stack status in the database upon completion.
func (d *StackBuilderDirector) spawnAsyncDeployment(ctx context.Context, stackID portainer.StackID, deploy func() error, hooks ...PostDeployHook) {
	go func() {
		deployErr := deploy()

		if err := d.dataStore.UpdateTx(func(tx dataservices.DataStoreTx) error {
			stack, err := tx.Stack().Read(stackID)
			if err != nil {
				return err
			}

			if deployErr != nil {
				stack.Status = portainer.StackStatusError
				stack.DeploymentStatus = append(stack.DeploymentStatus, portainer.StackDeploymentStatus{
					Status:  portainer.StackStatusError,
					Time:    time.Now().Unix(),
					Message: deployErr.Error(),
				})
			} else {
				stack.Status = portainer.StackStatusActive
				stack.DeploymentStatus = append(stack.DeploymentStatus, portainer.StackDeploymentStatus{
					Status: portainer.StackStatusActive,
					Time:   time.Now().Unix(),
				})

				for _, hook := range hooks {
					if err := hook(ctx, stack); err != nil {
						log.Error().Err(err).
							Int("stack_id", int(stackID)).
							Str("context", "StackBuilderDirector.spawnAsyncDeployment").
							Msg("Failed to run post-deployment hook")
					}
				}
			}

			return tx.Stack().Update(stack.ID, stack)
		}); err != nil {
			log.Error().Err(err).
				Int("stack_id", int(stackID)).
				Str("context", "StackBuilderDirector.spawnAsyncDeployment").
				Msg("Failed to update stack status after async deployment")
		}
	}()
}
