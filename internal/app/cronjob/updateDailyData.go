package cronjob

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	apprSvc "github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	logger "github.com/sirupsen/logrus"
)


const DAILY_JOB = "DAILY_JOB"
const SAY_HELLO_DAILY_CRON_JOB_INTERVAL_DAYS = 1

var SayHelloDailyJobTiming = JobTime{
	hours:   12,
	minutes: 40,
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
		logger.Warn(context.TODO(), "error occurred while scheduling %s, message %+v", cron.name, err.Error())
	}
}

func (cron *DailyJob) Task(ctx context.Context) {
	logger.Info(ctx, "in daily job task")
	for  i:=0;i<3;i++{
		logger.Info("cron job attempt:",i+1)
		isSuccess,err := cron.appreciationService.UpdateAppreciation(ctx)
		if err==nil && isSuccess{
			break
		}
	}
	
}



