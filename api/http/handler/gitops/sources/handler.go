package sources

import (
	"net/http"
	"time"

	gocache "github.com/patrickmn/go-cache"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/dataservices"
	"github.com/portainer/portainer/api/http/security"
	"github.com/portainer/portainer/api/kubernetes/cli"
	httperror "github.com/portainer/portainer/pkg/libhttp/error"

	"github.com/gorilla/mux"
)

const (
	cacheTTL             = 30 * time.Second
	cacheCleanupInterval = 10 * time.Minute
)

// Handler is the HTTP handler for the GitOps sources API.
type Handler struct {
	*mux.Router
	dataStore  dataservices.DataStore
	gitService portainer.GitService
	cache      *gocache.Cache
	k8sFactory *cli.ClientFactory
}

func NewHandler(bouncer security.BouncerService, dataStore dataservices.DataStore, gitService portainer.GitService, k8sFactory *cli.ClientFactory) *Handler {
	h := &Handler{
		Router:     mux.NewRouter(),
		dataStore:  dataStore,
		gitService: gitService,
		cache:      gocache.New(cacheTTL, cacheCleanupInterval),
		k8sFactory: k8sFactory,
	}

	adminRouter := h.PathPrefix("/gitops/sources").Subrouter()
	adminRouter.Use(bouncer.AdminAccess)
	adminRouter.Handle("", httperror.LoggerHandler(h.list)).Methods(http.MethodGet)
	adminRouter.Handle("/summary", httperror.LoggerHandler(h.summary)).Methods(http.MethodGet)
	adminRouter.Handle("/{id}", httperror.LoggerHandler(h.getSource)).Methods(http.MethodGet)
	return h
}
