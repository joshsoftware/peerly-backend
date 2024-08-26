package cronjob

import (
	"context"
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	orgSvc "github.com/joshsoftware/peerly-backend/internal/app/organizationConfig"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	log "github.com/joshsoftware/peerly-backend/internal/pkg/logger"

	logger "github.com/sirupsen/logrus"
)

const MONTHLY_JOB = "MONTHLY_JOB"
var MONTHLY_CRON_JOB_INTERVAL_MONTHS = 1

var MonthlyJobTiming = JobTime{
	hours:   23,
	minutes: 59,
	seconds: 0,
}

type MonthlyJob struct {
	CronJob
	userService user.Service
	organizationConfigService orgSvc.Service
}

func NewMontlyJob(userSvc user.Service,organizationConfigService orgSvc.Service,scheduler gocron.Scheduler) Job {
	return &MonthlyJob{
		userService: userSvc,
		organizationConfigService: organizationConfigService,
		CronJob: CronJob{
			name:      MONTHLY_JOB,
			scheduler: scheduler,
		},
	}
}

func (cron *MonthlyJob) Schedule() error{
	var err error
	err = cron.setMonthlyInterval()
	if err != nil{
		return err
	}
	cron.job, err = cron.scheduler.NewJob(
		gocron.MonthlyJob(uint(MONTHLY_CRON_JOB_INTERVAL_MONTHS), gocron.NewDaysOfTheMonth(-1),
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
		logger.Warn(context.TODO(), fmt.Sprintf("error occurred while scheduling %s, message %+v", cron.name, err.Error()))
		return err
	}
	return nil
}

func (cron *MonthlyJob) Task(ctx context.Context)  {
	logger.Info(ctx, "in monthly job task")
	var err error
	for i := 0; i < 3; i++ {
		logger.Info("cron job attempt:", i+1)
		err = cron.userService.UpdateRewardQuota(ctx)
		if err == nil {
			sendRewardQuotaRefilledNotificationToAll()
			return 
		}
		log.Info(fmt.Sprintf("cronjob fail error: %v",err.Error()))
	}
	// return err
	return
}

func sendRewardQuotaRefilledNotificationToAll() {
	msg := notification.Message{
		Title: "Reward Quota is Refilled",
		Body:  "Your reward quota is reset! You now recognize your colleagues.",
	}
	msg.SendNotificationToTopic("peerly")
}

func (cron *MonthlyJob) setMonthlyInterval()error{
	orgInfo, err := cron.organizationConfigService.GetOrganizationConfig(context.Background())
	if err != nil{
		return err
	}
	MONTHLY_CRON_JOB_INTERVAL_MONTHS = orgInfo.RewardQuotaRenewalFrequency;
	log.Info(fmt.Sprintf("MONTHLY_CRON_JOB_INTERVAL_MONTHS = %d",MONTHLY_CRON_JOB_INTERVAL_MONTHS))
	return nil
}