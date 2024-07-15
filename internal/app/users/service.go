package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type service struct {
	userRepo repository.UserStorer
}

type Service interface {
	ValidatePeerly(ctx context.Context, authToken string) (data dto.ValidateResp, err error)
	GetIntranetUserData(ctx context.Context, req dto.GetIntranetUserDataReq) (data dto.IntranetUserData, err error)
	LoginUser(ctx context.Context, u dto.IntranetUserData) (dto.LoginUserResp, error)
	RegisterUser(ctx context.Context, u dto.IntranetUserData) (user dto.GetUserResp, err error)
	GetUserListIntranet(ctx context.Context, reqData dto.GetUserListReq) (data []dto.IntranetUserData, err error)
	UpdateRewardQuota(ctx context.Context)(err error)
	GetActiveUserList(ctx context.Context) ([]dto.ActiveUser,error)
	GetUserList(ctx context.Context, reqData dto.UserListReq) (resp dto.UserListWithTotalCount, err error)
	GetUserById(ctx context.Context) (user dto.GetUserByIdResp, err error)
}

func NewService(userRepo repository.UserStorer) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (us *service) ValidatePeerly(ctx context.Context, authToken string) (data dto.ValidateResp, err error) {
	client := &http.Client{}
	validationReq, err := http.NewRequest("POST", "https://pg-stage-intranet.joshsoftware.com/api/peerly/v1/sessions/login", nil)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	validationReq.Header.Add(constants.AuthorizationHeader, authToken)
	validationReq.Header.Add(constants.ClientCode, config.IntranetClientCode())
	resp, err := client.Do(validationReq)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in intranet validation api. Status returned:  ", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error("Status returned ", resp.StatusCode)
		err = apperrors.IntranetValidationFailed
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in readall parsing")
		err = apperrors.JSONParsingErrorResp
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in unmarshal parsing")
		err = apperrors.JSONParsingErrorResp
		return
	}

	return
}

func (us *service) GetIntranetUserData(ctx context.Context, req dto.GetIntranetUserDataReq) (data dto.IntranetUserData, err error) {

	client := &http.Client{}
	url := fmt.Sprintf("https://pg-stage-intranet.joshsoftware.com/api/peerly/v1/users/%d", req.UserId)
	intranetReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}

	intranetReq.Header.Add(constants.AuthorizationHeader, req.Token)
	resp, err := client.Do(intranetReq)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in intranet get user api. Status returned:  ", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.WithField("err", "err").Error("Status returned ", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in io.readall")
		err = apperrors.JSONParsingErrorResp
	}

	var respData dto.IntranetGetUserDataResp

	err = json.Unmarshal(body, &respData)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in unmarshalling data")
		err = apperrors.JSONParsingErrorResp
		return
	}

	data = respData.Data

	return
}

func (us *service) LoginUser(ctx context.Context, u dto.IntranetUserData) (dto.LoginUserResp, error) {
	var resp dto.LoginUserResp
	resp.NewUserCreated = false
	// user, err := us.userRepo.GetUserByEmail(ctx, u.Email)
	// if err == apperrors.InternalServerError {
	// 	return resp, err
	// }

	// if err == apperrors.UserNotFound {

	// 	user, err = us.RegisterUser(ctx, u)
	// 	if err != nil {
	// 		return resp, err
	// 	}

	// 	resp.NewUserCreated = true
	// }

	user, err := us.RegisterUser(ctx, u)
	if err != nil && err != apperrors.RepeatedUser {
		return resp, err
	}

	if err == nil {
		resp.NewUserCreated = true
	}

	//sync user data
	syncNeeded, dataToBeUpdated := syncData(u, user)
	if syncNeeded {

		grade, err := us.userRepo.GetGradeByName(ctx, dataToBeUpdated.Grade)
		if err != nil {
			return resp, err
		}
		dataToBeUpdated.GradeId = grade.Id

		err = us.userRepo.SyncData(ctx, dataToBeUpdated)
		if err != nil {
			err = apperrors.InternalServerError
			return resp, err
		}
		user, err = us.userRepo.GetUserByEmail(ctx, u.Email)
		if err == apperrors.InternalServerError {
			return resp, err
		}

	}

	//login user

	expirationTime := time.Now().Add(time.Hour * time.Duration(config.JWTExpiryDurationHours()))

	claims := &dto.Claims{
		Id:   user.Id,
		Role: constants.UserRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	var jwtKey = config.JWTKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error generating authtoken")
		err = apperrors.InternalServerError
		return resp, err
	}

	resp.User = user
	resp.AuthToken = tokenString

	return resp, nil

}

