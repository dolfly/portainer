package stackbuilders

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/datastore"
	httperrors "github.com/portainer/portainer/api/http/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Stubs

type stubFileContentBuilder struct {
	FileContentMethodStackBuildProcess
	store      *datastore.Store
	savedStack *portainer.Stack
	saveErr    error
	deployErr  error
}

func (s *stubFileContentBuilder) SetGeneralInfo(_ *StackPayload, _ *portainer.Endpoint) FileContentMethodStackBuildProcess {
	s.savedStack.Status = portainer.StackStatusDeploying
	s.savedStack.DeploymentStatus = []portainer.StackDeploymentStatus{
		{Status: portainer.StackStatusDeploying, Time: time.Now().Unix()},
	}
	return s
}
func (s *stubFileContentBuilder) SetUniqueInfo(_ *StackPayload) FileContentMethodStackBuildProcess {
	return s
}
func (s *stubFileContentBuilder) SetFileContent(_ *StackPayload) FileContentMethodStackBuildProcess {
	return s
}
func (s *stubFileContentBuilder) SaveStack() (*portainer.Stack, error) {
	if s.saveErr != nil {
		return nil, s.saveErr
	}
	return s.savedStack, s.store.Stack().Create(s.savedStack)
}
func (s *stubFileContentBuilder) Deploy(_ context.Context, _ *StackPayload, _ *portainer.Endpoint) FileContentMethodStackBuildProcess {
	return s
}
func (s *stubFileContentBuilder) Error() error { return s.deployErr }

type stubFileUploadBuilder struct {
	FileUploadMethodStackBuildProcess
	store      *datastore.Store
	savedStack *portainer.Stack
	saveErr    error
	deployErr  error
}

func (s *stubFileUploadBuilder) SetGeneralInfo(_ *StackPayload, _ *portainer.Endpoint) FileUploadMethodStackBuildProcess {
	s.savedStack.Status = portainer.StackStatusDeploying
	s.savedStack.DeploymentStatus = []portainer.StackDeploymentStatus{
		{Status: portainer.StackStatusDeploying, Time: time.Now().Unix()},
	}
	return s
}
func (s *stubFileUploadBuilder) SetUniqueInfo(_ *StackPayload) FileUploadMethodStackBuildProcess {
	return s
}
func (s *stubFileUploadBuilder) SetUploadedFile(_ *StackPayload) FileUploadMethodStackBuildProcess {
	return s
}
func (s *stubFileUploadBuilder) SaveStack() (*portainer.Stack, error) {
	if s.saveErr != nil {
		return nil, s.saveErr
	}
	return s.savedStack, s.store.Stack().Create(s.savedStack)
}
func (s *stubFileUploadBuilder) Deploy(_ context.Context, _ *StackPayload, _ *portainer.Endpoint) FileUploadMethodStackBuildProcess {
	return s
}
func (s *stubFileUploadBuilder) Error() error { return s.deployErr }

type stubUrlBuilder struct {
	UrlMethodStackBuildProcess
	store      *datastore.Store
	savedStack *portainer.Stack
	saveErr    error
	deployErr  error
}

func (s *stubUrlBuilder) SetGeneralInfo(_ *StackPayload, _ *portainer.Endpoint) UrlMethodStackBuildProcess {
	s.savedStack.Status = portainer.StackStatusDeploying
	s.savedStack.DeploymentStatus = []portainer.StackDeploymentStatus{
		{Status: portainer.StackStatusDeploying, Time: time.Now().Unix()},
	}
	return s
}
func (s *stubUrlBuilder) SetUniqueInfo(_ *StackPayload) UrlMethodStackBuildProcess { return s }
func (s *stubUrlBuilder) SetURL(_ *StackPayload) UrlMethodStackBuildProcess        { return s }
func (s *stubUrlBuilder) SaveStack() (*portainer.Stack, error) {
	if s.saveErr != nil {
		return nil, s.saveErr
	}
	return s.savedStack, s.store.Stack().Create(s.savedStack)
}
func (s *stubUrlBuilder) Deploy(_ context.Context, _ *StackPayload, _ *portainer.Endpoint) UrlMethodStackBuildProcess {
	return s
}
func (s *stubUrlBuilder) Error() error { return s.deployErr }

type stubGitBuilder struct {
	GitMethodStackBuildProcess
	store      *datastore.Store
	savedStack *portainer.Stack
	saveErr    error
	deployErr  error
	hookCalled bool
}

func (s *stubGitBuilder) SetGeneralInfo(_ *StackPayload, _ *portainer.Endpoint) GitMethodStackBuildProcess {
	return s
}
func (s *stubGitBuilder) SetUniqueInfo(_ *StackPayload) GitMethodStackBuildProcess { return s }
func (s *stubGitBuilder) SetGitRepository(_ context.Context, _ *StackPayload) GitMethodStackBuildProcess {
	return s
}
func (s *stubGitBuilder) SaveStack() (*portainer.Stack, error) {
	if s.saveErr != nil {
		return nil, s.saveErr
	}
	return s.savedStack, s.store.Stack().Create(s.savedStack)
}
func (s *stubGitBuilder) Deploy(_ context.Context, _ *StackPayload, _ *portainer.Endpoint) GitMethodStackBuildProcess {
	return s
}
func (s *stubGitBuilder) Error() error { return s.deployErr }
func (s *stubGitBuilder) EnableAutoUpdate(_ context.Context, _ *portainer.Stack) error {
	s.hookCalled = true
	return nil
}

// Helpers

