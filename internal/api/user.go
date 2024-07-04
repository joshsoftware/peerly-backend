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
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		validateResp, err := userSvc.ValidatePeerly(req.Context(), authToken)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		reqData := dto.GetIntranetUserDataReq{
			Token:  validateResp.Data.JwtToken,
			UserId: validateResp.Data.UserId,
		}

		user, err := userSvc.GetIntranetUserData(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		err = validation.GetIntranetUserDataValidation(user)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		resp, err := userSvc.LoginUser(req.Context(), user)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Login successful", resp)

	}
}

func getIntranetUserListHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header.Get(constants.IntranetAuth)
		if authToken == "" {
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		page := req.URL.Query().Get("page")
		if page == "" {
			err := apperrors.PageParamNotFound
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}
		pageInt, _ := strconv.Atoi(page)

		validateResp, err := userSvc.ValidatePeerly(req.Context(), authToken)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		reqData := dto.GetUserListReq{
			AuthToken: validateResp.Data.JwtToken,
			Page:      pageInt,
		}

		usersData, err := userSvc.GetUserListIntranet(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Intranet users listed", usersData)
	}
}

func registerUser(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var user dto.IntranetUserData
		err := json.NewDecoder(req.Body).Decode(&user)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error while decoding request data")
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}
		resp, err := userSvc.RegisterUser(req.Context(), user)
		if err != nil {
			dto.ErrorRepsonse(rw, apperrors.GetHTTPStatusCode(err), err.Error(), nil)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "User registered successfully", resp)
	}
}
