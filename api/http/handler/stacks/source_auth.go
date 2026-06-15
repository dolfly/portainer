package stacks

import (
	"fmt"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"
)

// validateSourceForStack checks that the given Source exists and is a git Source, and returns it.
// TODO(BE-12905): enforce per-user access policies once Source ownership is introduced.
func validateSourceForStack(tx dataservices.DataStoreTx, sourceID portainer.SourceID) (*portainer.Source, *httperror.HandlerError) {
	src, err := tx.Source().Read(sourceID)
	if err != nil {
		if tx.IsErrObjectNotFound(err) {
			return nil, httperror.NotFound("Source not found", err)
		}
		return nil, httperror.InternalServerError("Unable to read source", err)
	}

	if src.Type != portainer.SourceTypeGit {
		return nil, httperror.BadRequest(fmt.Sprintf("source %d is not a git source", sourceID), nil)
	}

	return src, nil
}
