package sources

import (
	"net/http"
	"net/http/httptest"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	gittypes "github.com/portainer/portainer/api/git/types"
	"github.com/portainer/portainer/api/http/security"
	"github.com/portainer/portainer/api/internal/testhelpers"

	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/require"
)

func newTestHandler(t *testing.T, store dataservices.DataStore) *Handler {
	t.Helper()
	return NewHandler(testhelpers.NewTestRequestBouncer(), store, nil, nil)
}

func buildListReq(t *testing.T, userID portainer.UserID, query string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/gitops/sources?"+query, nil)
	req = req.WithContext(security.StoreTokenData(req, &portainer.TokenData{ID: userID}))
	req = req.WithContext(security.StoreRestrictedRequestContext(req, &security.RestrictedRequestContext{
		UserID: userID, IsAdmin: true,
	}))
	return req
}

func buildGetReq(t *testing.T, userID portainer.UserID, id string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/gitops/sources/"+id, nil)
	req = req.WithContext(security.StoreTokenData(req, &portainer.TokenData{ID: userID}))
	req = req.WithContext(security.StoreRestrictedRequestContext(req, &security.RestrictedRequestContext{
		UserID: userID, IsAdmin: true,
	}))
	return req
}

func decodeSources(t *testing.T, rr *httptest.ResponseRecorder) []Source {
	t.Helper()
	require.Equal(t, http.StatusOK, rr.Code, "unexpected status: %s", rr.Body.String())
	var items []Source
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&items))
	return items
}

func decodeSourceDetail(t *testing.T, rr *httptest.ResponseRecorder) SourceDetail {
	t.Helper()
	require.Equal(t, http.StatusOK, rr.Code, "unexpected status: %s", rr.Body.String())
	var item SourceDetail
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&item))
	return item
}

func gitCfg(url string) *gittypes.RepoConfig {
	return &gittypes.RepoConfig{
		URL:            url,
		ConfigFilePath: "docker-compose.yml",
		ReferenceName:  "refs/heads/main",
	}
}
