package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type badgeStore struct {
	DB         *sqlx.DB
	BadgeTable string
}

func NewBadgeRepo(db *sqlx.DB) repository.BadgeStorer {
	return &badgeStore{
		DB:         db,
		BadgeTable: constants.BadgeTable,
	}
}

var BadgeColumns = []string{"id", "name", "reward_points", "updated_by"}

func (bs *badgeStore) ListBadges(ctx context.Context) (badges []repository.Badge, err error) {
	queryBuilder := repository.Sq.Select(BadgeColumns...).From(bs.BadgeTable).OrderBy("id")
	listBadgesQuery, _, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}
	err = bs.DB.SelectContext(
		ctx,
		&badges,
		listBadgesQuery,
	)

	if err != nil {
		err = fmt.Errorf("error while getting badges, err: %w", err)
		return
	}

	return
}

func (bs *badgeStore) EditBadge(ctx context.Context, reqData dto.UpdateBadgeReq) (err error) {
	queryBuilder := repository.Sq.Update(bs.BadgeTable).Set("reward_points", reqData.RewardPoints).Set("updated_by", reqData.UserId).Where(squirrel.Eq{"id": reqData.Id})
	updateBadgeQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}
	_, err = bs.DB.ExecContext(ctx, updateBadgeQuery, args...)
	if err != nil {
		err = fmt.Errorf("error in updating badge points, err: %w", err)
		return
	}

	return
}
