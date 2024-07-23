package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

var AppreciationColumns = []string{"id", "core_value_id", "description", "quarter", "sender", "receiver"}

type appreciationsStore struct {
	BaseRepository
	TableAppreciations string
	TableRewards       string
	TableUsers         string
	TableCoreValues    string
}

func NewAppreciationRepo(db *sqlx.DB) repository.AppreciationStorer {
	return &appreciationsStore{
		BaseRepository:     BaseRepository{db},
		TableAppreciations: "appreciations",
		TableRewards:       "rewards",
		TableUsers:         "users",
		TableCoreValues:    "core_values",
	}
}

func (appr *appreciationsStore) CreateAppreciation(ctx context.Context, tx repository.Transaction, appreciation dto.Appreciation) (repository.Appreciation, error) {

	queryExecutor := appr.InitiateQueryExecutor(tx)

	insertQuery, args, err := sq.
		Insert(appr.TableAppreciations).Columns(AppreciationColumns[1:]...).
		Values(appreciation.CoreValueID, appreciation.Description, appreciation.Quarter, appreciation.Sender, appreciation.Receiver).
		Suffix("RETURNING id,core_value_id, description,total_reward_points,quarter,sender,receiver,created_at,updated_at").
		ToSql()
	if err != nil {
		logger.Error(fmt.Sprintf("error in generating squirrel query, err: %v", err))
		return repository.Appreciation{}, apperrors.InternalServerError
	}

	var resAppr repository.Appreciation
	err = queryExecutor.QueryRowx(insertQuery, args...).StructScan(&resAppr)
	if err != nil {
		logger.Error(fmt.Sprintf("Error executing create appreciation insert query: %v", err))
		return repository.Appreciation{}, apperrors.InternalServer
	}

	return resAppr, nil
}

