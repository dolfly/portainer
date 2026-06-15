package sources

import (
	"net/http"
	"testing"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/datastore"
	gittypes "github.com/portainer/portainer/api/git/types"
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

	_, httpErr := ValidateGitSourceAccess(store, src.ID)
	assert.Nil(t, httpErr)
}

func TestValidateSourceForStack_SourceNotFound_Returns404(t *testing.T) {
	t.Parallel()
	_, store := datastore.MustNewTestStore(t, false, false)

	_, httpErr := ValidateGitSourceAccess(store, portainer.SourceID(999))
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

	_, httpErr := ValidateGitSourceAccess(store, src.ID)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.StatusCode)
}
