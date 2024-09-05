package cronjob

import (
	"context"
	"fmt"
	"time"

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
	// Load the location for Asia/Kolkata
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		logger.Warn(context.TODO(), "error loading location: %+v", err.Error())
		return err
	}

	// Get the current date in Asia/Kolkata
	currentTimeInKolkata := time.Now().In(location)

	// Create a new time for today's date with MonthlyJobTiming hours, minutes, and seconds
	jobTimeInKolkata := time.Date(
		currentTimeInKolkata.Year(),         // Year: current year
		currentTimeInKolkata.Month(),        // Month: current month
		currentTimeInKolkata.Day(),          // Day: today's date
		int(SayHelloDailyJobTiming.hours),   // Hour: from MonthlyJobTiming
		int(SayHelloDailyJobTiming.minutes), // Minute: from MonthlyJobTiming
		int(SayHelloDailyJobTiming.seconds), // Second: from MonthlyJobTiming
		0,                                   // Nanosecond: 0
		location,                            // Timezone: Asia/Kolkata
	)

	logger.Info(context.TODO(), "IST time check: ", jobTimeInKolkata)

	// Convert to UTC
	jobTimeInUTC := jobTimeInKolkata.UTC()
	logger.Info(context.TODO(), "UTC time check: ", jobTimeInUTC)
	cron.job, err = cron.scheduler.NewJob(
		gocron.DailyJob(
			SAY_HELLO_DAILY_CRON_JOB_INTERVAL_DAYS,
			gocron.NewAtTimes(
				gocron.NewAtTime(
					uint(jobTimeInUTC.Hour()),
					uint(jobTimeInUTC.Minute()),
					uint(jobTimeInUTC.Second()),
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
	for i := 0; i < 3; i++ {
		logger.Info(context.TODO(), "cron job attempt:", i+1)
		isSuccess, err := cron.appreciationService.UpdateAppreciation(ctx, "Asia/Kolkata")
		if err == nil && isSuccess {
			break
		}
		logger.Info(ctx, fmt.Sprintf("daily cron job err: %v ", err))
	}
}
