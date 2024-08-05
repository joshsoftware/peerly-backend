package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
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
	peerlySubrouter := router.PathPrefix("/peerly").Subrouter()

	peerlySubrouter.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	//corevalues
	peerlySubrouter.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(getCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/core_values", middleware.JwtAuthMiddleware(listCoreValuesHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/core_values", middleware.JwtAuthMiddleware(createCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodPut).Headers(versionHeader, v1)

	//login
	peerlySubrouter.Handle("/user/register", registerUser(deps.UserService)).Methods(http.MethodPost)

	peerlySubrouter.Handle("/user/login", loginUser(deps.UserService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/users", listIntranetUsersHandler(deps.UserService)).Methods(http.MethodGet)

	peerlySubrouter.Handle("/users/all", middleware.JwtAuthMiddleware(listUsersHandler(deps.UserService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	//appreciations

	peerlySubrouter.Handle("/appreciations/{id:[0-9]+}", middleware.JwtAuthMiddleware(getAppreciationByIDHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/appreciations", middleware.JwtAuthMiddleware(listAppreciationsHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/appreciations/{id:[0-9]+}", middleware.JwtAuthMiddleware(deleteAppreciationHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodDelete).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/appreciations", middleware.JwtAuthMiddleware(createAppreciationHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)
	// organization config
	peerlySubrouter.Handle("/organizationconfig", middleware.JwtAuthMiddleware(getOrganizationConfigHandler(deps.OrganizationConfigService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	//organization config data inserted by seed file
	// peerlySubrouter.Handle("/organizationconfig", middleware.JwtAuthMiddleware(createOrganizationConfigHandler(deps.OrganizationConfigService),[]string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/organizationconfig", middleware.JwtAuthMiddleware(updateOrganizationConfigHandler(deps.OrganizationConfigService), []string{constants.UserRole})).Methods(http.MethodPut).Headers(versionHeader, v1)

	return router
}
