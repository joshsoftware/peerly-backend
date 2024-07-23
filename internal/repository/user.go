package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type UserStorer interface {
	GetUserByEmail(ctx context.Context, email string) (user User, err error)
	GetRoleByName(ctx context.Context, name string) (roleId int64, err error)
	CreateNewUser(ctx context.Context, user dto.User) (resp User, err error)
	GetGradeByName(ctx context.Context, name string) (grade Grade, err error)
	GetRewardMultiplier(ctx context.Context) (value int64, err error)
	SyncData(ctx context.Context, updateData dto.User) (err error)
	GetUserList(ctx context.Context, reqData dto.UserListReq) (resp []dto.GetUserListResp, err error)
	GetTotalUserCount(ctx context.Context, reqData dto.UserListReq) (totalCount int64, err error)
}

// User - basic struct representing a User
type User struct {
	Id                  int64         `db:"id"`
	EmployeeId          string        `db:"employee_id"`
	FirstName           string        `db:"first_name"`
	LastName            string        `db:"last_name"`
	Email               string        `db:"email"`
	ProfileImageURL     string        `db:"profile_image_url"`
	GradeId             int64         `db:"grade_id"`
	Designation         string        `db:"designation"`
	RoleID              int64         `db:"role_id"`
	RewardsQuotaBalance int64         `db:"reward_quota_balance"`
	Status              int64         `db:"status"`
	SoftDelete          bool          `db:"soft_delete"`
	SoftDeleteBy        sql.NullInt64 `db:"soft_delete_by"`
	SoftDeleteOn        sql.NullTime  `db:"soft_delete_on"`
	CreatedAt           int64         `db:"created_at"`
}

type Role struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Grade struct {
	Id     int64  `db:"id"`
	Name   string `db:"name"`
	Points int64  `db:"points"`
}
