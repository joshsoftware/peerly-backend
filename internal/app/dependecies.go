package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	corevalues "github.com/joshsoftware/peerly-backend/internal/app/coreValues"
	repo "github.com/joshsoftware/peerly-backend/internal/repository/postgresdb"

	user "github.com/joshsoftware/peerly-backend/internal/app/users"

	repository "github.com/joshsoftware/peerly-backend/internal/repository/postgresdb"
)

// Dependencies holds the dependencies required by the application.
type Dependencies struct {
	CoreValueService    corevalues.Service
	AppreciationService appreciation.Service
	UserService         user.Service
}

// NewServices initializes and returns a Dependencies instance with the given database connection.
func NewServices(db *sqlx.DB) Dependencies {
	// Initialize repository dependencies using the provided database connection.

	coreValueRepo := repository.NewCoreValueRepo(db)
	userRepo := repository.NewUserRepo(db)
	coreValueService := corevalues.NewService(coreValueRepo)

	appreciationRepo := repo.NewAppreciationRepo(db)
	appreciationService := appreciation.NewService(appreciationRepo, coreValueRepo)
	userService := user.NewService(userRepo)

	return Dependencies{
		CoreValueService:    coreValueService,
		AppreciationService: appreciationService,
		UserService:         userService,
	}

}
