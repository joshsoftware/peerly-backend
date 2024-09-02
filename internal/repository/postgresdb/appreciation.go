package repository

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

var AppreciationColumns = []string{"id", "core_value_id", "description", "quarter", "sender", "receiver"}

type appreciationsStore struct {
	BaseRepository
	AppreciationsTable string
	RewardsTable       string
	UsersTable         string
	CoreValuesTable    string
}

func NewAppreciationRepo(db *sqlx.DB) repository.AppreciationStorer {
	return &appreciationsStore{
		BaseRepository:     BaseRepository{db},
		AppreciationsTable: constants.AppreciationsTable,
		RewardsTable:       constants.RewardsTable,
		UsersTable:         constants.UsersTable,
		CoreValuesTable:    constants.CoreValuesTable,
	}
}

func (appr *appreciationsStore) CreateAppreciation(ctx context.Context, tx repository.Transaction, appreciation dto.Appreciation) (repository.Appreciation, error) {

	logger.Debug(ctx,"repoAppr: CreateAppreciation: appreciation: ",appreciation)
	queryExecutor := appr.InitiateQueryExecutor(tx)

	insertQuery, args, err := repository.Sq.
		Insert(appr.AppreciationsTable).Columns(AppreciationColumns[1:]...).
		Values(appreciation.CoreValueID, appreciation.Description, appreciation.Quarter, appreciation.Sender, appreciation.Receiver).
		Suffix("RETURNING id,core_value_id, description,total_reward_points,quarter,sender,receiver,created_at,updated_at").
		ToSql()
	if err != nil {
		logger.Errorf(ctx,"repoAppr: error in generating squirrel query, err: %v", err)
		return repository.Appreciation{}, apperrors.InternalServerError
	}

	logger.Infof(ctx,"repoAppr: insertQuery: %s , args: %v",insertQuery,args)

	var resAppr repository.Appreciation
	err = queryExecutor.QueryRowx(insertQuery, args...).StructScan(&resAppr)
	if err != nil {
		logger.Errorf(ctx,"repoAppr: Error executing create appreciation insert query: %v", err)
		return repository.Appreciation{}, apperrors.InternalServer
	}

	logger.Debug(ctx,"repoAppr: createappreciation response: ",resAppr)
	return resAppr, nil
}