func waitForStackStatus(t *testing.T, store *datastore.Store, id portainer.StackID, wantStatus portainer.StackStatus) *portainer.Stack {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		stack, err := store.Stack().Read(id)
		if err == nil && stack.Status == wantStatus {
			return stack
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for stack %d to reach status %d", id, wantStatus)
	return nil
}

// Tests

func TestDirector_Build_UnknownBuilder_ReturnsBadRequest(t *testing.T) {
	director := NewStackBuilderDirector(nil, "not a builder")

	_, herr := director.Build(t.Context(), &StackPayload{}, &portainer.Endpoint{})

	require.NotNil(t, herr)
	assert.Equal(t, http.StatusBadRequest, herr.StatusCode)
}

func TestDirector_Build_GitMethod_UnauthorizedCredential_ReturnsForbidden(t *testing.T) {
	director := NewStackBuilderDirector(nil, &stubGitBuilder{saveErr: httperrors.ErrUnauthorized})

	_, herr := director.Build(t.Context(), &StackPayload{}, &portainer.Endpoint{})

	require.NotNil(t, herr)
	assert.Equal(t, http.StatusForbidden, herr.StatusCode)
}

func TestDirector_Build_GitMethod_SaveError_ReturnsInternalServerError(t *testing.T) {
	director := NewStackBuilderDirector(nil, &stubGitBuilder{saveErr: errors.New("db error")})

	_, herr := director.Build(context.TODO(), &StackPayload{}, &portainer.Endpoint{})

	require.NotNil(t, herr)
	assert.Equal(t, http.StatusInternalServerError, herr.StatusCode)
}

func TestDirector_SpawnAsync_DeploySuccess_UpdatesStackStatusToActive(t *testing.T) {
	tc := []struct {
		name    string
		builder func(store *datastore.Store, stack *portainer.Stack) any
	}{
		{
			name: "file content builder",
			builder: func(store *datastore.Store, stack *portainer.Stack) any {
				return &stubFileContentBuilder{store: store, savedStack: stack}
			},
		},
		{
			name: "file upload builder",
			builder: func(store *datastore.Store, stack *portainer.Stack) any {
				return &stubFileUploadBuilder{store: store, savedStack: stack}
			},
		},
		{
			name: "url builder",
			builder: func(store *datastore.Store, stack *portainer.Stack) any {
				return &stubUrlBuilder{store: store, savedStack: stack}
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			_, store := datastore.MustNewTestStore(t, true, false)
			stack := &portainer.Stack{ID: 1}
			director := NewStackBuilderDirector(store, tt.builder(store, stack))
			_, herr := director.Build(context.TODO(), &StackPayload{}, &portainer.Endpoint{})
			require.Nil(t, herr)

			updated := waitForStackStatus(t, store, stack.ID, portainer.StackStatusActive)

			assert.Equal(t, portainer.StackStatusActive, updated.Status)
			require.Len(t, updated.DeploymentStatus, 2)
			assert.Equal(t, portainer.StackStatusDeploying, updated.DeploymentStatus[0].Status)
			assert.Equal(t, portainer.StackStatusActive, updated.DeploymentStatus[1].Status)
		})
	}
}

func TestDirector_SpawnAsync_DeployFailure_UpdatesStackStatusToError(t *testing.T) {
	deployErr := errors.New("failed to pull image nginx:999")
	_, store := datastore.MustNewTestStore(t, true, false)
	stack := &portainer.Stack{ID: 1}
	director := NewStackBuilderDirector(store, &stubFileContentBuilder{store: store, savedStack: stack, deployErr: deployErr})
	_, herr := director.Build(context.TODO(), &StackPayload{}, &portainer.Endpoint{})
	require.Nil(t, herr)

	updated := waitForStackStatus(t, store, stack.ID, portainer.StackStatusError)

	assert.Equal(t, portainer.StackStatusError, updated.Status)
	require.Len(t, updated.DeploymentStatus, 2)
	assert.Equal(t, portainer.StackStatusDeploying, updated.DeploymentStatus[0].Status)
	lastEntry := updated.DeploymentStatus[1]
	assert.Equal(t, portainer.StackStatusError, lastEntry.Status)
	assert.Equal(t, deployErr.Error(), lastEntry.Message)
}

func TestDirector_SpawnAsync_PostDeployHook_CalledOnSuccess(t *testing.T) {
	_, store := datastore.MustNewTestStore(t, true, false)
	stack := &portainer.Stack{ID: 1}
	builder := &stubGitBuilder{store: store, savedStack: stack}
	director := NewStackBuilderDirector(store, builder)
	_, herr := director.Build(context.TODO(), &StackPayload{}, &portainer.Endpoint{})
	require.Nil(t, herr)

	waitForStackStatus(t, store, stack.ID, portainer.StackStatusActive)

	assert.True(t, builder.hookCalled, "post-deploy hook should be called after a successful deployment")
}

func TestDirector_SpawnAsync_PostDeployHook_NotCalledOnDeployFailure(t *testing.T) {
	_, store := datastore.MustNewTestStore(t, true, false)
	stack := &portainer.Stack{ID: 1}
	builder := &stubGitBuilder{store: store, savedStack: stack, deployErr: errors.New("failed to deploy")}
	director := NewStackBuilderDirector(store, builder)
	_, herr := director.Build(context.TODO(), &StackPayload{}, &portainer.Endpoint{})
	require.Nil(t, herr)

	waitForStackStatus(t, store, stack.ID, portainer.StackStatusError)

	assert.False(t, builder.hookCalled, "post-deploy hook should not be called after a failed deployment")
}
