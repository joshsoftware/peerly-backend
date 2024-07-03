package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type userStore struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) repository.UserStorer {
	return &userStore{
		DB: db,
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

	getGradeId = `SELECT id FROM grades WHERE name = $1`

	getRewardQuotaBalanceDefault = "select reward_multiplier from organization_config where id = 1"

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

func (us *userStore) GetGradeByName(ctx context.Context, name string) (id int, err error) {
	err = us.DB.Get(&id, getGradeId, name)
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

func (us *userStore) GetRewardOuotaDefault(ctx context.Context) (id int, err error) {
	err = us.DB.Get(&id, getRewardQuotaBalanceDefault)
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
