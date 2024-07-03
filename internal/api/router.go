package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	intranet "github.com/joshsoftware/peerly-backend/internal/dummyIntranet"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
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
	router.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(getCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/core_values", middleware.JwtAuthMiddleware(listCoreValuesHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/core_values", middleware.JwtAuthMiddleware(createCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodPut).Headers(versionHeader, v1)

	//login

	router.Handle("/intranet/validate", intranet.ValidatePeerly()).Methods(http.MethodGet)

	router.Handle("/intranet/getuser/{user_id:[0-9]+}", intranet.IntranetGetUserApi()).Methods(http.MethodGet)

	router.Handle("/user/login", loginUser(deps.UserService)).Methods(http.MethodPost)

	//appreciations

	router.Handle("/appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(getAppreciationByIdHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/appreciations", middleware.JwtAuthMiddleware(getAppreciationsHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(validateAppreciationHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.Handle("/appreciation", middleware.JwtAuthMiddleware(createAppreciationHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	// organization config
	router.Handle("/organizationconfig", middleware.JwtAuthMiddleware(getOrganizationConfigHandler(deps.OrganizationService),[]string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizationconfig", middleware.JwtAuthMiddleware(createOrganizationConfigHandler(deps.OrganizationService),[]string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/organizationconfig", middleware.JwtAuthMiddleware(updateOrganizationConfigHandler(deps.OrganizationService),[]string{constants.UserRole})).Methods(http.MethodPut).Headers(versionHeader, v1)
	

	return router
}
