package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

var BadgeColumns = []string{"id", "name", "reward_points"}

type badgesStore struct {
	BaseRepository
	TableBadge string
}

func NewBadgeRepo(db *sqlx.DB) repository.BadgesStorer {
	return &badgesStore{
		BaseRepository: BaseRepository{db},
		TableBadge:     "badges",
	}
}

func (bdg *badgesStore) CreateBadge(ctx context.Context, tx repository.Transaction, badge dto.Badge) (repository.Badge, error) {

	// queryExecutor := bdg.InitiateQueryExecutor(tx)

	queryBuilder := repository.Sq.Insert(bdg.TableBadge).Columns(BadgeColumns[1], BadgeColumns[2]).Values(badge.Name, badge.RewardPoints).Suffix("RETURNING id, name, reward_points")

	createCoreValueQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Error(fmt.Sprintf("error in generating squirrel query, err: %v", err))
		return repository.Badge{}, err
	}

	// fmt.Println("query: ",createCoreValueQuery)
	var res repository.Badge
	err = bdg.DB.GetContext(
		ctx,
		&res,
		createCoreValueQuery,
		args...,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("error while creating core value, err: %v", err))
		return repository.Badge{}, err
	}

	return res, nil
}

func (bdg *badgesStore) ListBadges(ctx context.Context, tx repository.Transaction) ([]repository.Badge, error) {

	queryExecutor := bdg.InitiateQueryExecutor(tx)

	query, args, err := repository.Sq.Select("*").
		From(bdg.TableBadge).
		ToSql()

	if err != nil {
		logger.Error("err ", err.Error())
		return []repository.Badge{}, apperrors.InternalServer
	}

	res := make([]repository.Badge, 0)

	err = sqlx.Select(queryExecutor, &res, query, args...)
	if err != nil {
		logger.Error("failed to execute query: ", err.Error())
		return []repository.Badge{}, apperrors.InternalServer
	}

	return res, nil
}

func (bdg *badgesStore) GetBadge(ctx context.Context, tx repository.Transaction, badgeID int8) (repository.Badge, error) {

	queryExecutor := bdg.InitiateQueryExecutor(tx)

	query, args, err := repository.Sq.Select("*").
		From(bdg.TableBadge).
		Where(squirrel.Eq{"id": badgeID}).
		ToSql()

	if err != nil {
		logger.Error("err ", err.Error())
		return repository.Badge{}, apperrors.InternalServer
	}

	var res repository.Badge

	err = queryExecutor.QueryRowx(query, args...).StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(fmt.Sprintf("no appreciation found with id: %d", badgeID))
			return repository.Badge{}, apperrors.BadgeNotFound
		}
		logger.Error(fmt.Sprintf("failed to execute query: %v", err))
		return repository.Badge{}, apperrors.InternalServer
	}

	return res, nil
}

func (bdg *badgesStore) DeleteBadge(ctx context.Context, tx repository.Transaction, badgeID int8) error {

	queryExecutor := bdg.InitiateQueryExecutor(tx)

	query, args, err := repository.Sq.Delete(bdg.TableBadge).
		Where(squirrel.Eq{"id": badgeID}).
		ToSql()

	if err != nil {
		logger.Error("err ", err.Error())
		return apperrors.InternalServer
	}

	qres, err := queryExecutor.Exec(query, args...)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to execute query: %v", err))
		return apperrors.InternalServer
	}

	noOfRowsAffected, err := qres.RowsAffected()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to execute query: %v", err))
		return apperrors.InternalServer
	}

	if noOfRowsAffected == 0 {
		logger.Error("no row affected: " )
		return apperrors.BadgeNotFound
	}

	return nil
}

func (bdg *badgesStore) UpdateBadge(ctx context.Context, tx repository.Transaction, badgeUpdateInfo dto.Badge)(updatedBadge repository.Badge, err error){

	queryExecutor := bdg.InitiateQueryExecutor(tx)

	queryBuilder := repository.Sq.Update(bdg.TableBadge).
		Where(squirrel.Eq{"id": badgeUpdateInfo.ID}).
		Suffix("RETURNING id,name,reward_points")

	if badgeUpdateInfo.Name != "" {
		queryBuilder = queryBuilder.Set("name",badgeUpdateInfo.Name)
	}

	if badgeUpdateInfo.RewardPoints != 0 {
		queryBuilder = queryBuilder.Set("reward_points",badgeUpdateInfo.RewardPoints)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf("Error building update query : %v",err)
		return repository.Badge{},err
	}

	err = queryExecutor.QueryRowx(query, args...).StructScan(&updatedBadge)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("Just updated badge, but can't find it! %v", err)
			return repository.Badge{}, apperrors.BadgeNotFound
		}
	}
	return
}

func (bdg *badgesStore)	GetBadgeByName(ctx context.Context,tx repository.Transaction,badgeName string) int8{

	queryExecutor := bdg.InitiateQueryExecutor(tx)

	query, args, err := repository.Sq.Select("id").
		From(bdg.TableBadge).
		Where(squirrel.Eq{"name": badgeName}).
		ToSql()

	if err != nil {
		logger.Error("err ", err.Error())
		return 0
	}

	var badgeId int8
	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).Scan(&badgeId)
	if err != nil {
		logger.Error("failed to execute query: ", err.Error())
		return 0
	}

	return badgeId

}
func (bdg *badgesStore) GetBadgeByRewardPoints(ctx context.Context,tx repository.Transaction,badgeRewardPoints int16) int8{

	queryExecutor := bdg.InitiateQueryExecutor(tx)

	query, args, err := repository.Sq.Select("id").
		From(bdg.TableBadge).
		Where(squirrel.Eq{"reward_points": badgeRewardPoints}).
		ToSql()

	if err != nil {
		logger.Error("err ", err.Error())
		return 0
	}

	var badgeId int8
	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).Scan(&badgeId)
	if err != nil {
		logger.Error("failed to execute query: ", err.Error())
		return 0
	}

	return badgeId
}