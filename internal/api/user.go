package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/joshsoftware/peerly-backend/internal/api/validation"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/sirupsen/logrus"
)

func loginUser(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		notificationToken := req.URL.Query().Get("notification_token")
		authToken := req.Header.Get(constants.IntranetAuth)
		if authToken == "" {
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err)
			return
		}

		ctx := req.Context()

		validateResp, err := userSvc.ValidatePeerly(ctx, authToken)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		reqData := dto.GetIntranetUserDataReq{
			Token:  validateResp.Data.JwtToken,
			UserId: validateResp.Data.UserId,
		}

		user, err := userSvc.GetIntranetUserData(ctx, reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		err = validation.GetIntranetUserDataValidation(user)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		user.NotificationToken = notificationToken
		resp, err := userSvc.LoginUser(ctx, user)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "Login successful", resp)

	}
}

func loginAdmin(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var reqData dto.AdminLoginReq
		err := json.NewDecoder(req.Body).Decode(&reqData)
		if err != nil {
			logger.Errorf("error while decoding request data. err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}
		resp, err := userSvc.AdminLogin(req.Context(), reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Login successful", resp)
	}
}

func listIntranetUsersHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		authToken := req.Header.Get(constants.IntranetAuth)
		if authToken == "" {
			err := apperrors.InvalidAuthToken
			dto.ErrorRepsonse(rw, err)
			return
		}

		page := req.URL.Query().Get("page")
		if page == "" {
			logger.Error("page query parameter is required")
			err := apperrors.PageParamNotFound
			dto.ErrorRepsonse(rw, err)
			return
		}

		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			logger.Errorf("error page string to int64 conversion. err:%s ", err.Error())
			err = apperrors.InternalServerError
			dto.ErrorRepsonse(rw, err)
			return
		}

		ctx := req.Context()

		validateResp, err := userSvc.ValidatePeerly(ctx, authToken)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		reqData := dto.GetUserListReq{
			AuthToken: validateResp.Data.JwtToken,
			Page:      pageInt,
		}

		usersData, err := userSvc.ListIntranetUsers(ctx, reqData)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
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
			logger.Errorf("error while decoding request data. err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		err = validation.GetIntranetUserDataValidation(user)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		ctx := req.Context()

		resp, err := userSvc.RegisterUser(ctx, user)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, http.StatusOK, "User registered successfully", resp)
	}
}

func listUsersHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("page")
		if page == "" {
			err := apperrors.PageParamNotFound
			dto.ErrorRepsonse(rw, err)
			return
		}
		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			logger.Errorf("error in page  string to int64 conversion, err:%s", err.Error())
			err = apperrors.InternalServerError
			dto.ErrorRepsonse(rw, err)
		}
		if pageInt <= 0 {
			err := apperrors.InvalidPage
			dto.ErrorRepsonse(rw, err)
			return
		}

		perPage := req.URL.Query().Get("page_size")
		var perPageInt int64
		if perPage == "" {
			perPageInt = constants.DefaultPageSize
		} else {
			perPageInt, err = strconv.ParseInt(perPage, 10, 64)
			if err != nil {
				logger.Errorf("error in page size string to int64 conversion, err:%s", err.Error())
				err = apperrors.InternalServerError
				dto.ErrorRepsonse(rw, err)
			}
		}
		if perPageInt <= 0 {
			err := apperrors.InvalidPageSize
			dto.ErrorRepsonse(rw, err)
			return
		}
		names := strings.Split(req.URL.Query().Get("name"), " ")
		userListReq := dto.ListUsersReq{
			Name:     names,
			Page:     pageInt,
			PageSize: perPageInt,
		}
		resp, err := userSvc.ListUsers(req.Context(), userListReq)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Peerly users listed", resp)
	}
}

func getActiveUserListHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		resp, err := userSvc.GetActiveUserList(req.Context())
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, http.StatusOK, "Active Users list", resp)
	}
}
func getUserByIdHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		resp, err := userSvc.GetUserById(req.Context())
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, 200, "User fetched successfully", resp)

	}
}

func getTop10UserHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		resp, err := userSvc.GetTop10Users(req.Context())
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "Top 10 users fetched successfully", resp)
	}
}

func adminNotificationHandler(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var msg notification.Message
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			logger.Errorf("error while decoding request data. err: %s", err.Error())
			err = apperrors.JSONParsingErrorReq
			dto.ErrorRepsonse(rw, err)
			return
		}

		err = userSvc.NotificationByAdmin(req.Context(), msg, 214, false)
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		dto.SuccessRepsonse(rw, 200, "Notification sent successfully", nil)
	}
}

func downloadExcelReport(userSvc user.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		tempFileName, err := userSvc.DownloadExcel(req.Context())
		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}

		http.ServeFile(rw, req, tempFileName)

		if err != nil {
			dto.ErrorRepsonse(rw, err)
			return
		}
		dto.SuccessRepsonse(rw, 200, "Excel downloaded successfully", nil)
	}
}
