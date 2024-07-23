package dto

import (
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

type GetUserListRespData struct {
	Data []IntranetUserData `json:"data"`
}
