package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type BadgeStorer interface {
	ListBadges(ctx context.Context) (badges []Badge, err error)
	EditBadge(ctx context.Context, reqData dto.UpdateBadgeReq) (err error)
}

type Badge struct {
	Id           int64         `db:"id"`
	Name         string        `db:"name"`
	RewardPoints int64         `db:"reward_points"`
	UpdatedBy    sql.NullInt64 `db:"updated_by"`
}
