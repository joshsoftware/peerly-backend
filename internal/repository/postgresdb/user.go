package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

var (
	leaderBoardColumns = []string{"id", "first_name", "last_name", "profile_image_url", "badge_name", "appreciation_points"}
)

type userStore struct {
	BaseRepository
	UserTable      string
	UsersTable     string
	GradesTable    string
	RolesTable     string
	OrgConfigTable string
}

func NewUserRepo(db *sqlx.DB) repository.UserStorer {
	return &userStore{
		BaseRepository: BaseRepository{db},
		UsersTable:     constants.UsersTable,
		GradesTable:    constants.GradesTable,
		RolesTable:     constants.RolesTable,
		OrgConfigTable: constants.OrganizationConfigTable,
	}
}

var (
	userColumns      = []string{"id", "employee_id", "first_name", "last_name", "email", "profile_image_url", "role_id", "reward_quota_balance", "designation", "grade_id"}
	adminColumns     = []string{"id", "employee_id", "first_name", "last_name", "email", "password", "profile_image_url", "role_id", "reward_quota_balance", "designation", "grade_id"}
	rolesColumns     = []string{"id"}
	orgConfigColumns = []string{"reward_multiplier"}
)

// GetUserByEmail - Given an email address, return that user.
func (us *userStore) GetUserByEmail(ctx context.Context, email string) (user repository.User, err error) {

	queryBuilder := repository.Sq.Select(userColumns...).From(us.UsersTable).Where(squirrel.Like{"email": email})
	getUserByEmailQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf("error in generating query, err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}

	err = us.DB.GetContext(
		ctx,
		&user,
		getUserByEmailQuery,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = apperrors.UserNotFound
			return
		} else {
			// Possible that there's no rows in the result set
			logger.Errorf("error selecting user from database by email, err: %s", err.Error())
			err = apperrors.InternalServerError
			return
		}
	}

	return
}

func (us *userStore) GetRoleByName(ctx context.Context, name string) (roleId int64, err error) {

	queryBuilder := repository.Sq.Select(rolesColumns...).From(us.RolesTable).Where(squirrel.Like{"name": name})

	getRoleByNameQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating query, err: %w", err)
		return
	}

	err = us.DB.GetContext(ctx, &roleId, getRoleByNameQuery, args...)
	if err != nil {
		err = fmt.Errorf("error selecting role from database in GetRoleByName, grade: %s, err: %w", name, err)
		return
	}
	return
}

// CreateNewUser - creates a new user in the database
func (us *userStore) CreateNewUser(ctx context.Context, user dto.User) (resp repository.User, err error) {

	queryBuilder := repository.Sq.Insert(us.UsersTable).Columns(userColumns[1:]...).Values(user.EmployeeId, user.FirstName, user.LastName, user.Email, user.ProfileImgUrl, user.RoleId, user.RewardQuotaBalance, user.Designation, user.GradeId).Suffix("RETURNING id, employee_id, email, first_name, last_name, profile_image_url, role_id, reward_quota_balance, created_at, grade_id, designation")

	createUser, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating query, err: %w", err)
		return
	}

	err = us.DB.GetContext(
		ctx,
		&resp,
		createUser,
		args...,
	)

	if err != nil {
		// FAIL: Could not run insert query
		err = fmt.Errorf("error inserting user into database, email:%s, err: %w", user.Email, err)
		return
	}

	return
}

func (us *userStore) GetGradeByName(ctx context.Context, name string) (grade repository.Grade, err error) {

	queryBuilder := repository.Sq.Select(gradeColumns...).From(us.GradesTable).Where(squirrel.Like{"name": name})
	getGradeId, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf("error in generating query, err: %s", err)
		err = apperrors.InternalServerError
		return
	}

	err = us.DB.GetContext(ctx, &grade, getGradeId, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = apperrors.GradeNotFound
			return
		}
		logger.Errorf("error in retriving grade id, grade: %s, err: %s", name, err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}

func (us *userStore) GetRewardMultiplier(ctx context.Context) (value int64, err error) {

	queryBuilder := repository.Sq.Select(orgConfigColumns...).From(us.OrgConfigTable).Where(squirrel.Eq{"id": 1})
	getRewardMultiplier, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating query, err: %w", err)
		return
	}

	err = us.DB.GetContext(ctx, &value, getRewardMultiplier, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("no fields in organization config, err: %w", err)
			return
		}
		err = fmt.Errorf("error in retriving reward_multiplier from organization config, err: %w", err)
		return
	}
	return
}

