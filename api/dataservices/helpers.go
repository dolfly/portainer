package dataservices

import (
	"errors"
	"fmt"

	portainer "github.com/portainer/portainer/api"
	perrors "github.com/portainer/portainer/api/dataservices/errors"
	gittypes "github.com/portainer/portainer/api/git/types"

	"github.com/rs/zerolog/log"
)

// ErrStop signals the stop of computation when filtering results
var ErrStop = errors.New("stop")

// gitSourceReader is a minimal interface satisfied by both CE and EE DataStoreTx,
// used to avoid duplicating GitSourceForWorkflow across packages.
type gitSourceReader interface {
	Workflow() WorkflowService
	Source() SourceService
}

// GitSourceForWorkflow returns the first git-type Source for the given workflow.
// Returns nil, nil when workflowID is 0 or no git source is found.
func GitSourceForWorkflow(tx gitSourceReader, workflowID portainer.WorkflowID) (*portainer.Source, error) {
	if workflowID == 0 {
		return nil, nil
	}

	wf, err := tx.Workflow().Read(workflowID)
	if err != nil {
		return nil, err
	}

	for _, srcID := range wf.SourceIDs {
		src, err := tx.Source().Read(srcID)
		if err != nil {
			return nil, err
		}

		if src.Type == portainer.SourceTypeGit {
			return src, nil
		}
	}

	return nil, nil
}

// GitConfigHashForWorkflow returns the git commit hash for the first git-type Source of the given workflow.
// Returns "" when workflowID is 0, no git source is found, or on error.
func GitConfigHashForWorkflow(tx gitSourceReader, workflowID portainer.WorkflowID) string {
	src, err := GitSourceForWorkflow(tx, workflowID)
	if err != nil || src == nil || src.GitConfig == nil {
		return ""
	}

	return src.GitConfig.ConfigHash
}

// FindOrCreateGitSource returns an existing Source whose URL and authentication match src,
// or creates a new one. ConfigHash is excluded from matching as it is transient runtime state.
func FindOrCreateGitSource(tx gitSourceReader, src *portainer.Source) (*portainer.Source, error) {
	existing, err := tx.Source().ReadAll(func(s portainer.Source) bool {
		return s.Type == portainer.SourceTypeGit &&
			s.GitConfig != nil &&
			s.GitConfig.URL == src.GitConfig.URL &&
			gitAuthMatches(s.GitConfig.Authentication, src.GitConfig.Authentication)
	})
	if err != nil {
		return nil, err
	}

	if len(existing) > 0 {
		return &existing[0], nil
	}

	if err := tx.Source().Create(src); err != nil {
		return nil, err
	}

	return src, nil
}

func gitAuthMatches(a, b *gittypes.GitAuthentication) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a.Username == b.Username && a.Password == b.Password && a.GitCredentialID == b.GitCredentialID
}

func IsErrObjectNotFound(e error) bool {
	return errors.Is(e, perrors.ErrObjectNotFound)
}

// AppendFn appends elements to the given collection slice
func AppendFn[T any](collection *[]T) func(obj any) (any, error) {
	return func(obj any) (any, error) {
		element, ok := obj.(*T)
		if !ok {
			log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("type assertion failed")
			return nil, fmt.Errorf("failed to convert to %T object: %#v", new(T), obj)
		}

		*collection = append(*collection, *element)

		return new(T), nil
	}
}

// FilterFn appends elements to the given collection when the predicate is true
func FilterFn[T any](collection *[]T, predicate func(T) bool) func(obj any) (any, error) {
	return func(obj any) (any, error) {
		element, ok := obj.(*T)
		if !ok {
			log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("type assertion failed")
			return nil, fmt.Errorf("failed to convert to %T object: %#v", new(T), obj)
		}

		if predicate(*element) {
			*collection = append(*collection, *element)
		}

		return new(T), nil
	}
}

// FirstFn sets the element to the first one that satisfies the predicate and stops the computation, returns ErrStop on
// success
func FirstFn[T any](element *T, predicate func(T) bool) func(obj any) (any, error) {
	return func(obj any) (any, error) {
		e, ok := obj.(*T)
		if !ok {
			log.Debug().Str("obj", fmt.Sprintf("%#v", obj)).Msg("type assertion failed")
			return nil, fmt.Errorf("failed to convert to %T object: %#v", new(T), obj)
		}

		if predicate(*e) {
			*element = *e
			return new(T), ErrStop
		}

		return new(T), nil
	}
}
