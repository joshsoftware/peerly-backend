package cronjob

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	apprSvc "github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

const DAILY_JOB = "DAILY_JOB"
const SAY_HELLO_DAILY_CRON_JOB_INTERVAL_DAYS = 1

var SayHelloDailyJobTiming = JobTime{
	hours:   0,
	minutes: 0,
	seconds: 0,
}

type DailyJob struct {
	CronJob
	appreciationService apprSvc.Service
}

func NewDailyJob(
	appreciationService apprSvc.Service,
	scheduler gocron.Scheduler,
) Job {
	return &DailyJob{
		appreciationService: appreciationService,
		CronJob: CronJob{
			name:      DAILY_JOB,
			scheduler: scheduler,
		},
	}
}

func (cron *DailyJob) Schedule() {
	var err error
	cron.job, err = cron.scheduler.NewJob(
		gocron.DailyJob(
			SAY_HELLO_DAILY_CRON_JOB_INTERVAL_DAYS,
			gocron.NewAtTimes(
				gocron.NewAtTime(
					SayHelloDailyJobTiming.hours,
					SayHelloDailyJobTiming.minutes,
					SayHelloDailyJobTiming.seconds,
				),
			),
		),
		gocron.NewTask(cron.Execute, cron.Task),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	cron.scheduler.Start()

	if err != nil {
		logger.Warn(context.TODO(), fmt.Sprintf("error occurred while scheduling %s, message %+v", cron.name, err.Error()))
	}
}
func (cron *DailyJob) Task(ctx context.Context) {
	ctx = context.WithValue(ctx, constants.RequestID, "dailyUpdate")
	logger.Info(ctx, "in daily job task")
	for  i:=0;i<3;i++{
		logger.Info(ctx,"cron job attempt:",i+1)
		isSuccess,err := cron.appreciationService.UpdateAppreciation(ctx)

		logger.Info(ctx," isSuccess: ",isSuccess," err: ",err)
		if err==nil && isSuccess{
			logger.Info(ctx,"cronjob: UpdateDaily data completed")
			break
		}
	}

}
