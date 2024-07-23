package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type rewardStore struct {
	BaseRepository
}

func NewRewardRepo(db *sqlx.DB) repository.RewardStorer {
	return &rewardStore{
		BaseRepository: BaseRepository{db},
	}
}

func (rwrd *rewardStore) GiveReward(ctx context.Context, tx repository.Transaction, reward dto.Reward) (repository.Reward, error) {

	queryExecutor := rwrd.InitiateQueryExecutor(tx)
	insertQuery, args, err := sq.
		Insert("rewards").
		Columns(constants.CreateRewardColumns...).
		Values(reward.AppreciationId, reward.Point, reward.SenderId).
		Suffix("RETURNING \"id\",\"appreciation_id\", \"point\",\"sender\",\"created_at\"").
		ToSql()

	if err != nil {
		logger.Error("err: ", "error in creating query", err.Error())
		return repository.Reward{}, apperrors.InternalServer
	}

	var rewardInfo repository.Reward
	err = queryExecutor.QueryRowx(insertQuery, args...).Scan(&rewardInfo.Id, &rewardInfo.AppreciationId, &rewardInfo.Point, &rewardInfo.SenderId, &rewardInfo.CreatedAt)
	if err != nil {
		logger.Error("Error executing create reward insert query: ", err)
		return repository.Reward{}, apperrors.InternalServer
	}

	return rewardInfo, nil

}

func (rwrd *rewardStore) IsUserRewardForAppreciationPresent(ctx context.Context, tx repository.Transaction, apprId int64, senderId int64) (bool, error) {
	// Initialize the Squirrel query builder

	fmt.Println("appr id: ", apprId)
	fmt.Println("sender: ", senderId)
	// Build the SQL query
	query, args, err := sq.Select("COUNT(*)").
		From("rewards").
		Where(squirrel.And{
			squirrel.Eq{"appreciation_id": apprId},
			squirrel.Eq{"sender": senderId},
		}).
		ToSql()
	if err != nil {
		logger.Error("err ", err.Error())
		return false, apperrors.InternalServer
	}

	fmt.Println("query: ", query)

	queryExecutor := rwrd.InitiateQueryExecutor(tx)

	var count int
	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).Scan(&count)
	if err != nil {		
		logger.Error("failed to execute query 1: ", err.Error())
		return false, apperrors.InternalServer
	}
	fmt.Println("count: ", count)
	// Check if user and appreciation id is present
	return count > 0, nil
}

func (rwrd *rewardStore) DeduceRewardQuotaOfUser(ctx context.Context, tx repository.Transaction, userId int64, points int) (bool, error) {
	queryExecutor := rwrd.InitiateQueryExecutor(tx)
	// Build the SQL query to update the reward_quota_balance
	updateQuery, args, err := sq.
		Update("users").
		Set("reward_quota_balance", squirrel.Expr("reward_quota_balance - ? * (SELECT points FROM grades WHERE id = users.grade_id)", points)).
		Where(squirrel.Eq{"id": userId}).
		ToSql()

	if err != nil {
		logger.Error("err: building SQL Query ", err.Error())
		return false, err
	}

	// Execute the query within the transaction context
	result, err := queryExecutor.Exec(updateQuery, args...)
	if err != nil {
		logger.Error("err: error executing SQL query:", err.Error())
		return false, err
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("err: error getting rows affected:", err)
		return false, err
	}

	// Return true if at least one row was updated, false otherwise
	return rowsAffected > 0, nil
}

func (rwrd *rewardStore) UserHasRewardQuota(ctx context.Context, tx repository.Transaction, userID int64, points int64) (bool, error) {
	// Initialize the Squirrel query builder
	// psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Build the SQL query
	query := `
		SELECT COUNT(*)
		FROM users u
		JOIN grades g ON u.grade_id = g.id
		WHERE u.id = $1
		AND u.reward_quota_balance >= $2 * g.points
	`

	// Arguments for the query
	args := []interface{}{userID, points}

	fmt.Println("id: ", userID, "points: ", points)
	fmt.Println("query: ", query)

	queryExecutor := rwrd.InitiateQueryExecutor(tx)

	var count int
	// Execute the query
	err := queryExecutor.QueryRowx(query, args...).Scan(&count)
	if err != nil {
		logger.Error("failed to execute query 2: ", err.Error())
		return false, apperrors.InternalServer
	}
	fmt.Println("count: ", count)
	// Check if user is present
	return count > 0, nil
}
