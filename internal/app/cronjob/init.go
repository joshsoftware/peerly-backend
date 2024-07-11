package cronjob

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation"
)

func InitializeJobs(appreciationSvc appreciation.Service,scheduler gocron.Scheduler) {
	
	DailyJob := NewDailyJob(appreciationSvc,scheduler)
	DailyJob.Schedule()
}