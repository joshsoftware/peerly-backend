package dto

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
)

type PublicProfile struct {
	ProfileImgUrl string `json:"profile_image_url"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type Designation struct {
	Name string `json:"name"`
}

type EmployeeDetail struct {
	EmployeeId  string      `json:"employee_id"`
	Designation Designation `json:"designation"`
	Grade       string      `json:"grade"`
}
type IntranetUserData struct {
	Id             int64          `json:"id"`
	Email          string         `json:"email"`
	PublicProfile  PublicProfile  `json:"public_profile"`
	EmpolyeeDetail EmployeeDetail `json:"employee_detail"`
}

type IntranetGetUserDataResp struct {
	Data IntranetUserData `json:"data"`
}

type User struct {
	Id                 int64  `json:"id"`
	EmployeeId         string `json:"employee_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	ProfileImgUrl      string `json:"profile_image_url"`
	RoleId             int64  `json:"role_id"`
	RewardQuotaBalance int64  `json:"reward_quota_balance"`
	Designation        string `json:"designation"`
	GradeId            int64  `json:"grade_id"`
	CreatedAt          int64  `json:"created_at"`
}

type ValidateResp struct {
	Data IntranetValidateApiData `json:"data"`
}

type IntranetValidateApiData struct {
	JwtToken string `json:"jwt_token"`
	UserId   int64  `json:"user_id"`
}

type GetIntranetUserDataReq struct {
	Token  string
	UserId int64
}

type Claims struct {
	Id   int64
	Role string
	jwt.StandardClaims
}

type LoginUserResp struct {
	User           User
	NewUserCreated bool
	AuthToken      string
}
type GetUserListReq struct {
	AuthToken string
	Page      int64
}

type UserListResp struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserListMetadata struct {
	TotalCount  int64 `json:"total_count"`
	CurrentPage int64 `json:"current_page"`
	PageCount   int64 `json:"page_size"`
}

type UserListWithMetadata struct {
	UserList []UserListResp   `json:"user_list"`
	MetaData UserListMetadata `json:"metadata"`
}

type GetUserListRespData struct {
	Data []IntranetUserData `json:"data"`
}
type UserListReq struct {
	Page    int64
	PerPage int64
	Name    []string
}

type ActiveUser struct {
	ID                 int    `json:"id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	ProfileImageURL    string `json:"profile_image_url"`
	BadgeName          string `json:"badge_name"`
	AppreciationPoints int    `json:"appreciation_points"`
}
type Top10User struct {
	ID                 int    `json:"id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	ProfileImageURL    string `json:"profile_image_url"`
	BadgeName          string `json:"badge_name"`
	AppreciationPoints int    `json:"appreciation_points"`
}
type GetUserByIdReq struct {
	UserId          int64 `json:"user_id" db:"id"`
	QuaterTimeStamp int64
}

type GetUserByIdDbResp struct {
	UserId             int64          `json:"user_id" db:"id"`
	FirstName          string         `json:"first_name" db:"first_name"`
	LastName           string         `json:"last_name" db:"last_name"`
	Email              string         `json:"email" db:"email"`
	ProfileImgUrl      string         `json:"profile_image_url" db:"profile_image_url"`
	Designation        string         `json:"designation" db:"designation"`
	RewardQuotaBalance int64          `json:"reward_quota_balance" db:"reward_quota_balance"`
	GradeId            int64          `json:"grade_id" db:"grade_id"`
	EmployeeId         string         `json:"employee_id" db:"employee_id"`
	TotalPoints        int64          `json:"total_points" db:"total_points"`
	Badge              sql.NullString `json:"badge" db:"name"`
	BadgeCreatedAt     sql.NullInt64  `json:"badge_created_at" db:"badge_created_at"`
}

type GetUserByIdResp struct {
	UserId             int64  `json:"user_id" db:"id"`
	FirstName          string `json:"first_name" db:"first_name"`
	LastName           string `json:"last_name" db:"last_name"`
	Email              string `json:"email" db:"email"`
	ProfileImgUrl      string `json:"profile_image_url" db:"profile_image_url"`
	Designation        string `json:"designation" db:"designation"`
	RewardQuotaBalance int64  `json:"reward_quota_balance" db:"reward_quota_balance"`
	TotalRewardQuota   int64  `json:"total_reward_quota"`
	RefilDate          int64  `json:"refil_date"`
	GradeId            int64  `json:"grade_id" db:"grade_id"`
	EmployeeId         string `json:"employee_id" db:"employee_id"`
	TotalPoints        int64  `json:"total_points" db:"total_points"`
	Badge              string `json:"badge" db:"name"`
	BadgeCreatedAt     int64  `json:"badge_created_at" db:"badge_created_at"`
}
type AdminLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
