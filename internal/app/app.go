package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	corevalues "github.com/joshsoftware/peerly-backend/internal/app/coreValues"
	reward "github.com/joshsoftware/peerly-backend/internal/app/reward"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	organizationConfig "github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	repository "github.com/joshsoftware/peerly-backend/internal/repository/postgresdb"
)

// Dependencies holds the dependencies required by the application.
type Dependencies struct {
	CoreValueService    corevalues.Service
	AppreciationService appreciation.Service
	UserService         user.Service
	OrganizationConfigService organizationConfig.Service
	RewardService       reward.Service
}

// NewService initializes and returns a Dependencies instance with the given database connection.
func NewService(db *sqlx.DB) Dependencies {
	// Initialize repository dependencies using the provided database connection.

	coreValueRepo := repository.NewCoreValueRepo(db)
	userRepo := repository.NewUserRepo(db)
	appreciationRepo := repository.NewAppreciationRepo(db)
	orgConfigRepo := repository.NewOrganizationConfigRepo(db)
	rewardRepo := repository.NewRewardRepo(db)

	coreValueService := corevalues.NewService(coreValueRepo)
	userService := user.NewService(userRepo)
	appreciationService := appreciation.NewService(appreciationRepo, coreValueRepo)
	rewardService := reward.NewService(rewardRepo,appreciationRepo)
	orgConfigService := organizationConfig.NewService(orgConfigRepo)

	return Dependencies{
		CoreValueService:    coreValueService,
		AppreciationService: appreciationService,
		UserService:         userService,
		OrganizationConfigService: orgConfigService,
		RewardService:       rewardService,
	}

}