func (us *userStore) SyncData(ctx context.Context, updateData dto.User) (err error) {

	queryBuilder := repository.Sq.Update(us.UsersTable).
		Set("first_name", updateData.FirstName).
		Set("last_name", updateData.LastName).
		Set("profile_image_url", updateData.ProfileImgUrl).
		Set("designation", updateData.Designation).
		Set("grade_id", updateData.GradeId).
		Where(squirrel.Eq{"email": updateData.Email})

	updateUserQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating query, err: %w", err)
		return
	}

	_, err = us.DB.ExecContext(
		ctx,
		updateUserQuery,
		args...,
	)
	if err != nil {
		err = fmt.Errorf("error in data update query, err: %w", err)
		return
	}

	fmt.Println("Data update successful")

	return

}

func (us *userStore) GetTotalUserCount(ctx context.Context, reqData dto.ListUsersReq) (totalCount int64, err error) {

	queryBuilder := repository.Sq.Select("count(*)").From(us.UsersTable)
	conditions := []squirrel.Sqlizer{}
	for _, name := range reqData.Name {
		conditions = append(conditions, squirrel.Like{"lower(first_name)": "%" + name + "%"})
		conditions = append(conditions, squirrel.Like{"lower(last_name)": "%" + name + "%"})
	}
	if len(conditions) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Or(conditions))
	}

	getUserCountQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating query, err: %w", err)
		return
	}

	err = us.DB.GetContext(ctx, &totalCount, getUserCountQuery, args...)
	if err != nil {
		err = fmt.Errorf("error in getUserCountQuery, err:%w", err)
		return
	}
	return
}

func (us *userStore) ListUsers(ctx context.Context, reqData dto.ListUsersReq) (resp []repository.User, count int64, err error) {

	count, err = us.GetTotalUserCount(ctx, reqData)
	if err != nil {
		return
	}

	queryBuilder := repository.Sq.Select(userColumns...).From(us.UsersTable).OrderBy("first_name")
	conditions := []squirrel.Sqlizer{}
	for _, name := range reqData.Name {
		conditions = append(conditions, squirrel.Like{"lower(first_name)": "%" + name + "%"})
		conditions = append(conditions, squirrel.Like{"lower(last_name)": "%" + name + "%"})
	}
	if len(conditions) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Or(conditions))
	}
	offset := reqData.PageSize * (reqData.Page - 1)
	queryBuilder = queryBuilder.Limit(uint64(reqData.PageSize)).Offset(uint64(offset))

	listUsersQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating query, err: %w", err)
		return
	}

	err = us.DB.Select(&resp, listUsersQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("no fields returned, err:%s", err.Error())
			err = nil
			return
		}
		err = fmt.Errorf("error in fetching users from database, err: %w", err)
		return
	}

	return
}

