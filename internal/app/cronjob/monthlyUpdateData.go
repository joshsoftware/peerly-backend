package cronjob

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"

	logger "github.com/sirupsen/logrus"
)

const MONTHLY_JOB = "MONTHLY_JOB"
const MONTHLY_CRON_JOB_INTERVAL_MONTHS = 1

var MonthlyJobTiming = JobTime{
	hours:   16,
	minutes: 15,
	seconds: 0,
}

type MonthlyJob struct {
	CronJob
	userService user.Service
}

func NewMontlyJob(userSvc user.Service,scheduler gocron.Scheduler) Job {
	return &MonthlyJob{
		userService: userSvc,
		CronJob: CronJob{
			name:      MONTHLY_JOB,
			scheduler: scheduler,
		},
	}
}

func (cron *MonthlyJob) Schedule() {
	var err error
	cron.job, err = cron.scheduler.NewJob(
		gocron.MonthlyJob(MONTHLY_CRON_JOB_INTERVAL_MONTHS, gocron.NewDaysOfTheMonth(23),
			gocron.NewAtTimes(
				gocron.NewAtTime(
					MonthlyJobTiming.hours,
					MonthlyJobTiming.minutes,
					MonthlyJobTiming.seconds,
				),
			)),
		gocron.NewTask(cron.Execute, cron.Task),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	cron.scheduler.Start()
	if err != nil {
		logger.Warn(context.TODO(), "error occurred while scheduling %s, message %+v", cron.name, err.Error())
	}
}

func (cron *MonthlyJob) Task(ctx context.Context) {
	logger.Info(ctx, "in monthly job task")
	for i := 0; i < 3; i++ {
		logger.Info("cron job attempt:", i+1)
		err := cron.userService.UpdateRewardQuota(ctx)
		if err == nil {
			break
		}
	}

}
