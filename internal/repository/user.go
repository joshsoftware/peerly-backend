package repository

import (
	"context"
	"database/sql"

	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type UserStorer interface {
	GetUserByEmail(ctx context.Context, email string) (user dto.GetUserResp, err error)
	GetRoleByName(ctx context.Context, name string) (roleId int, err error)
	CreateNewUser(ctx context.Context, u dto.RegisterUser) (resp dto.GetUserResp, err error)
	GetGradeByName(ctx context.Context, name string) (grade Grade, err error)
	GetRewardMultiplier(ctx context.Context) (value int, err error)
	SyncData(ctx context.Context, updateData dto.UpdateUserData) (err error)
	GetUserList(ctx context.Context, reqData dto.UserListReq) (resp []dto.GetUserListResp, err error)
}

// User - basic struct representing a User
type User struct {
	ID                  int           `db:"id" json:"id"`
	FirstName           string        `db:"first_name" json:"first_name"`
	LastName            string        `db:"last_name" json:"last_name"`
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

type Grade struct {
	Id     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Points int    `db:"points" json:"points"`
}