func (us *userStore) GetActiveUserList(ctx context.Context, tx repository.Transaction) (activeUsers []repository.ActiveUser, err error) {
	queryExecutor := us.InitiateQueryExecutor(tx)
	afterTime := GetQuarterStartUnixTime()
	query := `WITH user_points AS (
    SELECT 
        u.id AS user_id,
        COALESCE(received.total_received_appreciations, 0) AS total_received_appreciations,
        COALESCE(sent.total_sent_appreciations, 0) AS total_sent_appreciations,
        COALESCE(given.total_given_rewards, 0) AS total_given_rewards,
        (3 * COALESCE(sent.total_sent_appreciations, 0) + 2 * COALESCE(received.total_received_appreciations, 0) + COALESCE(given.total_given_rewards, 0)) AS active_user_points
    FROM 
        users u
    LEFT JOIN 
        (SELECT receiver AS user_id, COUNT(*) AS total_received_appreciations 
         FROM appreciations 
		 WHERE
        Appreciations.is_valid = true AND appreciations.created_at >=$1
         GROUP BY receiver
		 ) AS received ON u.id = received.user_id
    LEFT JOIN 
        (SELECT sender AS user_id, COUNT(*) AS total_sent_appreciations 
         FROM appreciations
		 WHERE
        Appreciations.is_valid = true AND appreciations.created_at >=$2
         GROUP BY sender) AS sent ON u.id = sent.user_id
    LEFT JOIN 
        (SELECT sender AS user_id, COUNT(*) AS total_given_rewards 
         FROM rewards
		 WHERE
		 rewards.created_at >=$3
         GROUP BY sender) AS given ON u.id = given.user_id
    WHERE
        COALESCE(received.total_received_appreciations, 0) > 0 OR
        COALESCE(sent.total_sent_appreciations, 0) > 0 OR
        COALESCE(given.total_given_rewards, 0) > 0
    ORDER BY
        active_user_points DESC,
        total_sent_appreciations DESC,
        total_given_rewards DESC,
        total_received_appreciations DESC
)
SELECT 
    up.user_id,
    u.first_name ,
	u.last_name ,
    u.profile_image_url,
    b.name AS badge,
    COALESCE(ap.appreciation_points, 0) AS appreciation_points
FROM 
    user_points up
JOIN 
    users u ON up.user_id = u.id
LEFT JOIN 
    (SELECT receiver, SUM(total_reward_points) AS appreciation_points 
     FROM appreciations 
     GROUP BY receiver) AS ap ON u.id = ap.receiver
LEFT JOIN 
    (SELECT ub.user_id, b.name
     FROM user_badges ub 
     JOIN badges b ON ub.badge_id = b.id
     WHERE ub.id = (SELECT MAX(id) FROM user_badges WHERE user_id = ub.user_id)) AS b ON u.id = b.user_id
LIMIT 10;
`
	logger.Info("afterTime: ", afterTime)

	rows, err := queryExecutor.Query(query, afterTime, afterTime, afterTime)
	if err != nil {
		logger.Error("err: userStore ", err.Error())
		return []repository.ActiveUser{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user repository.ActiveUser
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.ProfileImageURL,
			&user.BadgeName,
			&user.AppreciationPoints,
		); err != nil {
			logger.Error("err: userStore ", err.Error())
			return nil, err
		}
		activeUsers = append(activeUsers, user)
	}

	if err = rows.Err(); err != nil {
		logger.Error("err: userStore ", err.Error())
		return []repository.ActiveUser{}, err
	}

	return activeUsers, nil
}

func (us *userStore) UpdateRewardQuota(ctx context.Context, tx repository.Transaction) (err error) {

	queryExecutor := us.InitiateQueryExecutor(tx)
	query := `UPDATE users
	SET reward_quota_balance = (
    SELECT oc.reward_multiplier * g.points
    FROM organization_config oc,grades g
    WHERE users.grade_id = g.id
	)`

	_, err = queryExecutor.Exec(query)
	if err != nil {
		logger.Error("err: userStore ", err.Error())
		return err
	}
	return
}

func (us *userStore) GetUserById(ctx context.Context, reqData dto.GetUserByIdReq) (user dto.GetUserByIdResp, err error) {

	getUserById := `select users.id, users.first_name, users.last_name, users.email, users.profile_image_url, users.designation, users.reward_quota_balance, users.grade_id, users.employee_id, 
		(
		select sum(total_reward_points) 
		from appreciations
		where
		receiver = users.id
		and
		appreciations.created_at >= $1
		) as total_points, 
	badges.name, user_badges.created_at as badge_created_at 
	from users
	left join user_badges 
	on user_badges.user_id = users.id
	left join badges
	on user_badges.badge_id = badges.id
	where users.id = $2
	group by users.id, badges.name, user_badges.id
	order by user_badges.created_at desc`

	var userList []dto.GetUserByIdDbResp

	err = us.DB.Select(&userList, getUserById, reqData.QuaterTimeStamp, reqData.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("No fields returned")
			err = apperrors.InvalidId
			return
		}
		logger.WithField("err", err.Error()).Error("Error in fetching users from database")
		err = apperrors.InternalServerError
		return
	}

	if (userList[0].BadgeCreatedAt.Valid && userList[0].BadgeCreatedAt.Int64 >= reqData.QuaterTimeStamp) || !userList[0].BadgeCreatedAt.Valid {
		user.UserId = userList[0].UserId
		user.FirstName = userList[0].FirstName
		user.LastName = userList[0].LastName
		user.Email = userList[0].Email
		user.ProfileImgUrl = userList[0].ProfileImgUrl
		user.Designation = userList[0].Designation
		user.RewardQuotaBalance = userList[0].RewardQuotaBalance
		user.GradeId = userList[0].GradeId
		user.EmployeeId = userList[0].EmployeeId
		user.TotalPoints = userList[0].TotalPoints.Int64
		user.Badge = userList[0].Badge.String
		user.BadgeCreatedAt = userList[0].BadgeCreatedAt.Int64
	}
	return
}

