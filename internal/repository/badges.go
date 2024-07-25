package repository

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type BadgesStorer interface {
	CreateBadge(ctx context.Context, tx Transaction, badge dto.Badge) (Badge, error)
	ListBadges(ctx context.Context, tx Transaction) ([]Badge, error)
	GetBadge(ctx context.Context, tx Transaction,badgeID int8) (Badge, error)
	DeleteBadge(ctx context.Context, tx Transaction,badgeID int8) error
	UpdateBadge(ctx context.Context, tx Transaction, badgeUpdateInfo dto.Badge) (Badge, error)
	GetBadgeByName(ctx context.Context,tx Transaction,badgeName string) int8
	GetBadgeByRewardPoints(ctx context.Context,tx Transaction,badgeRewardPoints int16) int8
}

// Badge - struct representing a badge object
type Badge struct {
	ID           int8   `db:"id"`
	Name         string `db:"name"`
	RewardPoints int16  `db:"reward_points"`
}
