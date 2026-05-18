package sources

import (
	"net/http/httptest"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	"github.com/portainer/portainer/api/datastore"
	gittypes "github.com/portainer/portainer/api/git/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourcesList_GroupsByURLAndCredentials(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, true)

	require.NoError(t, store.UpdateTx(func(tx dataservices.DataStoreTx) error {
		require.NoError(t, tx.Stack().Create(&portainer.Stack{
			ID: 1, Name: "stack-a", GitConfig: gitCfg("https://github.com/org/repo"),
		}))
		require.NoError(t, tx.Stack().Create(&portainer.Stack{
			ID: 2, Name: "stack-b", GitConfig: gitCfg("https://github.com/org/repo"),
		}))
		return tx.User().Create(&portainer.User{ID: 1, Role: portainer.AdministratorRole})
	}))

	h := newTestHandler(t, store)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, buildListReq(t, 1, ""))

	sources := decodeSources(t, rr)
	require.Len(t, sources, 1)
	assert.Equal(t, 2, sources[0].UsedBy)
}

func TestSourcesList_SeparatesCredentials(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, true)

	require.NoError(t, store.UpdateTx(func(tx dataservices.DataStoreTx) error {
		cfg1 := gitCfg("https://github.com/org/repo")
		cfg1.Authentication = &gittypes.GitAuthentication{Username: "alice", Password: "pass1"}
		cfg2 := gitCfg("https://github.com/org/repo")
		cfg2.Authentication = &gittypes.GitAuthentication{Username: "bob", Password: "pass2"}
		require.NoError(t, tx.Stack().Create(&portainer.Stack{ID: 1, Name: "stack-a", GitConfig: cfg1}))
		require.NoError(t, tx.Stack().Create(&portainer.Stack{ID: 2, Name: "stack-b", GitConfig: cfg2}))
		return tx.User().Create(&portainer.User{ID: 1, Role: portainer.AdministratorRole})
	}))

	h := newTestHandler(t, store)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, buildListReq(t, 1, ""))

	sources := decodeSources(t, rr)
	assert.Len(t, sources, 2)
}

func TestSourcesList_StatusFilter(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, true)

	// With nil gitService, source git-phase status is always StatusUnknown.
	require.NoError(t, store.UpdateTx(func(tx dataservices.DataStoreTx) error {
		require.NoError(t, tx.Stack().Create(&portainer.Stack{
			ID: 1, GitConfig: gitCfg("https://github.com/org/app"),
		}))
		return tx.User().Create(&portainer.User{ID: 1, Role: portainer.AdministratorRole})
	}))

	h := newTestHandler(t, store)

	t.Run("status=unknown matches sources with unknown status", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, buildListReq(t, 1, "status=unknown"))
		sources := decodeSources(t, rr)
		assert.Len(t, sources, 1)
	})

	t.Run("status=healthy excludes sources with unknown status", func(t *testing.T) {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, buildListReq(t, 1, "status=healthy"))
		sources := decodeSources(t, rr)
		assert.Empty(t, sources)
	})
}

func TestSourcesList_SearchByURL(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, true)

	require.NoError(t, store.UpdateTx(func(tx dataservices.DataStoreTx) error {
		require.NoError(t, tx.Stack().Create(&portainer.Stack{
			ID: 1, GitConfig: gitCfg("https://github.com/org/app"),
		}))
		require.NoError(t, tx.Stack().Create(&portainer.Stack{
			ID: 2, GitConfig: gitCfg("https://github.com/org/infra"),
		}))
		return tx.User().Create(&portainer.User{ID: 1, Role: portainer.AdministratorRole})
	}))

	h := newTestHandler(t, store)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, buildListReq(t, 1, "search=app"))

	sources := decodeSources(t, rr)
	require.Len(t, sources, 1)
	assert.Equal(t, "app", sources[0].Name)
}
