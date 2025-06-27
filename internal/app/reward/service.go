package reward

import (
	"context"
  "time"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	rewardRepo              repository.RewardStorer
	appreciationRepo        repository.AppreciationStorer
	reportedAppreciatonRepo repository.ReportAppreciationStorer
	userRepo                repository.UserStorer
}

type Service interface {
	GiveReward(ctx context.Context, rewardReq dto.Reward) (dto.Reward, error)
}

func NewService(rewardRepo repository.RewardStorer, appreciationRepo repository.AppreciationStorer, userRepo repository.UserStorer, reportedAppreciatonRepo repository.ReportAppreciationStorer) Service {
	return &service{
		rewardRepo:              rewardRepo,
		appreciationRepo:        appreciationRepo,
		userRepo:                userRepo,
		reportedAppreciatonRepo: reportedAppreciatonRepo,
	}
}

func IsAppreciationEligibleForReward(appreciationTime time.Time, now time.Time) bool {
	type Quarter struct {
		StartMonth time.Month
		Months     [3]time.Month
	}

	quarters := []Quarter{
		{time.March, [3]time.Month{time.March, time.April, time.May}},             
		{time.June, [3]time.Month{time.June, time.July, time.August}},             
		{time.September, [3]time.Month{time.September, time.October, time.November}}, 
		{time.December, [3]time.Month{time.December, time.January, time.February}},   
	}

	getQuarterIndex := func(month time.Month) int {
		for i, q := range quarters {
			for _, m := range q.Months {
				if m == month {
					return i
				}
			}
		}
		return -1
	}

	currQIndex := getQuarterIndex(now.Month())
	prevQIndex := (currQIndex + len(quarters) - 1) % len(quarters)

	currQ := quarters[currQIndex]
	prevQ := quarters[prevQIndex]

	for _, m := range currQ.Months {
		if appreciationTime.Month() == m {
			if m == time.January || m == time.February {
				if appreciationTime.Year() == now.Year() || appreciationTime.Year() == now.Year()-1 {
					return true
				}
			} else {
				if appreciationTime.Year() == now.Year() {
					return true
				}
			}
		}
	}

	lastMonthPrevQ := prevQ.Months[2] 
	firstMonthCurrQ := currQ.Months[0] 

	isLastMonthAppreciation := appreciationTime.Month() == lastMonthPrevQ
	isFirstMonthNow := now.Month() == firstMonthCurrQ

	isYearMatch := (appreciationTime.Year() == now.Year()-1 && now.Month() == time.January) ||
		(appreciationTime.Year() == now.Year())

	if isLastMonthAppreciation && isFirstMonthNow && isYearMatch {
		return true
	}

	return false
}


