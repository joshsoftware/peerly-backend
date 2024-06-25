package repository

import (
	"context"
	"database/sql"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

type userStore struct {
	DB *sqlx.DB
}

type UserStorer interface {
	GetUserByEmail(ctx context.Context, email string) (user dto.GetUserResp, err error)
	GetRoleByName(ctx context.Context, name string) (roleId int, err error)
	CreateNewUser(ctx context.Context, u dto.RegisterUser) (resp dto.GetUserResp, err error)
	GetGradeByName(ctx context.Context, name string) (id int, err error)
}

func NewUserRepo(db *sqlx.DB) UserStorer {
	return &userStore{
		DB: db,
	}
}

const (
	getUserByEmailQuery = `SELECT id, first_name, org_id, last_name, email, profile_image_url, role_id, reward_quota_balance, designation, grade_id FROM users WHERE email=$1 LIMIT 1`

	createUser = `INSERT INTO users (
		org_id, email, first_name, last_name, profile_image_url, role_id, reward_quota_balance, created_at, grade_id, designation
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	) RETURNING id, org_id, email, first_name, last_name, profile_image_url, role_id, reward_quota_balance, created_at, grade_id, designation`

	getRoleByNameQuery = `SELECT id FROM roles WHERE name=$1 LIMIT 1`

	getGradeId = `SELECT id FROM grade WHERE name = $1`
)

// User - basic struct representing a User
type User struct {
	ID                  int           `db:"id" json:"id"`
	FirstName           string        `db:"first_name" json:"first_name"`
	LastName            string        `db:"last_name" json:"last_name"`
	OrgID               int           `db:"org_id" json:"org_id"`
	Email               string        `db:"email" json:"email"`
	ProfileImageURL     string        `db:"profile_image_url" json:"profile_image_url"`
	Grade               int           `db:"grade" json:"grade"`
	Designation         string        `db:"designation" json:"designation"`
	RoleID              int           `db:"role_id" json:"role_id"`
	RewardsQuotaBalance int           `db:"rewards_quota_balance" json:"rewards_quota_balance"`
	Status              int           `db:"status" json:"status"`
	SoftDelete          bool          `db:"soft_delete" json:"soft_delete,omitempty"`
	SoftDeleteBy        sql.NullInt64 `db:"soft_delete_by" json:"soft_delete_by,omitempty"`
	SoftDeleteOn        sql.NullTime  `db:"soft_delete_on" json:"soft_delete_on,omitempty"`
	CreatedAt           time.Time     `db:"created_at" json:"created_at"`
}

type Role struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

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

	u.CreatedAt = time.Now().UTC()

	err = us.DB.GetContext(
		ctx,
		&resp,
		createUser,
		u.OrgId,
		u.User.Email,
		u.User.FirstName,
		u.User.LastName,
		u.User.ProfileImgUrl,
		u.RoleId,
		u.RewardQuotaBalance,
		u.CreatedAt,
		u.GradeId,
		u.User.Designation,
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
