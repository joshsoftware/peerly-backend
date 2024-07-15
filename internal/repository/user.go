package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type UserStorer interface {
	GetUserByEmail(ctx context.Context, email string) (user User, err error)
	GetRoleByName(ctx context.Context, name string) (roleId int, err error)
	CreateNewUser(ctx context.Context, user dto.User) (resp User, err error)
	GetGradeByName(ctx context.Context, name string) (grade Grade, err error)
	GetRewardMultiplier(ctx context.Context) (value int, err error)
	SyncData(ctx context.Context, updateData dto.User) (err error)
}

// User - basic struct representing a User
type User struct {
	Id                  int           `db:"id"`
	EmployeeId          string        `db:"employee_id"`
	FirstName           string        `db:"first_name"`
	LastName            string        `db:"last_name"`
	Email               string        `db:"email"`
	ProfileImageURL     string        `db:"profile_image_url"`
	GradeId             int           `db:"grade_id"`
	Designation         string        `db:"designation"`
	RoleID              int           `db:"role_id"`
	RewardsQuotaBalance int           `db:"reward_quota_balance"`
	Status              int           `db:"status"`
	SoftDelete          bool          `db:"soft_delete"`
	SoftDeleteBy        sql.NullInt64 `db:"soft_delete_by"`
	SoftDeleteOn        sql.NullTime  `db:"soft_delete_on"`
	CreatedAt           int64         `db:"created_at"`
}

type Role struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Grade struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Points int    `db:"points"`
}
