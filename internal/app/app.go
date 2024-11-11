package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/app/badges"
	corevalues "github.com/joshsoftware/peerly-backend/internal/app/coreValues"
	"github.com/joshsoftware/peerly-backend/internal/app/grades"
	reportappreciations "github.com/joshsoftware/peerly-backend/internal/app/reportAppreciations"

	organizationConfig "github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	reward "github.com/joshsoftware/peerly-backend/internal/app/reward"

	user "github.com/joshsoftware/peerly-backend/internal/app/users"

	appreciation "github.com/joshsoftware/peerly-backend/internal/app/appreciation"

	repository "github.com/joshsoftware/peerly-backend/internal/repository/postgresdb"
)

// Dependencies holds the dependencies required by the application.
type Dependencies struct {
	CoreValueService          corevalues.Service
	AppreciationService       appreciation.Service
	UserService               user.Service
	ReportAppreciationService reportappreciations.Service
	RewardService             reward.Service
	GradeService              grades.Service
	OrganizationConfigService organizationConfig.Service
	BadgeService              badges.Service
}

// NewService initializes and returns a Dependencies instance with the given database connection.
func NewService(db *sqlx.DB) Dependencies {
	// Initialize repository dependencies using the provided database connection.

	coreValueRepo := repository.NewCoreValueRepo(db)
	userRepo := repository.NewUserRepo(db)
	reportAppreciationRepo := repository.NewReportRepo(db)
	appreciationRepo := repository.NewAppreciationRepo(db)
	rewardRepo := repository.NewRewardRepo(db)
	gradeRepo := repository.NewGradesRepo(db)
	orgConfigRepo := repository.NewOrganizationConfigRepo(db)
	badgeRepo := repository.NewBadgeRepo(db)

	coreValueService := corevalues.NewService(coreValueRepo)
	appreciationService := appreciation.NewService(appreciationRepo, coreValueRepo, userRepo)
	userService := user.NewService(userRepo)
	reportAppreciationService := reportappreciations.NewService(reportAppreciationRepo, userRepo, appreciationRepo)
	rewardService := reward.NewService(rewardRepo, appreciationRepo, userRepo, reportAppreciationRepo)
	gradeService := grades.NewService(gradeRepo, userRepo)
	orgConfigService := organizationConfig.NewService(orgConfigRepo)
	badgeService := badges.NewService(badgeRepo, userRepo)

	return Dependencies{
		CoreValueService:          coreValueService,
		AppreciationService:       appreciationService,
		UserService:               userService,
		ReportAppreciationService: reportAppreciationService,
		RewardService:             rewardService,
		GradeService:              gradeService,
		OrganizationConfigService: orgConfigService,
		BadgeService:              badgeService,
	}

}
