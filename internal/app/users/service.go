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
	RegisterUser(ctx context.Context, u dto.IntranetUserData) (user dto.User, err error)
	ListIntranetUsers(ctx context.Context, reqData dto.GetUserListReq) (data []dto.IntranetUserData, err error)
	ListUsers(ctx context.Context, reqData dto.UserListReq) (resp dto.UserListWithMetadata, err error)
}

func NewService(userRepo repository.UserStorer) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (us *service) ValidatePeerly(ctx context.Context, authToken string) (data dto.ValidateResp, err error) {
	client := &http.Client{}
	validationReq, err := http.NewRequest(http.MethodPost, config.IntranetBaseUrl()+constants.PeerlyValidationPath, nil)
	if err != nil {
		logger.Errorf("error in creating new validation request err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}
	validationReq.Header.Add(constants.AuthorizationHeader, authToken)
	validationReq.Header.Add(constants.ClientCode, config.IntranetClientCode())
	resp, err := client.Do(validationReq)
	if err != nil {
		logger.Errorf("error in intranet validation api. status returned: %d, err: %s", resp.StatusCode, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("error returned,  status returned: %d", resp.StatusCode)
		err = apperrors.IntranetValidationFailed
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("error in readall parsing. err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Errorf("error in unmarshal parsing. err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}

	return
}

func (us *service) GetIntranetUserData(ctx context.Context, req dto.GetIntranetUserDataReq) (data dto.IntranetUserData, err error) {

	client := &http.Client{}
	url := fmt.Sprintf("%s%s%d", config.IntranetBaseUrl(), constants.GetIntranetUserDataPath, req.UserId)
	intranetReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("error in creating new get user request. err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}

	intranetReq.Header.Add(constants.AuthorizationHeader, req.Token)
	resp, err := client.Do(intranetReq)
	if err != nil {
		logger.Errorf("error in intranet get user api. status returned: %d, err: %s  ", resp.StatusCode, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("error in intranet get user api. status returned: %d ", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("error in io.readall. err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}

	var respData dto.IntranetGetUserDataResp

	err = json.Unmarshal(body, &respData)
	if err != nil {
		logger.Errorf("error in unmarshalling data. err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}

	data = respData.Data

	return
}

func (us *service) LoginUser(ctx context.Context, u dto.IntranetUserData) (dto.LoginUserResp, error) {
	var resp dto.LoginUserResp
	resp.NewUserCreated = false

	user, err := us.RegisterUser(ctx, u)
	if err != nil && err != apperrors.RepeatedUser {
		return resp, err
	}

	if err == nil {
		resp.NewUserCreated = true
	}

	//sync user data
	syncNeeded, dataToBeUpdated, err := us.syncData(ctx, u, user)
	if err != nil {
		logger.Error(err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}
	if syncNeeded {

		err = us.userRepo.SyncData(ctx, dataToBeUpdated)
		if err != nil {
			logger.Error(err.Error())
			err = apperrors.InternalServerError
			return resp, err
		}

		dbResp, err := us.userRepo.GetUserByEmail(ctx, u.Email)
		if err == apperrors.InternalServerError {
			return resp, err
		}

		user = mapUserDbToService(dbResp)

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
		logger.Errorf("error generating authtoken. err: %s", err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}

	resp.User = user
	resp.AuthToken = tokenString

	return resp, nil

}

func (us *service) RegisterUser(ctx context.Context, u dto.IntranetUserData) (user dto.User, err error) {

	dbUser, err := us.userRepo.GetUserByEmail(ctx, u.Email)
	if err == apperrors.InternalServerError || err == nil {
		user = mapUserDbToService(dbUser)
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
		logger.Error(err.Error())
		err = apperrors.InternalServerError
		return
	}
	reward_quota_balance := grade.Points * reward_multiplier

	//get role by name
	roleId, err := us.userRepo.GetRoleByName(ctx, constants.UserRole)
	if err != nil {
		logger.Error(err.Error())
		err = apperrors.InternalServerError
		return
	}

	svcData := mapIntranetUserDataToSvcUser(u)

	svcData.GradeId = grade.Id
	svcData.RewardQuotaBalance = reward_quota_balance
	svcData.RoleId = roleId

	//register user
	dbResp, err := us.userRepo.CreateNewUser(ctx, svcData)
	if err != nil {
		logger.Error(err.Error())
		err = apperrors.InternalServerError
		return
	}

	user = mapUserDbToService(dbResp)

	return
}

func (us *service) ListIntranetUsers(ctx context.Context, reqData dto.GetUserListReq) (data []dto.IntranetUserData, err error) {
	client := &http.Client{}
	url := config.IntranetBaseUrl() + fmt.Sprintf(constants.ListIntranetUsersPath, reqData.Page, constants.DefaultPageSize)
	intranetReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Errorf("error in creating new intranet user list request. err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}

	intranetReq.Header.Add(constants.AuthorizationHeader, reqData.AuthToken)
	resp, err := client.Do(intranetReq)
	if err != nil {
		logger.Errorf("error in intranet get user api. status returned: %d, err: %s ", resp.StatusCode, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("erro in intranet user list request. status returned: %d", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	defer resp.Body.Close()

	var respData dto.GetUserListRespData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("error in io.readall, err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
	}

	err = json.Unmarshal(body, &respData)
	if err != nil {
		logger.Errorf("error in unmarshalling data, err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}

	data = respData.Data
	return
}

func (us *service) ListUsers(ctx context.Context, reqData dto.UserListReq) (resp dto.UserListWithMetadata, err error) {

	var names []string
	for _, data := range reqData.Name {
		names = append(names, strings.ToLower(data))
	}

	reqData.Name = names

	dbResp, totalCount, err := us.userRepo.ListUsers(ctx, reqData)
	if err != nil {
		logger.Errorf(err.Error())
		err = apperrors.InternalServerError
		return
	}

	var users []dto.UserListResp

	for _, dbUser := range dbResp {
		user := mapDbUserToUserListResp(dbUser)
		users = append(users, user)
	}

	resp.UserList = users
	resp.MetaData.TotalCount = totalCount
	resp.MetaData.CurrentPage = reqData.Page
	resp.MetaData.PageCount = reqData.PageSize

	return
}

func (us *service) syncData(ctx context.Context, intranetUserData dto.IntranetUserData, peerlyUserData dto.User) (syncNeeded bool, dataToBeUpdated dto.User, err error) {
	syncNeeded = false
	grade, err := us.userRepo.GetGradeByName(ctx, intranetUserData.EmpolyeeDetail.Grade)
	if err != nil {
		err = fmt.Errorf("error in selecting grade in syncData err: %w", err)
		return
	}

	if intranetUserData.PublicProfile.FirstName != peerlyUserData.FirstName || intranetUserData.PublicProfile.LastName != peerlyUserData.LastName || intranetUserData.PublicProfile.ProfileImgUrl != peerlyUserData.ProfileImgUrl || intranetUserData.EmpolyeeDetail.Designation.Name != peerlyUserData.Designation || grade.Id != peerlyUserData.GradeId {
		syncNeeded = true
		dataToBeUpdated.FirstName = intranetUserData.PublicProfile.FirstName
		dataToBeUpdated.LastName = intranetUserData.PublicProfile.LastName
		dataToBeUpdated.ProfileImgUrl = intranetUserData.PublicProfile.ProfileImgUrl
		dataToBeUpdated.Designation = intranetUserData.EmpolyeeDetail.Designation.Name
		dataToBeUpdated.GradeId = grade.Id
		dataToBeUpdated.Email = intranetUserData.Email
	}
	return
}

func mapUserDbToService(dbStruct repository.User) (svcStruct dto.User) {
	svcStruct.Id = dbStruct.Id
	svcStruct.EmployeeId = dbStruct.EmployeeId
	svcStruct.FirstName = dbStruct.FirstName
	svcStruct.LastName = dbStruct.LastName
	svcStruct.Email = dbStruct.Email
	svcStruct.ProfileImgUrl = dbStruct.ProfileImageURL.String
	svcStruct.RoleId = dbStruct.RoleID
	svcStruct.RewardQuotaBalance = dbStruct.RewardsQuotaBalance
	svcStruct.Designation = dbStruct.Designation
	svcStruct.GradeId = dbStruct.GradeId
	svcStruct.CreatedAt = dbStruct.CreatedAt

	return svcStruct
}

func mapIntranetUserDataToSvcUser(intranetData dto.IntranetUserData) (svcData dto.User) {
	svcData.Email = intranetData.Email
	svcData.EmployeeId = intranetData.EmpolyeeDetail.EmployeeId
	svcData.ProfileImgUrl = intranetData.PublicProfile.ProfileImgUrl
	svcData.FirstName = intranetData.PublicProfile.FirstName
	svcData.LastName = intranetData.PublicProfile.LastName
	svcData.Designation = intranetData.EmpolyeeDetail.Designation.Name
	return svcData
}

func mapDbUserToUserListResp(dbStruct repository.User) (svcData dto.UserListResp) {
	svcData.Id = dbStruct.Id
	svcData.FirstName = dbStruct.FirstName
	svcData.LastName = dbStruct.LastName
	svcData.Email = dbStruct.Email
	return svcData
}