func (appr *appreciationsStore) GetAppreciationById(ctx context.Context, tx repository.Transaction, apprId int32) (repository.AppreciationResponse, error) {

	logger.Debug(ctx,"repoAppr: GetAppreciationById: apprId: ",apprId)
	queryExecutor := appr.InitiateQueryExecutor(tx)

	// Get logged-in user ID
	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		logger.Error(ctx,"err in parsing userID from token")
		return repository.AppreciationResponse{}, apperrors.InternalServer
	}

	// Build the SQL query
	query, args, err := repository.Sq.Select(
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
				SELECT SUM(r2.point) 
				FROM rewards r2 
				WHERE r2.appreciation_id = a.id AND r2.sender = %d
			), 0) AS given_reward_point`, userID),
	).From(appr.AppreciationsTable+" a").
		LeftJoin(appr.UsersTable+" u_sender ON a.sender = u_sender.id").
		LeftJoin(appr.UsersTable+" u_receiver ON a.receiver = u_receiver.id").
		LeftJoin(appr.CoreValuesTable+" cv ON a.core_value_id = cv.id").
		LeftJoin(appr.RewardsTable+" r ON a.id = r.appreciation_id").
		Where(squirrel.And{
			squirrel.Eq{"a.id": apprId},
			squirrel.Eq{"a.is_valid": true},
		}).
		GroupBy("a.id", "cv.name", "u_sender.id", "u_receiver.id").
		ToSql()

	if err != nil {
		logger.Errorf(ctx,"error in generating squirrel query, err: %v", err)
		return repository.AppreciationResponse{}, apperrors.InternalServer
	}

	logger.Infof(ctx,"repoAppr: insertQuery: %s , args: %v",query,args)

	var resAppr repository.AppreciationResponse

	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).StructScan(&resAppr)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf(ctx,"repoAppr: no appreciation found with id: %d", apprId)
			return repository.AppreciationResponse{}, apperrors.AppreciationNotFound
		}
		logger.Errorf(ctx,"repoAppr: failed to execute query: %v", err)
		return repository.AppreciationResponse{}, apperrors.InternalServer
	}
	logger.Debug(ctx,"repoAppr:  appreciationById: ",resAppr)
	return resAppr, nil
}
func (appr *appreciationsStore) ListAppreciations(ctx context.Context, tx repository.Transaction, filter dto.AppreciationFilter) ([]repository.AppreciationResponse, repository.Pagination, error) {

	logger.Debug(ctx,"repoAppr:  ListAppreciations: filter: ",filter)
	queryExecutor := appr.InitiateQueryExecutor(tx)

	// Get logged-in user ID
	data := ctx.Value(constants.UserId)
	userID, ok := data.(int64)
	if !ok {
		logger.Error(ctx,"err in parsing userID from token")
		return []repository.AppreciationResponse{}, repository.Pagination{}, apperrors.InternalServerError
	}

	// query builder for counting total records
	queryBuilder := repository.Sq.Select("COUNT(*)").
		From("appreciations a").
		LeftJoin("users u_sender ON a.sender = u_sender.id").
		LeftJoin("users u_receiver ON a.receiver = u_receiver.id").
		LeftJoin("core_values cv ON a.core_value_id = cv.id").
		Where(squirrel.Eq{"a.is_valid": true})

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

	countSql, countArgs, err := queryBuilder.ToSql()
	if err != nil {
		logger.Error(ctx,"repoAppr: failed to build count query: ", err.Error())
		return []repository.AppreciationResponse{}, repository.Pagination{}, apperrors.InternalServerError
	}

	logger.Debug(ctx,fmt.Sprintf("repoAppr: listAppreciation: countSql: %s, countArgs: %v",countSql,countArgs))

	var totalRecords int32
	err = queryExecutor.QueryRowx(countSql, countArgs...).Scan(&totalRecords)
	if err != nil {
		logger.Error(ctx,"failed to execute count query: ", err.Error())
		return []repository.AppreciationResponse{}, repository.Pagination{}, apperrors.InternalServerError
	}

	pagination := getPaginationMetaData(filter.Page, filter.Limit, totalRecords)
	logger.Debug(ctx," pagination: ",pagination)
	queryBuilder = queryBuilder.RemoveColumns()
	queryBuilder = queryBuilder.Columns(
		"a.id",
		"cv.name AS core_value_name", fmt.Sprintf(
			`COALESCE((
				SELECT SUM(r2.point) 
				FROM rewards r2 
				WHERE r2.appreciation_id = a.id AND r2.sender = %d
			), 0) AS given_reward_point`, userID),
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
				SELECT SUM(r2.point) 
				FROM rewards r2 
				WHERE r2.appreciation_id = a.id AND r2.sender = %d
			), 0) AS given_reward_point`, userID),
	).
		LeftJoin("rewards r ON a.id = r.appreciation_id").
		GroupBy("a.id", "cv.name", "cv.description", "u_sender.id", "u_receiver.id")

	if filter.SortOrder != "" {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("a.created_at %s", filter.SortOrder))
	}

	offset := (filter.Page - 1) * filter.Limit
	// Add pagination
	queryBuilder = queryBuilder.Limit(uint64(filter.Limit)).Offset(uint64(offset))
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Error(ctx,"failed to build query: ", err.Error())
		return nil, repository.Pagination{}, apperrors.InternalServerError
	}

	logger.Debug(ctx,fmt.Sprintf("repoAppr: listAprreciation: sqlQuery: %s,args: %v",sql,args))
	queryExecutor = appr.InitiateQueryExecutor(tx)
	res := make([]repository.AppreciationResponse, 0)
	err = sqlx.Select(queryExecutor, &res, sql, args...)
	if err != nil {
		logger.Error(ctx,"repoAppr:failed to execute query appreciation: ", err.Error())
		return nil, repository.Pagination{}, apperrors.InternalServerError
	}
	logger.Error(ctx,"repoAppr: res data: ", res)
	id := ctx.Value(constants.UserId)
	logger.Debug(ctx," id -> ", id)
	userId, ok := ctx.Value(constants.UserId).(int64)
	if !ok {
		logger.Error(ctx,"unable to convert context user id to int64")
		return nil, repository.Pagination{}, apperrors.InternalServerError
	}

	logger.Debug(ctx," userId: ",userId)
	for idx, appreciation := range res {
		var userIds []int64
		queryBuilder = repository.Sq.Select("reported_by").From("resolutions").Where(squirrel.Eq{"appreciation_id": appreciation.ID})
		query, args, err := queryBuilder.ToSql()
		if err != nil {
			logger.Errorf(ctx,"error in generating squirrel query, err: %s", err.Error())
			return nil, repository.Pagination{}, apperrors.InternalServerError
		}

		logger.Debug(ctx,fmt.Sprintf("repoAppr: resolutions query: %s,args: %v",query,args))
		err = appr.DB.SelectContext(ctx, &userIds, query, args...)
		if err != nil {
			logger.Errorf(ctx,"error in reported flag query, err: %s", err.Error())
			return nil, repository.Pagination{}, apperrors.InternalServerError
		}
		res[idx].ReportedFlag = slices.Contains(userIds, userId)
	}

	logger.Debug(ctx,fmt.Sprintf("repoAppr: res: %v, pagination : %v",res,pagination))
	return res, pagination, nil
}

