package stackbuilders

import (
	"context"
	"fmt"
	"strconv"
	"time"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/filesystem"
	gittypes "github.com/portainer/portainer/api/git/types"
	"github.com/portainer/portainer/api/scheduler"
	"github.com/portainer/portainer/api/stacks/deployments"
	"github.com/portainer/portainer/api/stacks/stackutils"
)

type GitMethodStackBuildProcess interface {
	// Set general stack information
	SetGeneralInfo(payload *StackPayload, endpoint *portainer.Endpoint) GitMethodStackBuildProcess
	// Set unique stack information, e.g. swarm stack has swarmID, kubernetes stack has namespace
	SetUniqueInfo(payload *StackPayload) GitMethodStackBuildProcess
	// Deploy stack based on the configuration
	Deploy(ctx context.Context, payload *StackPayload, endpoint *portainer.Endpoint) GitMethodStackBuildProcess
	// Save the stack information to database
	SaveStack() (*portainer.Stack, error)
	// Get response from HTTP request. Use if it is needed
	GetResponse() string
	// Set git repository configuration
	SetGitRepository(ctx context.Context, payload *StackPayload) GitMethodStackBuildProcess
	Error() error
	EnableAutoUpdate(ctx context.Context, stack *portainer.Stack) error
}

type GitMethodStackBuilder struct {
	StackBuilder
	gitService portainer.GitService
	scheduler  *scheduler.Scheduler
}

func (b *GitMethodStackBuilder) SetGeneralInfo(payload *StackPayload, endpoint *portainer.Endpoint) GitMethodStackBuildProcess {
	stackID := b.dataStore.Stack().GetNextIdentifier()
	b.stack.ID = portainer.StackID(stackID)
	b.stack.EndpointID = endpoint.ID
	b.stack.AdditionalFiles = payload.AdditionalFiles
	now := time.Now().Unix()
	b.stack.Status = portainer.StackStatusDeploying
	b.stack.CreationDate = now
	b.stack.DeploymentStatus = []portainer.StackDeploymentStatus{
		{Status: portainer.StackStatusDeploying, Time: now},
	}
	b.stack.AutoUpdate = payload.AutoUpdate

	return b
}

func (b *GitMethodStackBuilder) SetUniqueInfo(payload *StackPayload) GitMethodStackBuildProcess {
	b.stack.AutoUpdate = payload.AutoUpdate
	return b
}

func (b *GitMethodStackBuilder) SetGitRepository(ctx context.Context, payload *StackPayload) GitMethodStackBuildProcess {
	if b.hasError() {
		return b
	}

	var repoConfig gittypes.RepoConfig
	if payload.Authentication {
		repoConfig.Authentication = &gittypes.GitAuthentication{
			Username: payload.Username,
			Password: payload.Password,
		}
	}

	repoConfig.URL = payload.URL
	repoConfig.ReferenceName = payload.ReferenceName
	repoConfig.TLSSkipVerify = payload.TLSSkipVerify

	repoConfig.ConfigFilePath = payload.ComposeFile
	if payload.ComposeFile == "" {
		repoConfig.ConfigFilePath = filesystem.ComposeFileDefaultName
	}

	// If a manifest file is specified (for kube git apps), then use it instead of the default compose file name
	if payload.ManifestFile != "" {
		repoConfig.ConfigFilePath = payload.ManifestFile
	}

	stackFolder := strconv.Itoa(int(b.stack.ID))
	// Set the project path on the disk
	b.stack.ProjectPath = b.fileService.GetStackProjectPath(stackFolder)

	getProjectPath := func() string {
		stackFolder := fmt.Sprintf("%d", b.stack.ID)
		return b.fileService.GetStackProjectPath(stackFolder)
	}

	commitHash, err := stackutils.DownloadGitRepository(ctx, repoConfig, b.gitService, getProjectPath)
	if err != nil {
		b.err = fmt.Errorf("failed to download git repository: %w", err)
		return b
	}

	// Update the latest commit id
	repoConfig.ConfigHash = commitHash
	b.stack.GitConfig = &repoConfig

	return b
}

func (b *GitMethodStackBuilder) Deploy(ctx context.Context, payload *StackPayload, endpoint *portainer.Endpoint) GitMethodStackBuildProcess {
	if b.hasError() {
		return b
	}

	// Deploy the stack
	b.err = b.deploymentConfiger.Deploy(ctx)

	return b
}

func (b *GitMethodStackBuilder) GetResponse() string {
	return ""
}

func (b *GitMethodStackBuilder) EnableAutoUpdate(ctx context.Context, stack *portainer.Stack) error {
	if stack.AutoUpdate == nil || stack.AutoUpdate.Interval == "" {
		return nil
	}

	jobID, err := deployments.StartAutoupdate(ctx, stack.ID,
		stack.AutoUpdate.Interval,
		b.scheduler,
		b.stackDeployer,
		b.dataStore,
		b.gitService)
	if err != nil {
		return err
	}

	stack.AutoUpdate.JobID = jobID
	return nil
}