func (us *service) RegisterUser(ctx context.Context, u dto.IntranetUserData) (user dto.GetUserResp, err error) {

	user, err = us.userRepo.GetUserByEmail(ctx, u.Email)
	if err == apperrors.InternalServerError || err == nil {
		err = apperrors.RepeatedUser
		return
	}

	//get grade id
	grade, err := us.userRepo.GetGradeByName(ctx, u.EmpolyeeDetail.Grade)
	if err != nil {
		return
	}

	//reward_multiplier from organization config
	reward_multiplier, err := us.userRepo.GetRewardMultiplier(ctx)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	reward_quota_balance := grade.Points * reward_multiplier

	//get role by name
	roleId, err := us.userRepo.GetRoleByName(ctx, constants.UserRole)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}

	var userData dto.RegisterUser
	userData.User = u
	userData.GradeId = grade.Id
	userData.RewardQuotaBalance = reward_quota_balance
	userData.RoleId = roleId

	//register user
	user, err = us.userRepo.CreateNewUser(ctx, userData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}

	return
}

func (us *service) GetUserListIntranet(ctx context.Context, reqData dto.GetUserListReq) (data []dto.IntranetUserData, err error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://pg-stage-intranet.joshsoftware.com/api/peerly/v1/users?page=%d&per_page=%d", reqData.Page, constants.PerPage)
	intranetReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}

	intranetReq.Header.Add(constants.AuthorizationHeader, reqData.AuthToken)
	resp, err := client.Do(intranetReq)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in intranet get user api. Status returned:  ", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.WithField("err", "err").Error("Status returned ", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	defer resp.Body.Close()

	var respData dto.GetUserListRespData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in io.readall")
		err = apperrors.JSONParsingErrorResp
	}

	err = json.Unmarshal(body, &respData)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error in unmarshalling data")
		err = apperrors.JSONParsingErrorResp
		return
	}

	data = respData.Data
	return
}

func (us *service) GetUserList(ctx context.Context, reqData dto.UserListReq) (resp dto.UserListWithTotalCount, err error) {

	var names []string
	for _, data := range reqData.Name {
		names = append(names, strings.ToLower(data))
	}

	reqData.Name = names

	totalCount, err := us.userRepo.GetTotalUserCount(ctx, reqData)
	if err != nil {
		return
	}

	users, err := us.userRepo.GetUserList(ctx, reqData)
	if err != nil {
		return
	}

	resp.UserList = users
	resp.TotalCount = totalCount

	return

}

func (us *service) GetUserById(ctx context.Context) (user dto.GetUserByIdResp, err error) {

	id := ctx.Value(constants.UserId)
	fmt.Printf("userId: %T", id)
	userId, ok := id.(int64)
	if !ok {
		logger.Error("Error in typecasting user id")
		err = apperrors.InternalServerError
		return
	}

	quaterTimeStamp := GetQuarterStartUnixTime()

	reqData := dto.GetUserByIdReq{
		UserId:          userId,
		QuaterTimeStamp: quaterTimeStamp,
	}

	user, err = us.userRepo.GetUserById(ctx, reqData)
	if err != nil {
		return
	}

	return
}

func (us *service) GetActiveUserList(ctx context.Context) ([]dto.ActiveUser,error){
	activeUserDb,err := us.userRepo.GetActiveUserList(ctx,nil)
	if err != nil{
		return []dto.ActiveUser{},err
	}
	res := make([]dto.ActiveUser,0)
	for _,activerUser := range activeUserDb{
		actUsr := MapActiveUserDbtoDto(activerUser)
		res = append(res, actUsr)
	}
	return res,nil
}
func (us *service) UpdateRewardQuota(ctx context.Context)(error){
	err := us.userRepo.UpdateRewardQuota(ctx,nil)
	return err
}
func GetQuarterStartUnixTime() int64 {
	// Example function to get the Unix timestamp of the start of the quarter
	now := time.Now()
	quarterStart := time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
	return quarterStart.Unix() * 1000 // convert to milliseconds
}
