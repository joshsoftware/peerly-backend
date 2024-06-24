package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/middleware"
)

const (
	versionHeader = "Accept-Version"
	authHeader    = "X-Auth-Token"
)

// NewRouter initializes and returns a new router with the specified dependencies.
func NewRouter(deps app.Dependencies) *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	//corevalues
	router.Handle("/organisations/{organisation_id:[0-9]+}/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(getCoreValueHandler(deps.CoreValueService))).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organisations/{organisation_id:[0-9]+}/core_values", middleware.JwtAuthMiddleware(listCoreValuesHandler(deps.CoreValueService))).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organisations/{organisation_id:[0-9]+}/core_values", middleware.JwtAuthMiddleware(createCoreValueHandler(deps.CoreValueService))).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/organisations/{organisation_id:[0-9]+}/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(deleteCoreValueHandler(deps.CoreValueService))).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.Handle("/organisations/{organisation_id:[0-9]+}/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateCoreValueHandler(deps.CoreValueService))).Methods(http.MethodPut).Headers(versionHeader, v1)

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	return router
}
