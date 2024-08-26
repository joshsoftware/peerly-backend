package cronjob

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"

	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

const MONTHLY_JOB = "MONTHLY_JOB"
const MONTHLY_CRON_JOB_INTERVAL_MONTHS = 1

var MonthlyJobTiming = JobTime{
	hours:   0,
	minutes: 0,
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
		gocron.MonthlyJob(MONTHLY_CRON_JOB_INTERVAL_MONTHS, gocron.NewDaysOfTheMonth(-1),
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
		logger.Infof(context.TODO(), "error occurred while scheduling %s, message %+v", cron.name, err.Error())
	}
}

func (cron *MonthlyJob) Task(ctx context.Context) {
	logger.Info(ctx, "in monthly job task")
	for i := 0; i < 3; i++ {
		logger.Info(ctx,"cron job attempt:", i+1)
		err := cron.userService.UpdateRewardQuota(ctx)
		logger.Error(ctx,"err: ",err)
		if err == nil {
			sendRewardQuotaRefilledNotificationToAll()
			break
		}
	}

}

func sendRewardQuotaRefilledNotificationToAll() {
	msg := notification.Message{
		Title: "Reward Quota is Refilled",
		Body:  "Your reward quota is reset! You now recognize your colleagues.",
	}

	logger.Debug(context.Background(),"msg:",msg)
	msg.SendNotificationToTopic("peerly")
}

