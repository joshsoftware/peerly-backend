package badges

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	badgesRepo repository.BadgeStorer
}

type Service interface {
	ListBadges(ctx context.Context) (resp []dto.Badge, err error)
	EditBadge(ctx context.Context, id string, rewardPoints int64) (err error)
}

func NewService(badgesRepo repository.BadgeStorer) Service {
	return &service{
		badgesRepo: badgesRepo,
	}
}

func (bs *service) ListBadges(ctx context.Context) (resp []dto.Badge, err error) {

	dbResp, err := bs.badgesRepo.ListBadges(ctx)
	if err != nil {
		logger.Error(err.Error())
		err = apperrors.InternalServerError
	}

	for _, item := range dbResp {
		svcItem := mapDbToSvc(item)
		resp = append(resp, svcItem)
	}

	return

}

func (gs *service) EditBadge(ctx context.Context, id string, rewardPoints int64) (err error) {
	badgeId, err := utils.VarsStringToInt(id, "badgeId")
	if err != nil {
		return
	}
	var reqData dto.UpdateBadgeReq
	reqData.Id = badgeId
	if rewardPoints < 0 {
		logger.Errorf("badge reward points cannot be negative, reward points: %d", rewardPoints)
		err = apperrors.NegativeBadgePoints
		return
	}
	reqData.RewardPoints = rewardPoints
	err = gs.badgesRepo.EditBadge(ctx, reqData)
	if err != nil {
		logger.Error(err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}

func mapDbToSvc(dbResp repository.Badge) (svcResp dto.Badge) {
	svcResp.Id = dbResp.Id
	svcResp.Name = dbResp.Name
	svcResp.RewardPoints = dbResp.RewardPoints
	return
}
