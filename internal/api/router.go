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

	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)
	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	//corevalues
	router.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(getCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/core_values", middleware.JwtAuthMiddleware(listCoreValuesHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/core_values", middleware.JwtAuthMiddleware(createCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateCoreValueHandler(deps.CoreValueService), []string{constants.UserRole})).Methods(http.MethodPut).Headers(versionHeader, v1)

	//login
	router.Handle("/user/register", registerUser(deps.UserService)).Methods(http.MethodPost)

	router.Handle("/user/login", loginUser(deps.UserService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/users", listIntranetUsersHandler(deps.UserService)).Methods(http.MethodGet)

	//badge
	router.Handle("/badges",middleware.JwtAuthMiddleware(createBadgeHandler(deps.BadgeService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)
	
	router.Handle("/badges",middleware.JwtAuthMiddleware(listBadgesHandler(deps.BadgeService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/badges/{id:[0-9]+}",middleware.JwtAuthMiddleware(getBadgeHandler(deps.BadgeService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/badges/{id:[0-9]+}",middleware.JwtAuthMiddleware(deleteBadgeHandler(deps.BadgeService), []string{constants.UserRole})).Methods(http.MethodDelete).Headers(versionHeader, v1)
	
	router.Handle("/badges/{id:[0-9]+}",middleware.JwtAuthMiddleware(updateBadgeHandler(deps.BadgeService), []string{constants.UserRole})).Methods(http.MethodPatch).Headers(versionHeader, v1)

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	return router
}