func (rwrdSvc *service) GiveReward(ctx context.Context, rewardReq dto.Reward) (dto.Reward, error) {

	logger.Debug(ctx, " rewardService: GiveReward: ", rewardReq)
	//add sender
	data := ctx.Value(constants.UserId)
	sender, ok := data.(int64)
	if !ok {
		logger.Error(ctx, "rewardService: err in parsing userid from token")
		return dto.Reward{}, apperrors.InternalServer
	}
	rewardReq.SenderId = sender

	appr, err := rwrdSvc.appreciationRepo.GetAppreciationById(ctx, nil, int32(rewardReq.AppreciationId))
	if err != nil {
		logger.Errorf(ctx, "rewardService: gerAppreciationById err : %v", err)
		return dto.Reward{}, err
	}
	logger.Debug(ctx, "rewardService: appr: ", appr)

	if appr.SenderID == sender {
		logger.Error(ctx, "rewardService: SelfAppreciationRewardError")
		return dto.Reward{}, apperrors.SelfAppreciationRewardError
	}

	if appr.ReceiverID == sender {
		logger.Error(ctx, "rewardService: SelfRewardError")
		return dto.Reward{}, apperrors.SelfRewardError
	}

	appreciationTime := time.UnixMilli(appr.CreatedAt)

	if !IsAppreciationEligibleForReward(appreciationTime, time.Now()) {
		return dto.Reward{}, apperrors.PreviousQuarterRatingNotAllowed
	}

	reportedAppr, err := rwrdSvc.reportedAppreciatonRepo.GetReportedAppreciationByAppreciationID(ctx, appr.ID)
	if err != nil && err != apperrors.InvalidId {
		logger.Errorf(ctx, "rewardService: GetReportedAppreciation: err: %v", err)
		return dto.Reward{}, err
	}
	if err == nil && reportedAppr.Status != "resolved" {
		return dto.Reward{}, apperrors.NotAllowedForReportedAppreciation
	}

	userChk, err := rwrdSvc.rewardRepo.UserHasRewardQuota(ctx, nil, rewardReq.SenderId, rewardReq.Point)
	if err != nil {
		logger.Errorf(ctx, "rewardService: UserHasRewardQuota: err: %v", err)
		return dto.Reward{}, err
	}
	logger.Debug(ctx, " userChk: ", userChk)

	if !userChk {
		logger.Error(ctx, "rewardService: RewardQuotaIsNotSufficient")
		return dto.Reward{}, apperrors.RewardQuotaIsNotSufficient
	}

	rwrdChk, err := rwrdSvc.rewardRepo.IsUserRewardForAppreciationPresent(ctx, nil, rewardReq.AppreciationId, rewardReq.SenderId)
	if err != nil {
		logger.Errorf(ctx, "rewardService: IsUserRewardForAppreciationPresent: err: %v", err)
		return dto.Reward{}, err
	}

	if rwrdChk {
		logger.Error(ctx, " rwrdChk: RewardAlreadyPresent")
		return dto.Reward{}, apperrors.RewardAlreadyPresent
	}

	//initializing database transaction
	tx, err := rwrdSvc.rewardRepo.BeginTx(ctx)
	if err != nil {
		logger.Error(ctx, "rewardService: error in BeginTx")
		return dto.Reward{}, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Infof(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := rwrdSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Infof(ctx, "error in creating transaction, err: %s", txErr.Error())
			return
		}
	}()
	repoRewardRes, err := rwrdSvc.rewardRepo.GiveReward(ctx, tx, rewardReq)
	if err != nil {
		logger.Errorf(ctx, "rewardService: GiveReward: err: %v", err)
		return dto.Reward{}, err
	}

	//deduce user rewardquota

	deduceChk, err := rwrdSvc.rewardRepo.DeduceRewardQuotaOfUser(ctx, tx, rewardReq.SenderId, int(rewardReq.Point))
	if err != nil {
		logger.Errorf(ctx, "rewardService: DeduceRewardQuotaOfUser: err: %v", err)
		return dto.Reward{}, err
	}

	if !deduceChk {
		logger.Error(ctx, "rewardService: error in reduce Reward Quota")
		return dto.Reward{}, apperrors.InternalServer
	}

	var reward dto.Reward
	reward.Id = repoRewardRes.Id
	reward.AppreciationId = repoRewardRes.AppreciationId
	reward.SenderId = repoRewardRes.SenderId
	reward.Point = repoRewardRes.Point
	quaterTimeStamp := user.GetQuarterStartUnixTime()

	req := dto.GetUserByIdReq{
		UserId:          sender,
		QuaterTimeStamp: quaterTimeStamp,
	}
	userInfo, err := rwrdSvc.userRepo.GetUserById(ctx, req)
	if err != nil {
		logger.Errorf(ctx, "rewardService: err in getting user data: %v", err)
	}
	rwrdSvc.sendRewardNotificationToSender(ctx, userInfo)
	rwrdSvc.sendRewardNotificationToReceiver(ctx, appr.ReceiverID)
	return reward, nil
}

func (rwrdSvc *service) sendRewardNotificationToSender(ctx context.Context, user dto.GetUserByIdResp) {

	logger.Debug(ctx, " rewardService: sendRewardNotificationToSender: user: ", user)
	notificationTokens, err := rwrdSvc.userRepo.ListDeviceTokensByUserID(ctx, user.UserId)
	logger.Debug(ctx, " notificationTokens: ", notificationTokens)
	if err != nil {
		logger.Errorf(ctx, "err in getting device tokens: %v", err)
		return
	}

	msg := notification.Message{
		Title: "Reward Given Successfully",
		Body:  "You have successfully given a reward! ",
	}

	logger.Debug(ctx, "msg:", msg, " notificationTokens: ", notificationTokens)
	for _, notificationToken := range notificationTokens {
		msg.SendNotificationToNotificationToken(notificationToken)
	}

}

func (rwrdSvc *service) sendRewardNotificationToReceiver(ctx context.Context, userID int64) {

	logger.Debug(ctx, " rewardService: sendRewardNotificationToReceiver")
	notificationTokens, err := rwrdSvc.userRepo.ListDeviceTokensByUserID(ctx, userID)
	logger.Debug(ctx, " notificationTokens: ", notificationTokens)
	if err != nil {
		logger.Errorf(ctx, " err in getting device tokens: %v", err)
		return
	}

	msg := notification.Message{
		Title: "Reward's incoming!",
		Body:  "You've been awarded a reward! Well done and keep up the JOSH!",
	}

	logger.Debug(ctx, " rewardService: msg: ", msg)
	for _, notificationToken := range notificationTokens {
		msg.SendNotificationToNotificationToken(notificationToken)
	}

}
