package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type userStore struct {
	BaseRepository
}

func NewUserRepo(db *sqlx.DB) repository.UserStorer {
	return &userStore{
		BaseRepository: BaseRepository{db},
	}
}

const (
	getUserByEmailQuery = `SELECT users.id, users.employee_id, users.first_name, users.last_name, users.email, users.profile_image_url, users.role_id, users.reward_quota_balance, users.designation, users.grade_id, grades.name FROM users JOIN grades ON grades.id = users.grade_id WHERE users.email = $1;
`

	createUser = `INSERT INTO users ( email, employee_id, first_name, last_name, profile_image_url, role_id, reward_quota_balance, grade_id, designation
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	) RETURNING id, employee_id, email, first_name, last_name, profile_image_url, role_id, reward_quota_balance, created_at, grade_id, designation`

	getRoleByNameQuery = `SELECT id FROM roles WHERE name=$1 LIMIT 1`

	getGradeId = `SELECT id, name, points FROM grades WHERE name = $1`

	getRewardMultiplier = "select reward_multiplier from organization_config where id = 1"

	updateUserQuery = `UPDATE users SET (first_name, last_name, profile_image_url, designation, grade_id) =
		($1, $2, $3, $4, $5) where email = $6`
)

// GetUserByEmail - Given an email address, return that user.
func (us *userStore) GetUserByEmail(ctx context.Context, email string) (user dto.GetUserResp, err error) {
	err = us.DB.Get(&user, getUserByEmailQuery, email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = apperrors.UserNotFound
			return
		}
		// Possible that there's no rows in the result set
		logger.WithField("err", err.Error()).Error("Error selecting user from database by email " + email)
		err = apperrors.InternalServerError
		return
	}

	return
}

func (us *userStore) GetRoleByName(ctx context.Context, name string) (roleId int, err error) {
	err = us.DB.Get(&roleId, getRoleByNameQuery, name)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting role from database in GetRoleByName for name ", name)
		return
	}
	return
}

// CreateNewUser - creates a new user in the database
func (us *userStore) CreateNewUser(ctx context.Context, u dto.RegisterUser) (resp dto.GetUserResp, err error) {

	err = us.DB.GetContext(
		ctx,
		&resp,
		createUser,
		u.User.Email,
		u.User.EmpolyeeDetail.EmployeeId,
		u.User.PublicProfile.FirstName,
		u.User.PublicProfile.LastName,
		u.User.PublicProfile.ProfileImgUrl,
		u.RoleId,
		u.RewardQuotaBalance,
		u.GradeId,
		u.User.EmpolyeeDetail.Designation.Name,
	)

	if err != nil {
		// FAIL: Could not run insert query
		logger.WithField("err", err.Error()).Error("Error inserting user into database: " + u.User.Email)

		return
	}

	return
}

func (us *userStore) GetGradeByName(ctx context.Context, name string) (grade repository.Grade, err error) {
	err = us.DB.Get(&grade, getGradeId, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = apperrors.GradeNotFound
			return
		}
		logger.WithField("err", err.Error()).Error("Error in retriving grade id of the grade ", name)
		err = apperrors.InternalServerError
		return
	}
	return
}

func (us *userStore) GetRewardMultiplier(ctx context.Context) (value int, err error) {
	err = us.DB.Get(&value, getRewardMultiplier)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("No fields in organization config")
			return
		}
		logger.WithField("err", err.Error()).Error("Error in retriving reward_multiplier from organization config")
		return
	}
	return
}

func (us *userStore) SyncData(ctx context.Context, updateData dto.UpdateUserData) (err error) {
	_, err = us.DB.ExecContext(
		ctx,
		updateUserQuery,
		updateData.FirstName,
		updateData.LastName,
		updateData.ProfileImgUrl,
		updateData.Designation,
		updateData.GradeId,
		updateData.Email,
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in data update query")
		return
	}

	fmt.Println("Data update successful")

	return

}

func (us *userStore) GetTotalUserCount(ctx context.Context, reqData dto.UserListReq) (totalCount int64, err error) {

	getUserCountQuery := "Select count(*) from users "

	if len(reqData.Name) >= 0 {
		getUserCountQuery += "where"
	}
	for i, name := range reqData.Name {
		if i == 0 {
			str := fmt.Sprint(" lower(first_name) like '%" + name + "%' or lower(last_name) like '%" + name + "%'")
			getUserCountQuery += str
		} else {
			str := fmt.Sprint(" or lower(first_name) like '%" + name + "%' or lower(last_name) like '%" + name + "%'")
			getUserCountQuery += str
		}
	}

	var resp []int64

	err = us.DB.Select(&resp, getUserCountQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in getUserCountQuery")
		err = apperrors.InternalServerError
		return
	}

	totalCount = resp[0]

	return
}

func (us *userStore) GetUserList(ctx context.Context, reqData dto.UserListReq) (resp []dto.GetUserListResp, err error) {
	getUserListQuery := "Select users.employee_id, users.email, users.first_name, users.last_name, grades.name, users.designation, users.profile_image_url from users join grades on grades.id = users.grade_id "

	if len(reqData.Name) >= 0 {
		getUserListQuery += "where"
	}
	for i, name := range reqData.Name {
		if i == 0 {
			str := fmt.Sprint(" lower(first_name) like '%" + name + "%' or lower(last_name) like '%" + name + "%'")
			getUserListQuery += str
		} else {
			str := fmt.Sprint(" or lower(first_name) like '%" + name + "%' or lower(last_name) like '%" + name + "%'")
			getUserListQuery += str
		}
	}

	str := fmt.Sprint(" limit " + strconv.Itoa(reqData.PerPage) + " offset " + strconv.Itoa(reqData.PerPage*(reqData.Page-1)))
	getUserListQuery += str

	err = us.DB.Select(&resp, getUserListQuery)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("err", err.Error()).Error("No fields returned")
			err = nil
			return
		}
		logger.WithField("err", err.Error()).Error("Error in fetching users from database")
		err = apperrors.InternalServerError
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
     WHERE ub.id = (SELECT MAX(id) FROM user_badges WHERE user_id = ub.user_id)) AS b ON u.id = b.user_id;
`

	rows, err := queryExecutor.Query(query,afterTime,afterTime,afterTime)
	if err != nil {
		logger.Error("err: userStore ",err.Error())
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
			logger.Error("err: userStore ",err.Error())
			return nil, err
		}
		activeUsers = append(activeUsers, user)
	}

	if err = rows.Err(); err != nil {
		logger.Error("err: userStore ",err.Error())
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
		select count(*) 
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
		user.TotalPoints = userList[0].TotalPoints
		user.Badge = userList[0].Badge.String
		user.BadgeCreatedAt = userList[0].BadgeCreatedAt.Int64
	}
	return
}
