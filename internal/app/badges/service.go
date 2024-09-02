package badges

import (
	"context"
	"fmt"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	badgesRepo repository.BadgeStorer
	userRepo   repository.UserStorer
}

type Service interface {
	ListBadges(ctx context.Context) (resp []dto.Badge, err error)
	EditBadge(ctx context.Context, id string, rewardPoints int64) (err error)
}

func NewService(badgesRepo repository.BadgeStorer, userRepo repository.UserStorer) Service {
	return &service{
		badgesRepo: badgesRepo,
		userRepo:   userRepo,
	}
}

func (bs *service) ListBadges(ctx context.Context) ([]dto.Badge, error) {

	var resp []dto.Badge

	fmt.Println("In list badges")

	dbResp, err := bs.badgesRepo.ListBadges(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return nil, err
	}

	for _, item := range dbResp {
		fmt.Println("updated by: ", item.UpdatedBy)
		if item.UpdatedBy.Valid {
			reqData := dto.GetUserByIdReq{
				UserId:          item.UpdatedBy.Int64,
				QuaterTimeStamp: GetQuarterStartUnixTime(),
			}
			user, err := bs.userRepo.GetUserById(ctx, reqData)
			if err != nil {
				return nil, err
			}
			svcItem := mapDbToSvc(item, user)
			resp = append(resp, svcItem)
		} else {
			svcItem := mapDbToSvc(item, dto.GetUserByIdResp{})
			resp = append(resp, svcItem)
		}
	}

	return resp, nil

}

func (gs *service) EditBadge(ctx context.Context, id string, rewardPoints int64) (err error) {
	badgeId, err := utils.VarsStringToInt(id, "badgeId")
	if err != nil {
		return
	}

	var reqData dto.UpdateBadgeReq
	reqData.Id = badgeId
	if rewardPoints < 0 {
		logger.Errorf(ctx, "badge reward points cannot be negative, reward points: %d", rewardPoints)
		err = apperrors.NegativeBadgePoints
		return
	}
	reqData.RewardPoints = rewardPoints
	userId := ctx.Value(constants.UserId)
	data, ok := userId.(int64)
	if !ok {
		logger.Error(context.Background(),"Error in typecasting user id")
		err = apperrors.InternalServerError
		return
	}
	reqData.UserId = data
	err = gs.badgesRepo.EditBadge(ctx, reqData)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}

func mapDbToSvc(dbResp repository.Badge, user dto.GetUserByIdResp) (svcResp dto.Badge) {
	svcResp.Id = dbResp.Id
	svcResp.Name = dbResp.Name
	svcResp.RewardPoints = dbResp.RewardPoints
	svcResp.UpdatedBy = user.FirstName + " " + user.LastName
	return
}

func GetQuarterStartUnixTime() int64 {
	// Example function to get the Unix timestamp of the start of the quarter
	now := time.Now()
	quarterStart := time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
	return quarterStart.Unix() * 1000 // convert to milliseconds
}
