package stacks

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	gittypes "github.com/portainer/portainer/api/git/types"
)

// loadGitConfigFromSource reads the GitConfig and SourceID from the first Git-type Source in the workflow.
func loadGitConfigFromSource(tx dataservices.DataStoreTx, workflowID portainer.WorkflowID) (*gittypes.RepoConfig, portainer.SourceID, error) {
	src, err := dataservices.GitSourceForWorkflow(tx, workflowID)
	if err != nil || src == nil {
		return nil, 0, err
	}

	return src.GitConfig, src.ID, nil
}

// saveSourceGitConfig updates the GitConfig on the Source record identified by sourceID within a transaction.
func saveSourceGitConfig(tx dataservices.DataStoreTx, sourceID portainer.SourceID, gitConfig *gittypes.RepoConfig) error {
	src, err := tx.Source().Read(sourceID)
	if err != nil {
		return fmt.Errorf("failed to read source: %w", err)
	}

	src.GitConfig = gitConfig

	return tx.Source().Update(src.ID, src)
}

// fillStackGitConfig loads GitConfig from Source and sets it on the stack for backwards-compatible responses.
func fillStackGitConfig(tx dataservices.DataStoreTx, stack *portainer.Stack) error {
	if stack.WorkflowID == 0 {
		return nil
	}

	gitConfig, _, err := loadGitConfigFromSource(tx, stack.WorkflowID)
	if err != nil {
		return err
	}

	stack.GitConfig = gittypes.SanitizeRepoConfig(gitConfig)

	return nil
}
