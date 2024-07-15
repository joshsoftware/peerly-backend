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
	TableUsers     string
	TableGrades    string
	TableRoles     string
	TableOrgConfig string
}

func NewUserRepo(db *sqlx.DB) repository.UserStorer {
	return &userStore{
		DB:             db,
		TableUsers:     "users",
		TableGrades:    "grades",
		TableRoles:     "roles",
		TableOrgConfig: "organization_config",
	}
}

var (
	userColumns      = []string{"id", "employee_id", "first_name", "last_name", "email", "profile_image_url", "role_id", "reward_quota_balance", "designation", "grade_id"}
	rolesColumns     = []string{"id"}
	gradeColumns     = []string{"id", "name", "points"}
	orgConfigColumns = []string{"reward_multiplier"}
)

const (
// getUserByEmailQuery = `SELECT users.id, users.employee_id, users.first_name, users.last_name, users.email, users.profile_image_url, users.role_id, users.reward_quota_balance, users.designation, users.grade_id, grades.name FROM users JOIN grades ON grades.id = users.grade_id WHERE users.email = $1;
// `

// createUser = `INSERT INTO users ( email, employee_id, first_name, last_name, profile_image_url, role_id, reward_quota_balance, grade_id, designation
// ) VALUES (
// 	$1, $2, $3, $4, $5, $6, $7, $8, $9
// ) RETURNING id, employee_id, email, first_name, last_name, profile_image_url, role_id, reward_quota_balance, created_at, grade_id, designation`

// getRoleByNameQuery = `SELECT id FROM roles WHERE name=$1 LIMIT 1`

// getGradeId = `SELECT id, name, points FROM grades WHERE name = $1`

// getRewardMultiplier = "select reward_multiplier from organization_config where id = 1"

// updateUserQuery = `UPDATE users SET (first_name, last_name, profile_image_url, designation, grade_id) =
//
//	($1, $2, $3, $4, $5) where email = $6`
)

// GetUserByEmail - Given an email address, return that user.
func (us *userStore) GetUserByEmail(ctx context.Context, email string) (user repository.User, err error) {

	queryBuilder := repository.Sq.Select(userColumns...).From(us.TableUsers).Where(squirrel.Like{"email": email})
	getUserByEmailQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Error(fmt.Sprintf("error in generating squirrel query, err: %s", err.Error()))
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
			logger.Error(fmt.Sprintf("error selecting user from database by email, err: %s", err.Error()))
			err = apperrors.InternalServerError
			return
		}
	}

	return
}

func (us *userStore) GetRoleByName(ctx context.Context, name string) (roleId int, err error) {

	queryBuilder := repository.Sq.Select(rolesColumns...).From(us.TableRoles).Where(squirrel.Like{"name": name}).Limit(1)

	getRoleByNameQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %s", err.Error())
		return
	}

	err = us.DB.GetContext(ctx, &roleId, getRoleByNameQuery, args...)
	if err != nil {
		err = fmt.Errorf("error selecting role from database in GetRoleByName, grade: %s, err: %s", name, err.Error())
		return
	}
	return
}

// CreateNewUser - creates a new user in the database
func (us *userStore) CreateNewUser(ctx context.Context, user dto.User) (resp repository.User, err error) {

	queryBuilder := repository.Sq.Insert(us.TableUsers).Columns(userColumns[1:]...).Values(user.EmployeeId, user.FirstName, user.LastName, user.Email, user.ProfileImgUrl, user.RoleId, user.RewardQuotaBalance, user.Designation, user.GradeId).Suffix("RETURNING id, employee_id, email, first_name, last_name, profile_image_url, role_id, reward_quota_balance, created_at, grade_id, designation")

	createUser, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %s", err.Error())
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
		err = fmt.Errorf("error inserting user into database, email:%s, err: %s", user.Email, err.Error())
		return
	}

	return
}

func (us *userStore) GetGradeByName(ctx context.Context, name string) (grade repository.Grade, err error) {

	queryBuilder := repository.Sq.Select(gradeColumns...).From(us.TableGrades).Where(squirrel.Like{"name": name})
	getGradeId, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.Error(fmt.Sprintf("error in generating squirrel query, err: %s", err.Error()))
		err = apperrors.InternalServerError
		return
	}

	err = us.DB.GetContext(ctx, &grade, getGradeId, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = apperrors.GradeNotFound
			return
		}
		logger.Error(fmt.Sprintf("error in retriving grade id, grade: %s, err: %s", name, err.Error()))
		err = apperrors.InternalServerError
		return
	}
	return
}

func (us *userStore) GetRewardMultiplier(ctx context.Context) (value int, err error) {

	queryBuilder := repository.Sq.Select(orgConfigColumns...).From(us.TableOrgConfig).Where(squirrel.Eq{"id": 1})
	getRewardMultiplier, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %s", err.Error())
		return
	}

	err = us.DB.GetContext(ctx, &value, getRewardMultiplier, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("no fields in organization config, err: %s", err.Error())
			return
		}
		err = fmt.Errorf("error in retriving reward_multiplier from organization config, err: %s", err.Error())
		return
	}
	return
}

func (us *userStore) SyncData(ctx context.Context, updateData dto.User) (err error) {

	queryBuilder := repository.Sq.Update(us.TableUsers).
		Set("first_name", updateData.FirstName).
		Set("last_name", updateData.LastName).
		Set("profile_image_url", updateData.ProfileImgUrl).
		Set("designation", updateData.Designation).
		Set("grade_id", updateData.GradeId).
		Where(squirrel.Eq{"email": updateData.Email})

	updateUserQuery, args, err := queryBuilder.ToSql()
	if err != nil {
		err = fmt.Errorf("error in generating squirrel query, err: %s", err.Error())
		return
	}

	_, err = us.DB.ExecContext(
		ctx,
		updateUserQuery,
		args...,
	)
	if err != nil {
		err = fmt.Errorf("rrror in data update query, err: %s", err.Error())
		return
	}

	fmt.Println("Data update successful")

	return

}