func (appr *appreciationsStore) GetAppreciationById(ctx context.Context, tx repository.Transaction, apprId int32) (repository.AppreciationInfo, error) {

	queryExecutor := appr.InitiateQueryExecutor(tx)

	// Get logged-in user ID
	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userID from token")
		return repository.AppreciationInfo{}, apperrors.InternalServer
	}

	// Build the SQL query
	query, args, err := sq.Select(
		"a.id",
		"cv.name AS core_value_name",
		"a.description",
		"a.is_valid",
		"a.total_reward_points",
		"a.quarter",
		"u_sender.id AS sender_id",
		"u_sender.first_name AS sender_first_name",
		"u_sender.last_name AS sender_last_name",
		"u_sender.profile_image_url AS sender_image_url",
		"u_sender.designation AS sender_designation",
		"u_receiver.id AS receiver_id",
		"u_receiver.first_name AS receiver_first_name",
		"u_receiver.last_name AS receiver_last_name",
		"u_receiver.profile_image_url AS receiver_image_url",
		"u_receiver.designation AS receiver_designation",
		"a.created_at",
		"a.updated_at",
		"COUNT(r.id) AS total_rewards",
		fmt.Sprintf(
			`COALESCE((
				SELECT r2.point 
				FROM rewards r2 
				WHERE r2.appreciation_id = a.id AND r2.sender = %d
			), 0) AS given_reward_point`, userID),
	).From(appr.TableAppreciations+" a").
		LeftJoin(appr.TableUsers+" u_sender ON a.sender = u_sender.id").
		LeftJoin(appr.TableUsers+" u_receiver ON a.receiver = u_receiver.id").
		LeftJoin(appr.TableCoreValues+" cv ON a.core_value_id = cv.id").
		LeftJoin(appr.TableRewards+" r ON a.id = r.appreciation_id").
		Where(squirrel.And{
			squirrel.Eq{"a.id": apprId},
			squirrel.Eq{"a.is_valid": true},
		}).
		GroupBy("a.id","cv.name","u_sender.id","u_receiver.id").
		ToSql()

	if err != nil {
		logger.Error(fmt.Sprintf("error in generating squirrel query, err: %v", err))
		return repository.AppreciationInfo{}, apperrors.InternalServer
	}

	var resAppr repository.AppreciationInfo

	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).StructScan(&resAppr)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(fmt.Sprintf("no appreciation found with id: %d", apprId))
			return repository.AppreciationInfo{}, apperrors.AppreciationNotFound
		}
		logger.Error(fmt.Sprintf("failed to execute query: %v", err))
		return repository.AppreciationInfo{}, apperrors.InternalServer
	}
	return resAppr, nil
}
func (appr *appreciationsStore) GetAppreciations(ctx context.Context, tx repository.Transaction, filter dto.AppreciationFilter) ([]repository.AppreciationInfo, repository.Pagination, error) {

	queryExecutor := appr.InitiateQueryExecutor(tx)

	// Get logged-in user ID
	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userID from token")
		return []repository.AppreciationInfo{}, repository.Pagination{}, apperrors.InternalServerError
	}

	// query builder for counting total records
	countQueryBuilder := sq.Select("COUNT(*)").
		From("appreciations a").
		LeftJoin("users u_sender ON a.sender = u_sender.id").
		LeftJoin("users u_receiver ON a.receiver = u_receiver.id").
		LeftJoin("core_values cv ON a.core_value_id = cv.id").
		Where(squirrel.Eq{"a.is_valid": true})
	if filter.Name != "" {
		lowerNameFilter := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Name))
		countQueryBuilder = countQueryBuilder.Where(
			"(LOWER(CONCAT(u_sender.first_name, ' ', u_sender.last_name)) LIKE ? OR "+
				"LOWER(CONCAT(u_receiver.first_name, ' ', u_receiver.last_name)) LIKE ?)",
			lowerNameFilter, lowerNameFilter,
		)
	}
	
	if filter.Self {
		countQueryBuilder = countQueryBuilder.Where(squirrel.Or{
			squirrel.Eq{"a.sender": userID},
			squirrel.Eq{"a.receiver": userID},
		})
	}

	countSql, countArgs, err := countQueryBuilder.ToSql()
	if err != nil {
		logger.Error("failed to build count query: ", err.Error())
		return []repository.AppreciationInfo{}, repository.Pagination{}, apperrors.InternalServerError
	}

	var totalRecords int32
	err = queryExecutor.QueryRowx(countSql, countArgs...).Scan(&totalRecords)
	if err != nil {
		logger.Error("failed to execute count query: ", err.Error())
		return []repository.AppreciationInfo{}, repository.Pagination{}, apperrors.InternalServerError
	}

	pagination := getPaginationMetaData(filter.Page, filter.Limit, totalRecords)
	// Initialize the Squirrel query builder
	queryBuilder := sq.Select(
		"a.id",
		"cv.name AS core_value_name",
		"cv.description AS core_value_description",
		"a.description",
		"a.is_valid",
		"a.total_reward_points",
		"a.quarter",
		"u_sender.id AS sender_id",
		"u_sender.first_name AS sender_first_name",
		"u_sender.last_name AS sender_last_name",
		"u_sender.profile_image_url AS sender_image_url",
		"u_sender.designation AS sender_designation",
		"u_receiver.id AS receiver_id",
		"u_receiver.first_name AS receiver_first_name",
		"u_receiver.last_name AS receiver_last_name",
		"u_receiver.profile_image_url AS receiver_image_url",
		"u_receiver.designation AS receiver_designation",
		"a.created_at",
		"a.updated_at",
		"COUNT(r.id) AS total_rewards",
		fmt.Sprintf(
			`COALESCE((
				SELECT r2.point 
				FROM rewards r2 
				WHERE r2.appreciation_id = a.id AND r2.sender = %d
			), 0) AS given_reward_point`, userID),
	).From("appreciations a").
		LeftJoin("users u_sender ON a.sender = u_sender.id").
		LeftJoin("users u_receiver ON a.receiver = u_receiver.id").
		LeftJoin("core_values cv ON a.core_value_id = cv.id").
		LeftJoin("rewards r ON a.id = r.appreciation_id").
		Where(squirrel.Eq{"a.is_valid": true}).
		GroupBy("a.id","cv.name","cv.description","u_sender.id","u_receiver.id")

	if filter.Name != "" {
		lowerNameFilter := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Name))
		queryBuilder = queryBuilder.Where(
			"(LOWER(CONCAT(u_sender.first_name, ' ', u_sender.last_name)) LIKE ? OR "+
				"LOWER(CONCAT(u_receiver.first_name, ' ', u_receiver.last_name)) LIKE ?)",
			lowerNameFilter, lowerNameFilter,
		)
	}

	if filter.Self {
		queryBuilder = queryBuilder.Where(squirrel.Or{
			squirrel.Eq{"a.sender": userID},
			squirrel.Eq{"a.receiver": userID},
		})
	}


	if filter.SortOrder != "" {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("a.created_at %s", filter.SortOrder))
	}

	offset := (filter.Page - 1) * filter.Limit

	// Add pagination
	queryBuilder = queryBuilder.Limit(uint64(filter.Limit)).Offset(uint64(offset))
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Error("failed to build query: ", err.Error())
		return nil, repository.Pagination{}, apperrors.InternalServerError
	}

	queryExecutor = appr.InitiateQueryExecutor(tx)
	res := make([]repository.AppreciationInfo, 0)
	err = sqlx.Select(queryExecutor, &res, sql, args...)
	if err != nil {
		logger.Error("failed to execute query: ", err.Error())
		return nil, repository.Pagination{}, apperrors.InternalServerError
	}

	return res, pagination, nil
}

func (appr *appreciationsStore) DeleteAppreciation(ctx context.Context, tx repository.Transaction, apprId int32) (bool, error) {
	query, args, err := sq.Update("appreciations").
		Set("is_valid", false).
		Where(squirrel.And{
			squirrel.Eq{"id": apprId},
			squirrel.Eq{"is_valid": true},
		}).
		ToSql()

	if err != nil {
		logger.Error("Error building SQL: ", err.Error())
		return false, apperrors.InternalServer
	}

	queryExecutor := appr.InitiateQueryExecutor(tx)

	result, err := queryExecutor.Exec(query, args...)
	if err != nil {
		logger.Error("Error executing SQL: ", err.Error())
		return false, apperrors.InternalServer
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected: ", err.Error())
		return false, apperrors.InternalServer
	}

	if rowsAffected == 0 {
		logger.Warn("No rows affected")
		return false, apperrors.AppreciationNotFound
	}

	return true, nil
}

func (appr *appreciationsStore) IsUserPresent(ctx context.Context, tx repository.Transaction, userID int64) (bool, error) {

	// Build the SQL query
	query, args, err := sq.Select("COUNT(*)").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	if err != nil {
		logger.Error("err ", err.Error())
		return false, apperrors.InternalServer
	}

	queryExecutor := appr.InitiateQueryExecutor(tx)

	var count int
	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).Scan(&count)
	if err != nil {
		logger.Error("failed to execute query: ", err.Error())
		return false, apperrors.InternalServer
	}

	// Check if user is present
	return count > 0, nil
}
