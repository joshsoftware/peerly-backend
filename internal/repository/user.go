package repository

import (
	"context"
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
)

type UserStorer interface {
	RepositoryTransaction

	GetUserByEmail(ctx context.Context, email string) (user User, err error)
	GetRoleByName(ctx context.Context, name string) (roleId int64, err error)
	CreateNewUser(ctx context.Context, user dto.User) (resp User, err error)
	GetGradeByName(ctx context.Context, name string) (grade Grade, err error)
	GetRewardMultiplier(ctx context.Context) (value int64, err error)
	SyncData(ctx context.Context, updateData dto.User) (err error)
	ListUsers(ctx context.Context, reqData dto.ListUsersReq) (resp []User, count int64, err error)

	UpdateRewardQuota(ctx context.Context, tx Transaction) (err error)
	GetActiveUserList(ctx context.Context, tx Transaction) (activeUsers []ActiveUser, err error)
	GetUserById(ctx context.Context, reqData dto.GetUserByIdReq) (user dto.GetUserByIdResp, err error)
	GetTop10Users(ctx context.Context, quarterTimestamp int64) (users []Top10Users, err error)
	GetGradeById(ctx context.Context, id int64) (grade Grade, err error)
	GetAdmin(ctx context.Context, email string) (user User, err error)
	AddDeviceToken(ctx context.Context, userID int64, deviceToken string) (err error)
	ListDeviceTokensByUserID(ctx context.Context, userID int64) (notificationTokens []string, err error)
}

// User - basic struct representing a User
type User struct {
	Id                  int64          `db:"id"`
	EmployeeId          string         `db:"employee_id"`
	FirstName           string         `db:"first_name"`
	LastName            string         `db:"last_name"`
	Email               string         `db:"email"`
	Password            sql.NullString `db:"password"`
	ProfileImageURL     sql.NullString `db:"profile_image_url"`
	GradeId             int64          `db:"grade_id"`
	Designation         string         `db:"designation"`
	RoleID              int64          `db:"role_id"`
	RewardsQuotaBalance int64          `db:"reward_quota_balance"`
	Status              int64          `db:"status"`
	SoftDelete          bool           `db:"soft_delete"`
	SoftDeleteBy        sql.NullInt64  `db:"soft_delete_by"`
	SoftDeleteOn        sql.NullTime   `db:"soft_delete_on"`
	CreatedAt           int64          `db:"created_at"`
}

type Role struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type ActiveUser struct {
	ID                 int            `db:"id"`
	FirstName          string         `db:"first_name"`
	LastName           string         `db:"last_name"`
	ProfileImageURL    sql.NullString `db:"profile_image_url"`
	BadgeName          sql.NullString `db:"badge_name"`
	AppreciationPoints int            `db:"appreciation_points"`
}

type Top10Users struct {
	ID                 int            `db:"id"`
	FirstName          string         `db:"first_name"`
	LastName           string         `db:"last_name"`
	ProfileImageURL    sql.NullString `db:"profile_image_url"`
	BadgeName          sql.NullString `db:"name"`
	AppreciationPoints int            `db:"ap"`
}

type UserBadgeDetails struct {
	ID          int64          `db:"id"`
	FirstName   string         `db:"first_name"`
	LastName    string         `db:"last_name"`
	Email       string         `db:"email"`
	BadgeID     int8           `db:"badge_id"`
	BadgeName   sql.NullString `db:"badge_name"`
	BadgePoints int32          `db:"badge_points"`
}
