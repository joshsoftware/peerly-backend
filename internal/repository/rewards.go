package repository

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type RewardStorer interface {
	RepositoryTransaction

	GiveReward(ctx context.Context, tx Transaction, reward dto.Reward) (Reward, error)
	IsUserRewardForAppreciationPresent(ctx context.Context, tx Transaction, apprId int64, senderId int64) (bool, error)
	UserHasRewardQuota(ctx context.Context, tx Transaction, userID int64, points int64) (bool, error)
	DeduceRewardQuotaOfUser(ctx context.Context, tx Transaction, userId int64, points int) (bool, error)
}

type Reward struct {
	Id             int64 `db:"id"`
	AppreciationId int64 `db:"appreciation_id"`
	Point          int64 `db:"point"`
	SenderId       int64 `db:"sender"`
	CreatedAt      int64 `db:"created_at"`
}
