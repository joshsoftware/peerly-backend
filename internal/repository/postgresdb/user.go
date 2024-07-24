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

type userStore struct {
	DB             *sqlx.DB
	UsersTable     string
	GradesTable    string
	RolesTable     string
	OrgConfigTable string
}

func NewUserRepo(db *sqlx.DB) repository.UserStorer {
	return &userStore{
		DB:             db,
		UsersTable:     "users",
		GradesTable:    "grades",
		RolesTable:     "roles",
		OrgConfigTable: "organization_config",
	}
}

var (
	userColumns      = []string{"id", "employee_id", "first_name", "last_name", "email", "profile_image_url", "role_id", "reward_quota_balance", "designation", "grade_id"}
	rolesColumns     = []string{"id"}
	gradeColumns     = []string{"id", "name", "points"}
	orgConfigColumns = []string{"reward_multiplier"}
)

// GetUserByEmail - Given an email address, return that user.
func (us *userStore) GetUserByEmail(ctx context.Context, email string) (user repository.User, err error) {

	queryBuilder := repository.Sq.Select(userColumns...).From(us.UsersTable).Where(squirrel.Like{"email": email})
	getUserByEmailQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf("error in generating squirrel query, err: %s", err.Error())
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
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
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
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
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
		logger.Errorf("error in generating squirrel query, err: %s", err)
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
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
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
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
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

func (us *userStore) GetTotalUserCount(ctx context.Context, reqData dto.UserListReq) (totalCount int64, err error) {

	queryBuilder := repository.Sq.Select("count(*)").From("users")
	conditions := []squirrel.Sqlizer{}
	for _, name := range reqData.Name {
		if name != "" {
			conditions = append(conditions, squirrel.Like{"lower(first_name)": "%" + name + "%"})
			conditions = append(conditions, squirrel.Like{"lower(last_name)": "%" + name + "%"})
		}
	}
	if len(conditions) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Or(conditions))
	}

	getUserCountQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %w", err)
		return
	}

	var resp []int64

	err = us.DB.Select(&resp, getUserCountQuery, args...)
	if err != nil {
		err = fmt.Errorf("error in getUserCountQuery, err:%w", err)
		return
	}
	return
}

func (us *userStore) ListUsers(ctx context.Context, reqData dto.UserListReq) (resp []repository.User, err error) {

	queryBuilder := repository.Sq.Select(userColumns...).From(us.UsersTable)
	conditions := []squirrel.Sqlizer{}
	for _, name := range reqData.Name {
		if name != "" {
			conditions = append(conditions, squirrel.Like{"lower(first_name)": "%" + name + "%"})
			conditions = append(conditions, squirrel.Like{"lower(last_name)": "%" + name + "%"})
		}
	}
	if len(conditions) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Or(conditions))
	}
	offset := reqData.PerPage * (reqData.Page - 1)
	queryBuilder = queryBuilder.Limit(uint64(reqData.PerPage)).Offset(uint64(offset))

	listUsersQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Errorf("error in generating squirrel query, err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}

	err = us.DB.Select(&resp, listUsersQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("no fields returned, err:%s", err.Error())
			err = nil
			return
		}
		logger.Errorf("error in fetching users from database, err: %s", err.Error())
		err = apperrors.InternalServerError
		return
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

func (us *userStore) GetGradeById(ctx context.Context, id int64) (grade repository.Grade, err error) {
	getGradeById := `SELECT id, name, points FROM grades WHERE id = $1`
	err = us.DB.Get(&grade, getGradeById, id)
	if err != nil {
		err = fmt.Errorf("error in getGradeById. id:%d. err:%w", id, err)
		return
	}
	return

}
