package cronjob

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	orgSvc "github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	"github.com/joshsoftware/peerly-backend/internal/app/users"
)

func InitializeJobs(appreciationSvc appreciation.Service, userSvc user.Service, organizationConfigService orgSvc.Service, scheduler gocron.Scheduler) error {

	DailyJob := NewDailyJob(appreciationSvc, organizationConfigService, scheduler)
	err := DailyJob.Schedule()
	if err != nil {
		return err
	}
	MonthlyJob := NewMontlyJob(userSvc, organizationConfigService, scheduler)
	err = MonthlyJob.Schedule()
	if err != nil {
		return err
	}
	return nil
}
