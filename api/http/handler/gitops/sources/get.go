package sources

import (
	"net/http"

	portainer "github.com/portainer/portainer/api"
	gittypes "github.com/portainer/portainer/api/git/types"
	ce "github.com/portainer/portainer/api/gitops/workflows"
	"github.com/portainer/portainer/api/http/security"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"
	"github.com/portainer/portainer/pkg/libhttp/request"
	"github.com/portainer/portainer/pkg/libhttp/response"
)

type connectionInfo struct {
	ReferenceName  string `json:"referenceName"`
	ConfigFilePath string `json:"configFilePath"`
	TLSSkipVerify  bool   `json:"tlsSkipVerify"`
	Authentication bool   `json:"authentication,omitempty"`
}

type autoUpdateInfo struct {
	Mechanism     string `json:"mechanism,omitempty"`
	FetchInterval string `json:"fetchInterval,omitempty"`
}

// SourceDetail extends Source with connection settings and linked workflows.
type SourceDetail struct {
	Source
	Connection *connectionInfo `json:"connection,omitempty"`
	AutoUpdate *autoUpdateInfo `json:"autoUpdate,omitempty"`
	Workflows  []ce.Workflow   `json:"workflows"`
}

// @id GitOpsSourceGet
// @summary Get a GitOps source by ID
// @description Returns a single GitOps source with its connection settings and linked workflows.
// @description **Access policy**: admin
// @tags gitops
// @security ApiKeyAuth
// @security jwt
// @produce json
// @param id path string true "Source ID"
// @success 200 {object} SourceDetail
// @failure 403 "Access denied"
// @failure 404 "Source not found"
// @failure 500 "Server error"
// @router /gitops/sources/{id} [get]
func (h *Handler) getSource(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid source ID", err)
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve info from request context", err)
	}

	workflows, err := ce.FetchWorkflows(r.Context(), h.dataStore, h.gitService, h.k8sFactory, securityContext, nil)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve workflows", err)
	}

	byID := workflowsBySourceID(workflows)

	wfs, ok := byID[id]
	if !ok || len(wfs) == 0 || wfs[0].GitConfig == nil {
		return httperror.NotFound("Source not found", nil)
	}

	url := wfs[0].GitConfig.URL
	detail := SourceDetail{
		Source:     buildSource(id, url, wfs),
		Connection: buildConnectionInfo(wfs[0].GitConfig),
		AutoUpdate: buildAutoUpdateInfo(wfs[0].AutoUpdate),
		Workflows:  redactWorkflowCredentials(wfs),
	}
	return response.JSON(w, detail)
}

func buildConnectionInfo(cfg *gittypes.RepoConfig) *connectionInfo {
	if cfg == nil {
		return nil
	}
	return &connectionInfo{
		ReferenceName:  cfg.ReferenceName,
		ConfigFilePath: cfg.ConfigFilePath,
		TLSSkipVerify:  cfg.TLSSkipVerify,
		Authentication: cfg.Authentication != nil,
	}
}

func buildAutoUpdateInfo(autoUpdate *portainer.AutoUpdateSettings) *autoUpdateInfo {
	if autoUpdate == nil {
		return nil
	}

	switch {
	case autoUpdate.Interval != "":
		return &autoUpdateInfo{
			Mechanism:     "Interval",
			FetchInterval: autoUpdate.Interval,
		}
	case autoUpdate.Webhook != "":
		return &autoUpdateInfo{
			Mechanism: "Webhook",
		}
	default:
		return nil
	}
}
