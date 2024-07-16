package cronjob

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/app/users"
)

func InitializeJobs(appreciationSvc appreciation.Service,userSvc user.Service,scheduler gocron.Scheduler) {
	
	DailyJob := NewDailyJob(appreciationSvc,scheduler)
	DailyJob.Schedule()

	MonthlyJob := NewMontlyJob(userSvc,scheduler)
	MonthlyJob.Schedule()
}