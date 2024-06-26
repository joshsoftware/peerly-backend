package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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
}

func NewService(userRepo repository.UserStorer) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (cs *service) ValidatePeerly(ctx context.Context, authToken string) (data dto.ValidateResp, err error) {
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
		logger.WithField("err", "err").Error("Status returned ", resp.StatusCode)
		err = apperrors.IntranetValidationFailed
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = apperrors.JSONParsingErrorResp
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = apperrors.JSONParsingErrorResp
		return
	}

	return
}

func (cs *service) GetIntranetUserData(ctx context.Context, req dto.GetIntranetUserDataReq) (data dto.IntranetUserData, err error) {

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

func (cs *service) LoginUser(ctx context.Context, u dto.IntranetUserData) (dto.LoginUserResp, error) {
	var resp dto.LoginUserResp
	resp.NewUserCreated = false
	user, err := cs.userRepo.GetUserByEmail(ctx, u.Email)
	if err == apperrors.InternalServerError {
		return resp, err
	}

	if err == apperrors.UserNotFound {

		//get grade id
		gradeId, err := cs.userRepo.GetGradeByName(ctx, u.EmpolyeeDetail.Grade)
		if err != nil {
			return resp, err
		}

		//reward_quota_balance from organization config
		reward_quota_balance, err := cs.userRepo.GetRewardOuotaDefault(ctx)
		if err != nil {
			err = apperrors.InternalServerError
			return resp, err
		}

		//get role by name
		roleId, err := cs.userRepo.GetRoleByName(ctx, constants.UserRole)
		if err != nil {
			err = apperrors.InternalServerError
			return resp, err
		}

		var userData dto.RegisterUser
		userData.User = u
		userData.GradeId = gradeId
		userData.RewardQuotaBalance = reward_quota_balance
		userData.RoleId = roleId

		//register user
		user, err = cs.userRepo.CreateNewUser(ctx, userData)
		if err != nil {
			err = apperrors.InternalServerError
			return resp, err
		}

		resp.NewUserCreated = true
	}

	//login user

	expirationTime := time.Now().Add(time.Hour * 5)

	claims := &dto.Claims{
		Id:     user.Id,
		RoleId: user.RoleId,
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
