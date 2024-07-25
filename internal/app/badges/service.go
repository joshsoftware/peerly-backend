package badges

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	badgeRepo repository.BadgesStorer
}

type Service interface {
	CreateBadge(ctx context.Context, badge dto.Badge) (dto.Badge, error)
	ListBadges(ctx context.Context) ([]dto.Badge, error)
	GetBadge(ctx context.Context, badgeID int8) (dto.Badge, error)
	DeleteBadge(ctx context.Context, badgeID int8) error
	UpdateBadge(ctx context.Context, badgeUpdateInfo dto.Badge) (dto.Badge, error)
}

func NewService(badgesRepo repository.BadgesStorer) Service {
	return &service{
		badgeRepo: badgesRepo,
	}
}

func (bdg *service) CreateBadge(ctx context.Context, badge dto.Badge) (dto.Badge, error) {

	//check whether badge name is already present
	badgeId := bdg.badgeRepo.GetBadgeByName(ctx, nil, badge.Name)
	if badgeId != 0 {
		return dto.Badge{}, apperrors.BadgeNameAlreadyExists
	}

	//check whether badge reward points is already present
	badgeId = bdg.badgeRepo.GetBadgeByRewardPoints(ctx, nil, badge.RewardPoints)
	if badgeId != 0 {
		return dto.Badge{}, apperrors.BadgeRewardPointsAlreadyExists
	}

	resp, err := bdg.badgeRepo.CreateBadge(ctx, nil, badge)
	if err != nil {
		return dto.Badge{}, err
	}
	return dto.Badge{
		ID:           resp.ID,
		Name:         resp.Name,
		RewardPoints: resp.RewardPoints,
	}, nil
}

func (bdg *service) ListBadges(ctx context.Context) ([]dto.Badge, error) {

	dbResp, err := bdg.badgeRepo.ListBadges(ctx, nil)
	if err != nil {
		return []dto.Badge{}, err
	}

	resp := make([]dto.Badge, 0)
	for _, value := range dbResp {
		coreValue := mapRepoBadgeToDTOBadge(value)
		resp = append(resp, coreValue)
	}
	return resp, nil
}

func (bdg *service) GetBadge(ctx context.Context, badgeID int8) (dto.Badge, error) {
	resp, err := bdg.badgeRepo.GetBadge(ctx, nil, badgeID)
	if err != nil {
		return dto.Badge{}, err
	}

	return mapRepoBadgeToDTOBadge(resp), nil
}

func (bdg *service) DeleteBadge(ctx context.Context, badgeID int8) error {
	return bdg.badgeRepo.DeleteBadge(ctx, nil, badgeID)
}

func (bdg *service) UpdateBadge(ctx context.Context, badgeUpdateInfo dto.Badge) (dto.Badge, error) {

	//check whether badge is present or not
	_,err := bdg.badgeRepo.GetBadge(ctx,nil,badgeUpdateInfo.ID)
	if err != nil{
		return dto.Badge{},err
	}

	//check whether badge name is already present
	badgeId := bdg.badgeRepo.GetBadgeByName(ctx, nil, badgeUpdateInfo.Name)
	if badgeId != 0 && badgeUpdateInfo.ID != badgeId {
		return dto.Badge{}, apperrors.BadgeNameAlreadyExists
	}

	//check whether badge reward points is already present
	badgeId = bdg.badgeRepo.GetBadgeByRewardPoints(ctx, nil, badgeUpdateInfo.RewardPoints)
	if badgeId != 0 && badgeUpdateInfo.ID != badgeId{
		return dto.Badge{}, apperrors.BadgeRewardPointsAlreadyExists
	}

	resp, err := bdg.badgeRepo.UpdateBadge(ctx, nil, badgeUpdateInfo)
	if err != nil {
		logger.Errorf("err in updating badge: %v",err)
		return dto.Badge{}, err
	}

	return mapRepoBadgeToDTOBadge(resp), nil
}
