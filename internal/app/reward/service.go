package reward

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	rewardRepo       repository.RewardStorer
	appreciationRepo repository.AppreciationStorer
	userRepo         repository.UserStorer
}

type Service interface {
	GiveReward(ctx context.Context, rewardReq dto.Reward) (dto.Reward, error)
}

func NewService(rewardRepo repository.RewardStorer, appreciationRepo repository.AppreciationStorer, userRepo repository.UserStorer) Service {
	return &service{
		rewardRepo:       rewardRepo,
		appreciationRepo: appreciationRepo,
		userRepo:         userRepo,
	}
}

func (rwrdSvc *service) GiveReward(ctx context.Context, rewardReq dto.Reward) (dto.Reward, error) {

	//add sender
	data := ctx.Value(constants.UserId)
	sender, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userid from token")
		return dto.Reward{}, apperrors.InternalServer
	}
	rewardReq.SenderId = sender

	appr, err := rwrdSvc.appreciationRepo.GetAppreciationById(ctx, nil, int32(rewardReq.AppreciationId))
	if err != nil {
		return dto.Reward{}, err
	}

	if appr.SenderID == sender {
		return dto.Reward{}, apperrors.SelfAppreciationRewardError
	}

	if appr.ReceiverID == sender {
		return dto.Reward{}, apperrors.SelfRewardError
	}

	userChk, err := rwrdSvc.rewardRepo.UserHasRewardQuota(ctx, nil, rewardReq.SenderId, rewardReq.Point)
	if err != nil {
		return dto.Reward{}, err
	}

	if !userChk {
		return dto.Reward{}, apperrors.RewardQuotaIsNotSufficient
	}

	rwrdChk, err := rwrdSvc.rewardRepo.IsUserRewardForAppreciationPresent(ctx, nil, rewardReq.AppreciationId, rewardReq.SenderId)
	if err != nil {
		return dto.Reward{}, err
	}

	if rwrdChk {
		return dto.Reward{}, apperrors.RewardAlreadyPresent
	}

	//initializing database transaction
	tx, err := rwrdSvc.rewardRepo.BeginTx(ctx)
	if err != nil {
		return dto.Reward{}, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Info(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := rwrdSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Info(ctx, "error in creating transaction, err: %s", txErr.Error())
			return
		}
	}()
	repoRewardRes, err := rwrdSvc.rewardRepo.GiveReward(ctx, tx, rewardReq)
	if err != nil {
		return dto.Reward{}, err
	}

	//deduce user rewardquota

	deduceChk, err := rwrdSvc.rewardRepo.DeduceRewardQuotaOfUser(ctx, tx, rewardReq.SenderId, int(rewardReq.Point))
	if err != nil {
		return dto.Reward{}, err
	}

	if !deduceChk {
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
	userInfo ,err := rwrdSvc.userRepo.GetUserById(ctx,req)
	if err != nil {
		logger.Errorf("err in getting user data: %v",err)
	}
	rwrdSvc.sendRewardNotificationToSender(ctx,userInfo)
	rwrdSvc.sendRewardNotificationToReceiver(ctx,appr.ReceiverID)
	return reward, nil
}

func GetQuarterStartUnixTime() {
	panic("unimplemented")
}

func (rwrdSvc *service) sendRewardNotificationToSender(ctx context.Context, user dto.GetUserByIdResp) {

	notificationTokens, err := rwrdSvc.userRepo.ListDeviceTokensByUserID(ctx, user.UserId)
	if err != nil {
		logger.Errorf("err in getting device tokens: %v", err)
		return 
	}

	msg := notification.Message{
		Title: "Reward Given Successfully",
		Body:  "You have successfully given a reward! ",
	}

	for _, notificationToken := range notificationTokens {
		msg.SendNotificationToNotificationToken(notificationToken)
	}

}

func (rwrdSvc *service) sendRewardNotificationToReceiver(ctx context.Context, userID int64) {

	notificationTokens, err := rwrdSvc.userRepo.ListDeviceTokensByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("err in getting device tokens: %v", err)
		return 
	}

	msg := notification.Message{
		Title: "Reward's incoming!",
		Body:  "You've been awarded a reward! Well done and keep up the JOSH!",
	}

	for _, notificationToken := range notificationTokens {
		msg.SendNotificationToNotificationToken(notificationToken)
	}

}
