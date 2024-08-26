package cronjob

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	apprSvc "github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	orgSvc "github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
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
	appreciationService       apprSvc.Service
	organizationConfigService orgSvc.Service
}

func NewDailyJob(
	appreciationService apprSvc.Service,
	organizationConfigService orgSvc.Service,
	scheduler gocron.Scheduler,
) Job {
	return &DailyJob{
		appreciationService:       appreciationService,
		organizationConfigService: organizationConfigService,
		CronJob: CronJob{
			name:      DAILY_JOB,
			scheduler: scheduler,
		},
	}
}

func (cron *DailyJob) Schedule() error {
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
		logger.Warn(context.Background(), "error occurred while scheduling %s, message %+v", cron.name, err.Error())
		return err
	}
	return nil
}
func (cron *DailyJob) Task(ctx context.Context) {
	ctx = context.WithValue(ctx, constants.RequestID, "dailyUpdate")
	logger.Info(ctx, "in daily job task")

	orgInfo, err := cron.organizationConfigService.GetOrganizationConfig(ctx)
	if err != nil {
		logger.Info(ctx,fmt.Sprintf("daily cron job err: %v ", err))
		return
	}
	for i := 0; i < 3; i++ {
		logger.Info(ctx, "cron job attempt:", i+1)
		isSuccess, err := cron.appreciationService.UpdateAppreciation(ctx, orgInfo.Timezone)

		logger.Info(ctx, " isSuccess: ", isSuccess, " err: ", err)
		if err == nil && isSuccess {
			break
		}
		logger.Info(ctx,fmt.Sprintf("daily cron job err: %v ", err))
	}
}
