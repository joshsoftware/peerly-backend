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
	versionHeader = "Accept"
	authHeader    = "X-Auth-Token"
)

// NewRouter initializes and returns a new router with the specified dependencies.
func NewRouter(deps app.Dependencies) *mux.Router {

	router := mux.NewRouter()

	// No version requirement for /ping
	router.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	v1 := fmt.Sprintf("application/vnd.%s.v1", config.AppName())

	router.Handle("/organizations", middleware.JwtAuthMiddleware(listOrganizationHandler(deps.OrganizationService))).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizations/{id:[0-9]+}", middleware.JwtAuthMiddleware(getOrganizationHandler(deps.OrganizationService))).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizations/{domainName}", middleware.JwtAuthMiddleware(getOrganizationByDomainNameHandler(deps.OrganizationService))).Methods(http.MethodGet).Headers(versionHeader, v1)

	router.Handle("/organizations", middleware.JwtAuthMiddleware(createOrganizationHandler(deps.OrganizationService))).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/organizations/{id:[0-9]+}", middleware.JwtAuthMiddleware(deleteOrganizationHandler(deps.OrganizationService))).Methods(http.MethodDelete).Headers(versionHeader, v1)

	router.Handle("/organizations/{id:[0-9]+}", middleware.JwtAuthMiddleware(updateOrganizationHandler(deps.OrganizationService))).Methods(http.MethodPut).Headers(versionHeader, v1)
	
	router.Handle("/organizations/otp/verify",middleware.JwtAuthMiddleware(OTPVerificationHandler(deps.OrganizationService))).Methods(http.MethodPost).Headers(versionHeader, v1)

	router.Handle("/organizations/otp/{id:[0-9]+}",middleware.JwtAuthMiddleware(ResendOTPhandler(deps.OrganizationService))).Methods(http.MethodPost).Headers(versionHeader, v1)

	return router
}
