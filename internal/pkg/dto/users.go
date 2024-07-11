package dto

import (
	"github.com/dgrijalva/jwt-go"
)

type LoginReq struct {
	Authtoken string `json:"authtoken"`
}

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
	Id             int            `json:"id"`
	Email          string         `json:"email"`
	PublicProfile  PublicProfile  `json:"public_profile"`
	EmpolyeeDetail EmployeeDetail `json:"employee_detail"`
}

type IntranetGetUserDataResp struct {
	Data IntranetUserData `json:"data"`
}

type GetUserResp struct {
	Id                 int64  `json:"id" db:"id"`
	EmployeeId         string `json:"employee_id" db:"employee_id"`
	FirstName          string `json:"first_name" db:"first_name"`
	LastName           string `json:"last_name" db:"last_name"`
	Email              string `json:"email" db:"email"`
	ProfileImgUrl      string `json:"profile_image_url" db:"profile_image_url"`
	RoleId             int    `json:"role_id" db:"role_id"`
	RewardQuotaBalance int    `json:"reward_quota_balance" db:"reward_quota_balance"`
	Designation        string `json:"designation" db:"designation"`
	GradeId            int    `json:"grade_id" db:"grade_id"`
	Grade              string `json:"grade" db:"name"`
	CreatedAt          int64  `db:"created_at" json:"created_at"`
}

type RegisterUser struct {
	User               IntranetUserData
	RoleId             int `json:"role_id" db:"role_id"`
	RewardQuotaBalance int `json:"reward_quota_balance" db:"reward_quota_balance"`
	GradeId            int `json:"grade_id" db:"grade_id"`
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
	Id   int64
	Role string
	jwt.StandardClaims
}

type LoginUserResp struct {
	User           GetUserResp
	NewUserCreated bool
	AuthToken      string
}

type UpdateUserData struct {
	EmployeeId    string `json:"employee_id" db:"employee_id"`
	FirstName     string `json:"first_name" db:"first_name"`
	LastName      string `json:"last_name" db:"last_name"`
	ProfileImgUrl string `json:"profile_image_url" db:"profile_image_url"`
	Designation   string `json:"designation" db:"designation"`
	Grade         string `json:"grade" db:"name"`
	GradeId       int    `json:"grade_id" db:"grade_id"`
	Email         string `json:"email" db:"email"`
}

type GetUserListReq struct {
	AuthToken string
	Page      int
}

//	type GetUserListRespDB struct {
//		Id          int64          `json:"id" db:"id"`
//		EmployeeId  string         `json:"employee_id" db:"employee_id"`
//		Email       string         `json:"email" db:"email"`
//		FirstName   string         `json:"first_name" db:"first_name"`
//		LastName    string         `json:"last_name" db:"last_name"`
//		Grade       string         `json:"grade" db:"name"`
//		Designation string         `json:"designation" db:"designation"`
//		ProfileImg  sql.NullString `json:"profile_image_url" db:"profile_image_url"`
//	}
type GetUserListResp struct {
	Id        int64  `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}

type UserListMetadata struct {
	TotalCount  int64 `json:"total_count"`
	CurrentPage int64 `json:"current_page"`
	PageCount   int64 `json:"page_size"`
}

type UserListWithMetadata struct {
	UserList []GetUserListResp `json:"user_list"`
	MetaData UserListMetadata  `json:"metadata"`
}

type GetUserListRespData struct {
	Data []IntranetUserData `json:"data"`
}

type UserListReq struct {
	Page    int64
	PerPage int64
	Name    []string
}
