package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joshsoftware/peerly-backend/internal/api/validation"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func loginUser(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		authToken := req.Header.Get(constants.IntranetAuth)
		if authToken == "" {
			err := apperrors.InvalidAuthToken
			apperrors.ErrorResp(rw, err)
			return
		}

		validateResp, err := userSvc.ValidatePeerly(req.Context(), authToken)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		reqData := dto.GetIntranetUserDataReq{
			Token:  validateResp.Data.JwtToken,
			UserId: validateResp.Data.UserId,
		}

		user, err := userSvc.GetIntranetUserData(req.Context(), reqData)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		err = validation.GetIntranetUserDataValidation(user)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		resp, err := userSvc.LoginUser(req.Context(), user)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: resp})

	}
}

func getIntranetUserListHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header.Get(constants.IntranetAuth)
		if authToken == "" {
			err := apperrors.InvalidAuthToken
			apperrors.ErrorResp(rw, err)
			return
		}

		page := req.URL.Query().Get("page")
		if page == "" {
			err := apperrors.PageParamNotFound
			apperrors.ErrorResp(rw, err)
			return
		}
		pageInt, _ := strconv.Atoi(page)

		validateResp, err := userSvc.ValidatePeerly(req.Context(), authToken)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		reqData := dto.GetUserListReq{
			AuthToken: validateResp.Data.JwtToken,
			Page:      pageInt,
		}

		usersData, err := userSvc.GetUserListIntranet(req.Context(), reqData)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}

		// var users []dto.GetUserListResp

		// for i := 0; i < len(usersData); i++ {
		// 	var user dto.GetUserListResp
		// 	user.EmployeeId = usersData[i].EmpolyeeDetail.EmployeeId
		// 	user.FirstName = usersData[i].PublicProfile.FirstName
		// 	user.LastName = usersData[i].PublicProfile.LastName
		// 	user.Email = usersData[i].Email
		// 	user.Grade = usersData[i].EmpolyeeDetail.Grade
		// 	user.Designation = usersData[i].EmpolyeeDetail.Designation.Name
		// 	user.ProfileImg = usersData[i].PublicProfile.ProfileImgUrl
		// 	users = append(users, user)
		// }

		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: usersData})
	}
}

func registerUser(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var user dto.IntranetUserData
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding request data")
			err = apperrors.JSONParsingErrorReq
			apperrors.ErrorResp(rw, err)
			return
		}
		resp, err := userSvc.RegisterUser(req.Context(), user)
		if err != nil {
			apperrors.ErrorResp(rw, err)
			return
		}
		dto.Repsonse(rw, http.StatusOK, dto.SuccessResponse{Data: resp})
	}
}
