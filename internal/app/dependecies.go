package app

import (
	"github.com/jmoiron/sqlx"
	corevalues "github.com/joshsoftware/peerly-backend/internal/app/coreValues"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	repo "github.com/joshsoftware/peerly-backend/internal/repository/postgresdb"
)

// Dependencies holds the dependencies required by the application.
type Dependencies struct {
	CoreValueService corevalues.Service
    AppreciationService appreciation.Service
}

// NewServices initializes and returns a Dependencies instance with the given database connection.
func NewServices(db *sqlx.DB) Dependencies {
	// Initialize repository dependencies using the provided database connection.

	coreValueRepo := repository.NewCoreValueRepo(db)
	coreValueService := corevalues.NewService(coreValueRepo)

    appreciationRepo := repo.NewAppreciationRepo(db)
    appreciationService := appreciation.NewService(appreciationRepo,coreValueRepo)

    return Dependencies{
		CoreValueService: coreValueService,
        AppreciationService: appreciationService,
    }
}
