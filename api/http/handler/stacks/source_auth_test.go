package stacks

import (
	"net/http"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/datastore"
	gittypes "github.com/portainer/portainer/api/git/types"
	"github.com/portainer/portainer/api/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateSourceForStack_ValidGitSource_ReturnsNil(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, false)

	src := &portainer.Source{
		Type: portainer.SourceTypeGit,
		Git:  &gittypes.RepoConfig{URL: "https://github.com/org/repo"},
	}
	require.NoError(t, store.Source().Create(src))

	handler := NewHandler(testhelpers.NewTestRequestBouncer())
	handler.DataStore = store

	_, httpErr := validateSourceForStack(store, src.ID)
	assert.Nil(t, httpErr)
}

func TestValidateSourceForStack_SourceNotFound_Returns404(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, false)

	handler := NewHandler(testhelpers.NewTestRequestBouncer())
	handler.DataStore = store

	_, httpErr := validateSourceForStack(store, portainer.SourceID(999))
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusNotFound, httpErr.StatusCode)
}

func TestValidateSourceForStack_NonGitSource_Returns400(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, false)

	src := &portainer.Source{
		Type: portainer.SourceType(99), // not a git source
	}
	require.NoError(t, store.Source().Create(src))

	handler := NewHandler(testhelpers.NewTestRequestBouncer())
	handler.DataStore = store

	_, httpErr := validateSourceForStack(store, src.ID)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
}

func TestComposeGitPayload_ValidateWithSourceID_URLNotRequired(t *testing.T) {
	t.Parallel()
	payload := &composeStackFromGitRepositoryPayload{
		Name:     "mystack",
		SourceID: portainer.SourceID(1),
		// RepositoryURL intentionally omitted
	}

	err := payload.Validate(nil)
	assert.NoError(t, err)
}

func TestComposeGitPayload_ValidateWithoutSourceID_URLRequired(t *testing.T) {
	t.Parallel()
	payload := &composeStackFromGitRepositoryPayload{
		Name: "mystack",
		// SourceID and RepositoryURL both omitted
	}

	err := payload.Validate(nil)
	assert.Error(t, err)
}

func TestSwarmGitPayload_ValidateWithSourceID_URLNotRequired(t *testing.T) {
	t.Parallel()
	payload := &swarmStackFromGitRepositoryPayload{
		Name:     "myswarm",
		SwarmID:  "swarm-abc",
		SourceID: portainer.SourceID(1),
	}

	err := payload.Validate(nil)
	assert.NoError(t, err)
}

func TestSwarmGitPayload_ValidateWithoutSourceID_URLRequired(t *testing.T) {
	t.Parallel()
	payload := &swarmStackFromGitRepositoryPayload{
		Name:    "myswarm",
		SwarmID: "swarm-abc",
	}

	err := payload.Validate(nil)
	assert.Error(t, err)
}
