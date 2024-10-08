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

	// Add the RequestIDMiddleware to the subrouter
	peerlySubrouter.Use(middleware.RequestIDMiddleware)

	peerlySubrouter.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	peerlySubrouter.HandleFunc("/set_logger_level", loggerHandler).Methods(http.MethodPatch)

	// Version 1 API management
	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	//corevalues
	peerlySubrouter.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(getCoreValueHandler(deps.CoreValueService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/core_values", middleware.JwtAuthMiddleware(listCoreValuesHandler(deps.CoreValueService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/core_values", middleware.JwtAuthMiddleware(createCoreValueHandler(deps.CoreValueService), constants.Admin)).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/core_values/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateCoreValueHandler(deps.CoreValueService), constants.Admin)).Methods(http.MethodPut).Headers(versionHeader, v1)

	//users

	peerlySubrouter.Handle("/user/register", registerUser(deps.UserService)).Methods(http.MethodPost)

	peerlySubrouter.Handle("/user/login", loginUser(deps.UserService)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/admin/login", loginAdmin(deps.UserService)).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/intranet/users", listIntranetUsersHandler(deps.UserService)).Methods(http.MethodGet)

	peerlySubrouter.Handle("/users", middleware.JwtAuthMiddleware(listUsersHandler(deps.UserService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/user_profile", middleware.JwtAuthMiddleware(getUserByIdHandler(deps.UserService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/users/active", middleware.JwtAuthMiddleware(getActiveUserListHandler(deps.UserService), constants.User)).Methods(http.MethodGet)

	peerlySubrouter.Handle("/users/top10", middleware.JwtAuthMiddleware(getTop10UserHandler(deps.UserService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/admin/notification", middleware.JwtAuthMiddleware(adminNotificationHandler(deps.UserService), constants.Admin)).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/admin/appreciation_report", middleware.JwtAuthMiddleware(appreciationReportHandler(deps.UserService, deps.AppreciationService), constants.Admin)).Methods(http.MethodGet)

	peerlySubrouter.Handle("/admin/reported_appreciation_report", middleware.JwtAuthMiddleware(reportedAppreciationReportHandler(deps.UserService, deps.ReportAppreciationService), constants.Admin)).Methods(http.MethodGet)

	//appreciations

	peerlySubrouter.Handle("/appreciations/{id:[0-9]+}", middleware.JwtAuthMiddleware(getAppreciationByIDHandler(deps.AppreciationService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/appreciations", middleware.JwtAuthMiddleware(listAppreciationsHandler(deps.AppreciationService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/appreciations/{id:[0-9]+}", middleware.JwtAuthMiddleware(deleteAppreciationHandler(deps.AppreciationService), constants.Admin)).Methods(http.MethodDelete).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/appreciations", middleware.JwtAuthMiddleware(createAppreciationHandler(deps.AppreciationService), constants.User)).Methods(http.MethodPost).Headers(versionHeader, v1)

	//report appreciation
	peerlySubrouter.Handle("/report_appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(reportAppreciationHandler(deps.ReportAppreciationService), constants.User)).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/report_appreciations", middleware.JwtAuthMiddleware(listReportedAppreciations(deps.ReportAppreciationService), constants.Admin)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/moderate_appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(moderateAppreciation(deps.ReportAppreciationService), constants.Admin)).Methods(http.MethodPut).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/resolve_appreciation/{id:[0-9]+}", middleware.JwtAuthMiddleware(resolveAppreciation(deps.ReportAppreciationService), constants.Admin)).Methods(http.MethodPut).Headers(versionHeader, v1)

	//grades
	peerlySubrouter.Handle("/grades", middleware.JwtAuthMiddleware(listGradesHandler(deps.GradeService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/grades/{id:[0-9]+}", middleware.JwtAuthMiddleware(editGradesHandler(deps.GradeService), constants.Admin)).Methods(http.MethodPatch).Headers(versionHeader, v1)

	// reward appreciation
	peerlySubrouter.Handle("/reward/{id:[0-9]+}", middleware.JwtAuthMiddleware(giveRewardHandler(deps.RewardService), constants.User)).Methods(http.MethodPost).Headers(versionHeader, v1)

	// organization config
	peerlySubrouter.Handle("/organizationconfig", middleware.JwtAuthMiddleware(getOrganizationConfigHandler(deps.OrganizationConfigService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	//organization config data inserted by seed file
	// peerlySubrouter.Handle("/organizationconfig", middleware.JwtAuthMiddleware(createOrganizationConfigHandler(deps.OrganizationConfigService),[]string{constants.UserRole})).Methods(http.MethodPost).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/organizationconfig", middleware.JwtAuthMiddleware(updateOrganizationConfigHandler(deps.OrganizationConfigService), constants.Admin)).Methods(http.MethodPut).Headers(versionHeader, v1)

	//badges

	peerlySubrouter.Handle("/badges", middleware.JwtAuthMiddleware(listBadgesHandler(deps.BadgeService), constants.User)).Methods(http.MethodGet).Headers(versionHeader, v1)

	peerlySubrouter.Handle("/badges/{id:[0-9]+}", middleware.JwtAuthMiddleware(editBadgesHandler(deps.BadgeService), constants.Admin)).Methods(http.MethodPatch).Headers(versionHeader, v1)

	// No version requirement for /ping
	peerlySubrouter.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	sh := http.StripPrefix("/peerly/api_doc", http.FileServer(http.Dir("./apiDoc")))
	peerlySubrouter.PathPrefix("/api_doc").Handler(sh)

	// Serve static files from the "./assets" directory
	peerlySubrouter.PathPrefix("/assets/").Handler(http.StripPrefix("/peerly/assets/", http.FileServer(http.Dir("./assets"))))

	return router
}
