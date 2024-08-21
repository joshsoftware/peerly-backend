package cronjob

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"time"

	logger "github.com/sirupsen/logrus"
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

func NewMontlyJob(userSvc user.Service, scheduler gocron.Scheduler) Job {
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

	// Convert to UTC
	jobTimeInUTC := jobTimeInKolkata.UTC()

	cron.job, err = cron.scheduler.NewJob(
		gocron.MonthlyJob(MONTHLY_CRON_JOB_INTERVAL_MONTHS, gocron.NewDaysOfTheMonth(-1),
			gocron.NewAtTimes(
				gocron.NewAtTime(
					uint(jobTimeInUTC.Hour()),
					uint(jobTimeInUTC.Minute()),
					uint(jobTimeInUTC.Second()),
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
			sendRewardQuotaRefilledNotificationToAll()
			break
		}
	}

}

func sendRewardQuotaRefilledNotificationToAll() {
	msg := notification.Message{
		Title: "🚀 Reward Quota Reset! ",
		Body:  "Quota for Rewards Renewed. Time to Shower Your Peers with Kudos! 🎁",
	}
	msg.SendNotificationToTopic("peerly")
}
