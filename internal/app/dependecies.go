package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/app/organization"
	"github.com/joshsoftware/peerly-backend/internal/repository/postgresdb"
)

// Dependencies holds the dependencies required by the application.
type Dependencies struct {
    OrganizationService organization.Service
}

// NewServices initializes and returns a Dependencies instance with the given database connection.
func NewServices(db *sqlx.DB) Dependencies {
    // Initialize repository dependencies using the provided database connection.

    orgRepo := repository.NewOrganizationRepo(db)
	otpRepo := repository.NewOTPVerificationRepo(db)
	orgService := organization.NewService(orgRepo,otpRepo)
    return Dependencies{
        OrganizationService: orgService,
    }
}