func (us *userStore) GetTop10Users(ctx context.Context, quarterTimestamp int64) (users []repository.Top10Users, err error) {

	afterTime := GetQuarterStartUnixTime()
	getTop10UserQuery := `select users.id, users.first_name, users.last_name, users.profile_image_url, sum(appreciations.total_reward_points) as AP from users join appreciations on users.id = appreciations.receiver where appreciations.created_at >= $1 group by users.id, appreciations.receiver order by AP desc limit 10`

	err = us.DB.Select(&users, getTop10UserQuery, afterTime)
	if err != nil {
		err = fmt.Errorf("err in getTop10UsersQuery err: %w", err)
		return
	}

	getUserBadge := `select badges.name from badges join user_badges on user_badges.badge_id = badges.id where user_badges.user_id = $1 and created_at >= $2 group by badges.id, user_badges.created_at, user_badges.badge_id order by user_badges.badge_id desc limit 1`

	for i, user := range users {
		var badge []sql.NullString
		err = us.DB.Select(&badge, getUserBadge, user.ID, quarterTimestamp)
		if err != nil {
			err = fmt.Errorf("err in getUserBadge query. userId:%d, err: %w", user.ID, err)
			return
		}

		if len(badge) > 0 {
			fmt.Println("badge: ", badge[0], " for id: ", user.ID)
			user.BadgeName = badge[0]
			users[i] = user
		}

	}

	return
}

func (us *userStore) GetGradeById(ctx context.Context, id int64) (grade repository.Grade, err error) {
	getGradeById := `SELECT id, name, points FROM grades WHERE id = $1`
	err = us.DB.Get(&grade, getGradeById, id)
	if err != nil {
		err = fmt.Errorf("error in getGradeById. id:%d. err:%w", id, err)
		return
	}
	return
}

func (us *userStore) GetAdmin(ctx context.Context, email string) (user repository.User, err error) {
	queryBuilder := repository.Sq.Select(adminColumns...).From(us.UsersTable).Where(squirrel.Like{"email": email})
	getAdminQuery, args, err := queryBuilder.ToSql()
	err = us.DB.GetContext(
		ctx,
		&user,
		getAdminQuery,
		args...,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("invalid user email, err:%s", err.Error())
			err = apperrors.InvalidEmail
			return
		}
		logger.Errorf("error in get admin query, err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}
	return
}

func (us *userStore) AddDeviceToken(ctx context.Context, userID int64, notificationToken string) (err error) {

	if notificationToken == "" {
		return nil
	}
	insertQuery, args, err := repository.Sq.
		Insert("notification_tokens").Columns("user_id", "notification_token").
		Values(userID, notificationToken).
		Suffix("RETURNING id,user_id,notification_token").
		ToSql()
	if err != nil {
		logger.Errorf("error in generating squirrel query, err: %v", err)
		return apperrors.InternalServerError
	}

	type Device struct {
		ID                int32  `db:"id"`
		UserID            int64  `db:"user_id"`
		NotificationToken string `db:"notification_token"`
	}

	var device Device
	// Execute the query
	err = us.DB.QueryRowx(insertQuery, args...).StructScan(&device)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("device not found")
			return apperrors.InternalServerError
		}
		logger.Errorf("failed to execute query: %v", err)
		return apperrors.InternalServerError
	}
	return nil
}

func (us *userStore) ListDeviceTokensByUserID(ctx context.Context, userID int64) (notificationTokens []string, err error) {
	notificationTokenQuery := "SELECT notification_token FROM notification_tokens WHERE user_id = $1"
	err = us.DB.Select(&notificationTokens, notificationTokenQuery, userID)
	if err != nil {
		err = fmt.Errorf("error in ListDeviceTokensByUserID: %w", err)
		return
	}
	if len(notificationTokens) <= 0 {
		fmt.Println("notification tokens: ", notificationTokens)
	}
	return
}
