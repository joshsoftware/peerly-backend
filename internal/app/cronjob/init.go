package cronjob

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/app/googlesheets"
	orgSvc "github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	reportappreciations "github.com/joshsoftware/peerly-backend/internal/app/reportAppreciations"
	"github.com/joshsoftware/peerly-backend/internal/app/users"
)

func InitializeJobs(appreciationSvc appreciation.Service, userSvc user.Service, organizationConfigService orgSvc.Service, reportAppreciationSvc reportappreciations.Service, sheetSvc *googlesheets.Service, spreadsheetID string, scheduler gocron.Scheduler) error {

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

	if sheetSvc != nil {
		QuarterlyJob := NewQuarterlyReportJob(appreciationSvc, reportAppreciationSvc, sheetSvc, spreadsheetID, scheduler)
		err = QuarterlyJob.Schedule()
		if err != nil {
			return err
		}
	}

	return nil
}
