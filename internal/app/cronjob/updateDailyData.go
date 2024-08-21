package cronjob

import (
	"context"
	"time"
	"github.com/go-co-op/gocron/v2"
	apprSvc "github.com/joshsoftware/peerly-backend/internal/app/appreciation"
	logger "github.com/sirupsen/logrus"
)

const DAILY_JOB = "DAILY_JOB"
const SAY_HELLO_DAILY_CRON_JOB_INTERVAL_DAYS = 1

var SayHelloDailyJobTiming = JobTime{
	hours:   17,
	minutes: 34,
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
	// Load the location for Asia/Kolkata
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		logger.Warn(context.TODO(), "error loading location: %+v", err.Error())
		return
	}

	// Get the current date in Asia/Kolkata
	currentTimeInKolkata := time.Now().In(location)

	// Create a new time for today's date with MonthlyJobTiming hours, minutes, and seconds
	jobTimeInKolkata := time.Date(
		currentTimeInKolkata.Year(),   // Year: current year
		currentTimeInKolkata.Month(),  // Month: current month
		currentTimeInKolkata.Day(),    // Day: today's date
		int(MonthlyJobTiming.hours),   // Hour: from MonthlyJobTiming
		int(MonthlyJobTiming.minutes), // Minute: from MonthlyJobTiming
		int(MonthlyJobTiming.seconds), // Second: from MonthlyJobTiming
		0,                             // Nanosecond: 0
		location,                      // Timezone: Asia/Kolkata
	)

	logger.Info("IST time check: ",jobTimeInKolkata)

	// Convert to UTC
	jobTimeInUTC := jobTimeInKolkata.UTC()
	logger.Info("UTC time check: ",jobTimeInUTC)
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
		logger.Warn(context.TODO(), "error occurred while scheduling %s, message %+v", cron.name, err.Error())
	}
}

func (cron *DailyJob) Task(ctx context.Context) {
	logger.Info(ctx, "in daily job task")
	for i := 0; i < 3; i++ {
		logger.Info("cron job attempt:", i+1)
		isSuccess, err := cron.appreciationService.UpdateAppreciation(ctx)
		if err == nil && isSuccess {
			break
		}
	}

}
