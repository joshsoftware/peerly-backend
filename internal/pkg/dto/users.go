package dto

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type LoginReq struct {
	Authtoken string `json:"authtoken"`
}

type IntranetApiResp struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	ProfileImgUrl string `json:"profile_image_url"`
	Designation   string `json:"designation"`
	Grade         string `json:"grade"`
}

type GetUserResp struct {
	Id                 int       `json:"id" db:"id"`
	FirstName          string    `json:"first_name" db:"first_name"`
	OrgId              int       `json:"org_id" db:"org_id"`
	LastName           string    `json:"last_name" db:"last_name"`
	Email              string    `json:"email" db:"email"`
	ProfileImgUrl      string    `json:"profile_image_url" db:"profile_image_url"`
	RoleId             int       `json:"role_id" db:"role_id"`
	RewardQuotaBalance int       `json:"reward_quota_balance" db:"reward_quota_balance"`
	Designation        string    `json:"designation" db:"designation"`
	GradeId            int       `json:"grade_id" db:"grade_id"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

type RegisterUser struct {
	User               IntranetApiResp
	OrgId              int       `json:"org_id" db:"org_id"`
	RoleId             int       `json:"role_id" db:"role_id"`
	RewardQuotaBalance int       `json:"reward_quota_balance" db:"reward_quota_balance"`
	GradeId            int       `json:"grade_id" db:"grade_id"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

type ValidateResp struct {
	Data IntranetValidateApiData `json:"data"`
}

type IntranetValidateApiData struct {
	JwtToken string `json:"jwt_token"`
	UserId   int    `json:"user_id"`
}

type GetIntranetUserDataReq struct {
	Token  string
	UserId int
}

type Claims struct {
	Id     int
	RoleId int
	jwt.StandardClaims
}

type LoginUserResp struct {
	User           GetUserResp
	NewUserCreated bool
	AuthToken      string
}
