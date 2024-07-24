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

	router.Handle("/users/all", middleware.JwtAuthMiddleware(getUsersHandler(deps.UserService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/user_profile", middleware.JwtAuthMiddleware(getUserByIdHandler(deps.UserService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/users/active", getActiveUserListHandler(deps.UserService)).Methods(http.MethodGet)

	router.Handle("/users/top10", middleware.JwtAuthMiddleware(getTop10UserHandler(deps.UserService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	//appreciations

	router.Handle("/appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(getAppreciationByIDHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/appreciations", middleware.JwtAuthMiddleware(getAppreciationsHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(deleteAppreciationHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.Handle("/appreciation", middleware.JwtAuthMiddleware(createAppreciationHandler(deps.AppreciationService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	//report appreciation
	router.Handle("/report_appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(reportAppreciationHandler(deps.ReportAppreciationService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	// reward appreciation
	router.Handle("/reward/{id:[0-9]+}", middleware.JwtAuthMiddleware(giveRewardHandler(deps.RewardService), []string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	sh := http.StripPrefix("/api_doc", http.FileServer(http.Dir("./apiDoc")))
	router.PathPrefix("/api_doc").Handler(sh)

	return router
}