func (appr *appreciationsStore) DeleteAppreciation(ctx context.Context, tx repository.Transaction, apprId int32) error {
	logger.Debug(ctx,"repoAppr: apprId: ",apprId)
	query, args, err := repository.Sq.Update(appr.AppreciationsTable).
		Set("is_valid", false).
		Where(squirrel.And{
			squirrel.Eq{"id": apprId},
			squirrel.Eq{"is_valid": true},
		}).
		ToSql()

	logger.Debug(ctx,fmt.Sprintf("repoAppr: query: %s,args: %v",query,args))
	if err != nil {
		logger.Error(ctx,"Error building SQL: ", err.Error())
		return apperrors.InternalServer
	}

	queryExecutor := appr.InitiateQueryExecutor(tx)

	result, err := queryExecutor.Exec(query, args...)
	if err != nil {
		logger.Error(ctx,"Error executing SQL: ", err.Error())
		return apperrors.InternalServer
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(ctx,"Error getting rows affected: ", err.Error())
		return apperrors.InternalServer
	}

	if rowsAffected == 0 {
		logger.Warn(ctx,"No rows affected")
		return apperrors.AppreciationNotFound
	}

	return nil
}

func (appr *appreciationsStore) IsUserPresent(ctx context.Context, tx repository.Transaction, userID int64) (bool, error) {

	logger.Debug(ctx,"repoAppr: IsUserPresent: userID: ",userID)
	// Build the SQL query
	query, args, err := repository.Sq.Select("COUNT(*)").
		From(appr.UsersTable).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	logger.Debug(ctx," query: ",query)
	if err != nil {
		logger.Error(ctx,"err ", err.Error())
		return false, apperrors.InternalServer
	}

	logger.Debug(ctx,fmt.Sprintf("repoAppr: query: %s,args: %v",query,args))
	queryExecutor := appr.InitiateQueryExecutor(tx)

	var count int
	// Execute the query
	err = queryExecutor.QueryRowx(query, args...).Scan(&count)
	if err != nil {
		logger.Error(ctx,"failed to execute query: ", err.Error())
		return false, apperrors.InternalServer
	}

	logger.Debug(ctx,"repoAppr: count: ",count)
	// Check if user is present
	return count > 0, nil
}

func (appr *appreciationsStore) UpdateAppreciationTotalRewardsOfYesterday(ctx context.Context, tx repository.Transaction) (bool, error) {
	logger.Info(ctx,"appr: UpdateAppreciationTotalRewardsOfYesterday")
	// Initialize query executor
	queryExecutor := appr.InitiateQueryExecutor(tx)

	// Build the SQL update query with subquery
	query := `
UPDATE appreciations AS app
SET total_reward_points = total_reward_points + agg.total_points
FROM (
    SELECT appreciation_id, SUM(r.point * g.points) AS total_points
    FROM rewards r
    JOIN appreciations a ON r.appreciation_id = a.id
    JOIN users u ON r.sender = u.id
    JOIN grades g ON u.grade_id = g.id
    WHERE a.is_valid = true
      AND r.created_at >= EXTRACT(EPOCH FROM TIMESTAMP 'yesterday'::TIMESTAMP) * 1000
     AND r.created_at < EXTRACT(EPOCH FROM TIMESTAMP 'today'::TIMESTAMP) * 1000
    GROUP BY appreciation_id
) AS agg
WHERE app.id = agg.appreciation_id;
    `

	logger.Debug(ctx," query: ",query)
	// Execute the query using the query executor
	res, err := queryExecutor.Exec(query)
	if err != nil {
		logger.Error(ctx,"Error executing SQL query:", err.Error())
		return false, apperrors.InternalServer
	}

	rowsAffected,err := res.RowsAffected()
	if err != nil{
		logger.Error(ctx," err: ",err)
		return false,nil
	}
	logger.Info(ctx,"repoAppr: rowsAffected: ",rowsAffected)
	return true, nil
}

