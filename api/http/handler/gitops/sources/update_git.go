package sources

import (
	"errors"
	"net/http"
	"strings"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	gittypes "github.com/portainer/portainer/api/git/types"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"
	"github.com/portainer/portainer/pkg/libhttp/request"
	"github.com/portainer/portainer/pkg/libhttp/response"
)

var ErrNotGitSource = errors.New("source is not a Git source")

// @id GitOpsSourcesUpdateGit
// @summary Update a Git source
// @description Updates an existing GitOps source backed by a Git repository.
// @description **Access policy**: admin
// @tags gitops
// @security ApiKeyAuth
// @security jwt
// @accept json
// @produce json
// @param id path int true "Source identifier"
// @param body body GitSourceCreatePayload true "Git source details"
// @success 200 {object} portainer.Source
// @failure 400 "Invalid request payload"
// @failure 403 "Access denied"
// @failure 404 "Source not found"
// @failure 500 "Server error"
// @router /gitops/sources/{id} [put]
func (h *Handler) gitSourceUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	sourceID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid source identifier route variable", err)
	}

	var payload GitSourceCreatePayload

	if err := request.DecodeAndValidateJSONPayload(r, &payload); err != nil {
		return httperror.BadRequest("Invalid request payload", err)
	}

	var src *portainer.Source

	if err := h.dataStore.UpdateTx(func(tx dataservices.DataStoreTx) error {
		var err error

		if src, err = tx.Source().Read(portainer.SourceID(sourceID)); err != nil {
			return err
		}

		if src.Type != portainer.SourceTypeGit {
			return ErrNotGitSource
		}

		name := payload.Name
		if strings.TrimSpace(name) == "" {
			name = gittypes.RepoName(payload.URL)
		}

		src.Name = name

		var existingAuth *gittypes.GitAuthentication
		if src.GitConfig != nil {
			existingAuth = src.GitConfig.Authentication
		}

		src.GitConfig = &gittypes.RepoConfig{
			URL:           payload.URL,
			ReferenceName: payload.ReferenceName,
			TLSSkipVerify: payload.TLSSkipVerify,
		}

		if payload.Authentication != nil {
			src.GitConfig.Authentication = &gittypes.GitAuthentication{
				Username:          payload.Authentication.Username,
				Password:          payload.Authentication.Password,
				Provider:          payload.Authentication.Provider,
				AuthorizationType: payload.Authentication.AuthorizationType,
			}
		} else if payload.ClearAuthentication {
			src.GitConfig.Authentication = nil
		} else {
			src.GitConfig.Authentication = existingAuth
		}

		return tx.Source().Update(src.ID, src)
	}); h.dataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a source with the specified identifier", err)
	} else if errors.Is(err, ErrNotGitSource) {
		return httperror.BadRequest("Source is not a Git source", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to update source", err)
	}

	src.GitConfig = gittypes.SanitizeRepoConfig(src.GitConfig)

	return response.JSON(w, src)
}
