package reward

import (
	"context"
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	rewardRepo       repository.RewardStorer
	appreciationRepo repository.AppreciationStorer
}

type Service interface {
	GiveReward(ctx context.Context, rewardReq dto.Reward) (dto.Reward, error)
}

func NewService(rewardRepo repository.RewardStorer, appreciationRepo repository.AppreciationStorer) Service {
	return &service{
		rewardRepo:       rewardRepo,
		appreciationRepo: appreciationRepo,
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
		logger.Error("appreciationbyid: ",err.Error())
		return dto.Reward{}, err
	}
	fmt.Println("HIi")
	if appr.SenderId == sender {
		return dto.Reward{},apperrors.SelfAppreciationRewardError
	}

	if appr.ReceiverId == sender{
		return dto.Reward{},apperrors.SelfRewardError
	}

	userChk, err := rwrdSvc.rewardRepo.UserHasRewardQuota(ctx, nil, rewardReq.SenderId,rewardReq.Point)
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

	deduceChk, err := rwrdSvc.rewardRepo.DeduceRewardQuotaOfUser(ctx, tx, rewardReq.SenderId,int(rewardReq.Point))
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
	return reward, nil
}