func (appr *appreciationsStore) UpdateUserBadgesBasedOnTotalRewards(ctx context.Context, tx repository.Transaction) ([]repository.UserBadgeDetails, error) {
	logger.Info(ctx," appr: UpdateUserBadgesBasedOnTotalRewards")
	queryExecutor := appr.InitiateQueryExecutor(tx)
	afterTime := GetQuarterStartUnixTime()

	logger.Info(ctx," aftertime: ",afterTime)
	query := `
		-- Calculate total reward points for each receiver
WITH receiver_points AS (
    SELECT
        receiver,
        SUM(total_reward_points) AS total_points
    FROM
        appreciations
    WHERE
        Appreciations.is_valid = true AND appreciations.created_at >= $1
    GROUP BY
        receiver
),

-- Determine eligible badges for each receiver
eligible_badges AS (
    SELECT
        rp.receiver AS user_id,
        b.id AS badge_id,
        ROW_NUMBER() OVER (PARTITION BY rp.receiver ORDER BY b.reward_points DESC) AS rn
    FROM
        receiver_points rp
    JOIN
        badges b ON rp.total_points >= b.reward_points
),

-- Check for existing badges created within the same period
existing_recent_badges AS (
    SELECT DISTINCT ON (ub.user_id, ub.badge_id)
        ub.user_id,
        ub.badge_id,
        ub.created_at
    FROM
        user_badges ub
    WHERE
        ub.created_at >= $2
),

-- Filter eligible badges that are not conflicted 
eligible_non_conflicted_badges AS (
    SELECT
        eb.user_id,
        eb.badge_id,
        (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT AS created_at
    FROM
        eligible_badges eb
    LEFT JOIN
        existing_recent_badges erb ON eb.user_id = erb.user_id AND eb.badge_id = erb.badge_id
    WHERE
        erb.user_id IS NULL
),

-- Insert eligible non-conflicting badges into user_badges
inserted_badges AS (
    INSERT INTO user_badges (badge_id, user_id, created_at)
    SELECT
        badge_id,
        user_id,
        created_at
    FROM
        eligible_non_conflicted_badges
    RETURNING user_id, badge_id
)

-- Retrieve user details and badge points for new badges
SELECT
	u.id,
    u.email,
    u.first_name,
    u.last_name,
    b.id AS badge_id,
	b.name AS badge_name,
    b.reward_points AS badge_points
FROM
    inserted_badges ib
JOIN
    users u ON ib.user_id = u.id
JOIN
    badges b ON ib.badge_id = b.id;
	`

	logger.Debug(ctx,fmt.Sprintf("repoAppr: query: %s, args: %v %v",query,afterTime,afterTime))
	rows, err := queryExecutor.Query(query, afterTime, afterTime)
	if err != nil {
		logger.Error(ctx,"repoAppr: error in extecution query")
		return []repository.UserBadgeDetails{}, err
	}
	defer rows.Close()

	var userBadgeDetails []repository.UserBadgeDetails
	for rows.Next() {
		var detail repository.UserBadgeDetails
		if err := rows.Scan(&detail.ID, &detail.Email, &detail.FirstName, &detail.LastName, &detail.BadgeID, &detail.BadgeName, &detail.BadgePoints); err != nil {
			return []repository.UserBadgeDetails{}, err
		}
		userBadgeDetails = append(userBadgeDetails, detail)
	}

	logger.Debug(ctx," userBadgeDetails: ",userBadgeDetails)
	if err = rows.Err(); err != nil {
		return []repository.UserBadgeDetails{}, err
	}

	return userBadgeDetails, nil
}
