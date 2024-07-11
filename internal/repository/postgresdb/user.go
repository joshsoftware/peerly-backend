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

	for i, name := range reqData.Name {
		if name != "" {
			if i == 0 {
				getUserCountQuery += "where"
				str := fmt.Sprint(" lower(first_name) like '%" + name + "%' or lower(last_name) like '%" + name + "%'")
				getUserCountQuery += str
			} else {
				str := fmt.Sprint(" or lower(first_name) like '%" + name + "%' or lower(last_name) like '%" + name + "%'")
				getUserCountQuery += str
			}
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

	// getUserListQuery := "Select users.employee_id, users.email, users.first_name, users.last_name, grades.name, users.designation, users.profile_image_url from users join grades on grades.id = users.grade_id "

	getUserListQuery := "Select users.id, users.email, users.first_name, users.last_name from users "

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

	str := fmt.Sprint(" limit " + strconv.Itoa(int(reqData.PerPage)) + " offset " + strconv.Itoa(int(reqData.PerPage*(reqData.Page-1))))
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

	// for _, user := range dbResp {
	// 	var respUser dto.GetUserListResp
	// 	respUser.EmployeeId = user.EmployeeId
	// 	respUser.Email = user.Email
	// 	respUser.FirstName = user.FirstName
	// 	respUser.LastName = user.LastName
	// 	respUser.Grade = user.Grade
	// 	respUser.Designation = user.Designation
	// 	respUser.ProfileImg = user.ProfileImg.String
	// 	resp = append(resp, respUser)
	// }

	return
}
