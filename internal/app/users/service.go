package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"

	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xuri/excelize/v2"
	// "github.com/joshsoftware/peerly-backend/internal/app/email"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
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
	ListUsers(ctx context.Context, reqData dto.ListUsersReq) (resp dto.ListUsersResp, err error)
	GetUserById(ctx context.Context) (user dto.GetUserByIdResp, err error)
	UpdateRewardQuota(ctx context.Context) (err error)
	GetActiveUserList(ctx context.Context) ([]dto.ActiveUser, error)
	GetTop10Users(ctx context.Context) (users []dto.Top10User, err error)
	AdminLogin(ctx context.Context, loginReq dto.AdminLoginReq) (resp dto.LoginUserResp, err error)
	// sendRewardQuotaRefillEmailToAll(ctx context.Context)
	NotificationByAdmin(ctx context.Context, notificationReq dto.AdminNotificationReq) (err error)
	AllAppreciationReport(ctx context.Context, appreciations []dto.AppreciationResponse) (tempFileName string, err error)
	ReportedAppreciationReport(ctx context.Context, appreciations []dto.ReportedAppreciation) (tempFileName string, err error)
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
		logger.Errorf(ctx, "error in creating new validation request err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}
	validationReq.Header.Add(constants.AuthorizationHeader, authToken)
	validationReq.Header.Add(constants.ClientCode, config.IntranetClientCode())
	resp, err := client.Do(validationReq)
	if err != nil {
		logger.Errorf(ctx, "error in intranet validation api. status returned: %d, err: %s", resp.StatusCode, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf(ctx, "error returned,  status returned: %d", resp.StatusCode)
		err = apperrors.IntranetValidationFailed
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(ctx, "error in readall parsing. err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Errorf(ctx, "error in unmarshal parsing. err: %s", err.Error())
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
		logger.Errorf(ctx, "error in creating new get user request. err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}

	intranetReq.Header.Add(constants.AuthorizationHeader, req.Token)
	resp, err := client.Do(intranetReq)
	if err != nil {
		logger.Errorf(ctx, "error in intranet get user api. status returned: %d, err: %s  ", resp.StatusCode, err.Error())
		logger.Errorf(ctx, "error response: %v", resp)
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf(ctx, "error in intranet get user api. status returned: %d ", resp.StatusCode)
		logger.Errorf(ctx, "error response: %v", resp)
		err = apperrors.InternalServerError
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(ctx, "error in io.readall. err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}

	var respData dto.IntranetGetUserDataResp

	err = json.Unmarshal(body, &respData)
	if err != nil {
		logger.Errorf(ctx, "error in unmarshalling data. err: %s", err.Error())
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
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}
	if syncNeeded {

		err = us.userRepo.SyncData(ctx, dataToBeUpdated)
		if err != nil {
			logger.Error(ctx, err.Error())
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
		logger.Errorf(ctx, "error generating authtoken. err: %s", err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}

	resp.User = user
	resp.AuthToken = tokenString

	err = us.userRepo.AddDeviceToken(ctx, user.Id, u.NotificationToken)
	if err != nil {
		logger.Errorf(ctx, "err in adding device token: %v", err)
	}
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
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}
	reward_quota_balance := grade.Points * reward_multiplier

	//get role by name
	roleId, err := us.userRepo.GetRoleByName(ctx, constants.UserRole)
	if err != nil {
		logger.Error(ctx, err.Error())
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
		logger.Error(ctx, err.Error())
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
		logger.Errorf(ctx, "error in creating new intranet user list request. err: %s", err.Error())
		err = apperrors.InternalServerError
		return
	}

	intranetReq.Header.Add(constants.AuthorizationHeader, reqData.AuthToken)
	resp, err := client.Do(intranetReq)
	if err != nil {
		logger.Errorf(ctx, "error in intranet get user api. status returned: %d, err: %s ", resp.StatusCode, err.Error())
		err = apperrors.InternalServerError
		return
	}
	if resp.StatusCode != http.StatusOK {
		logger.Errorf(ctx, "erro in intranet user list request. status returned: %d", resp.StatusCode)
		err = apperrors.InternalServerError
		return
	}
	defer resp.Body.Close()

	var respData dto.ListIntranetUsersRespData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(ctx, "error in io.readall, err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
	}

	err = json.Unmarshal(body, &respData)
	if err != nil {
		logger.Errorf(ctx, "error in unmarshalling data, err: %s", err.Error())
		err = apperrors.JSONParsingErrorResp
		return
	}

	data = respData.Data
	return
}

func (us *service) ListUsers(ctx context.Context, reqData dto.ListUsersReq) (resp dto.ListUsersResp, err error) {

	var names []string
	for _, data := range reqData.Name {
		if data != "" {
			names = append(names, strings.ToLower(data))
		}
	}

	reqData.Name = names

	dbResp, totalCount, err := us.userRepo.ListUsers(ctx, reqData)
	if err != nil {
		logger.Errorf(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	userId := ctx.Value(constants.UserId)

	var users []dto.UserDetails

	for _, dbUser := range dbResp {
		if reqData.Self {
			if dbUser.Id != userId {
				user := mapDbUserToUserListResp(dbUser)
				users = append(users, user)
			}
		} else {
			user := mapDbUserToUserListResp(dbUser)
			users = append(users, user)
		}
	}

	if reqData.Self && totalCount > 0 {
		totalCount = totalCount - 1
	}

	resp.UserList = users
	resp.MetaData.TotalRecords = totalCount
	resp.MetaData.CurrentPage = reqData.Page
	resp.MetaData.PageSize = reqData.PageSize
	resp.MetaData.TotalPage = int64(math.Ceil(float64(totalCount) / float64(reqData.PageSize)))

	return
}
func (us *service) AdminLogin(ctx context.Context, loginReq dto.AdminLoginReq) (resp dto.LoginUserResp, err error) {

	dbUser, err := us.userRepo.GetAdmin(ctx, loginReq.Email)
	if err != nil {
		return
	}

	if dbUser.RoleID != 2 {
		logger.Errorf(ctx, "unathorized access")
		err = apperrors.RoleUnathorized
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password.String), []byte(loginReq.Password))
	if err != nil {
		logger.Errorf(ctx, "invalid password, err: %s", err.Error())
		err = apperrors.InvalidPassword
		return
	}

	user := mapUserDbToService(dbUser)

	expirationTime := time.Now().Add(time.Hour * time.Duration(config.JWTExpiryDurationHours()))

	claims := &dto.Claims{
		Id:   user.Id,
		Role: constants.AdminRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	var jwtKey = config.JWTKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		logger.Errorf(ctx, "error generating authtoken. err: %s", err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}

	resp.User = user
	resp.AuthToken = tokenString

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

func (us *service) GetUserById(ctx context.Context) (user dto.GetUserByIdResp, err error) {

	id := ctx.Value(constants.UserId)
	fmt.Printf("userId: %T", id)
	userId, ok := id.(int64)
	if !ok {
		logger.Error(ctx, "Error in typecasting user id")
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

	grade, err := us.userRepo.GetGradeById(ctx, user.GradeId)
	if err != nil {
		logger.Errorf(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	reward_multiplier, err := us.userRepo.GetRewardMultiplier(ctx)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	total_reward_quota := grade.Points * reward_multiplier

	user.TotalRewardQuota = int64(total_reward_quota)

	now := time.Now()

	nextMonth := now.AddDate(0, 1, 0)
	firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	timestamp := firstDayOfNextMonth.Unix()

	user.RefilDate = timestamp

	return
}

func (us *service) GetActiveUserList(ctx context.Context) ([]dto.ActiveUser, error) {
	activeUserDb, err := us.userRepo.GetActiveUserList(ctx, nil)
	if err != nil {
		return []dto.ActiveUser{}, err
	}
	res := make([]dto.ActiveUser, 0)
	for _, activerUser := range activeUserDb {
		actUsr := MapActiveUserDbtoDto(activerUser)
		res = append(res, actUsr)
	}
	return res, nil
}
func (us *service) UpdateRewardQuota(ctx context.Context) error {
	err := us.userRepo.UpdateRewardQuota(ctx, nil)

	if err == nil {
		// us.sendRewardQuotaRefillEmailToAll(ctx)
	}
	return err
}
func GetQuarterStartUnixTime() int64 {
	// Example function to get the Unix timestamp of the start of the quarter
	now := time.Now()
	quarterStart := time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
	return quarterStart.Unix() * 1000 // convert to milliseconds
}

func (us *service) GetTop10Users(ctx context.Context) (users []dto.Top10User, err error) {

	quaterTimeStamp := GetQuarterStartUnixTime()
	dbUsers, err := us.userRepo.GetTop10Users(ctx, quaterTimeStamp)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	for _, dbUser := range dbUsers {
		svcUser := mapDbTop10ToSvcTop10(dbUser)
		users = append(users, svcUser)
	}

	return
}

func mapDbTop10ToSvcTop10(dbStruct repository.Top10Users) (svcStruct dto.Top10User) {
	svcStruct.ID = dbStruct.ID
	svcStruct.FirstName = dbStruct.FirstName
	svcStruct.LastName = dbStruct.LastName
	svcStruct.ProfileImageURL = dbStruct.ProfileImageURL.String
	svcStruct.BadgeName = dbStruct.BadgeName.String
	svcStruct.AppreciationPoints = dbStruct.AppreciationPoints
	return
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

// func (us *service) sendRewardQuotaRefillEmailToAll(ctx context.Context) {

// 	reqData := dto.ListUsersReq{
// 		Page:     1,
// 		PageSize: 1000,
// 	}
// 	dbUsers, _, err := us.userRepo.ListUsers(ctx, reqData)
// 	if err != nil {
// 		logger.Errorf("error in getting users for email")
// 		return
// 	}

// 	usersEmails := make([]string, 0)
// 	for _, user := range dbUsers {
// 		usersEmails = append(usersEmails, user.Email)
// 	}

// 	return
// }


func mapDbUserToUserListResp(dbStruct repository.User) (svcData dto.UserDetails) {
	svcData.Id = dbStruct.Id
	svcData.FirstName = dbStruct.FirstName
	svcData.LastName = dbStruct.LastName
	svcData.Email = dbStruct.Email
	return svcData
}

func (us *service) NotificationByAdmin(ctx context.Context, notificationReq dto.AdminNotificationReq) (err error) {

	notificationTokens, err := us.userRepo.ListDeviceTokensByUserID(ctx, notificationReq.Id)
	if err != nil {
		logger.Errorf(ctx, "err in getting device tokens: %v", err)
		err = apperrors.InternalServerError
		return
	}

	if notificationReq.All {
		err = notificationReq.Message.SendNotificationToTopic("peerly")
		if err != nil {
			return
		}
		return
	}

	for _, notificationToken := range notificationTokens {
		err = notificationReq.Message.SendNotificationToNotificationToken(notificationToken)
		if err != nil {
			return
		}
	}

	return
}

func (us *service) AllAppreciationReport(ctx context.Context, appreciations []dto.AppreciationResponse) (tempFileName string, err error) {

	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet
	sheetName := "Appreciations"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		logger.Errorf(ctx, "err in generating newsheet, err: %v", err)
		return
	}

	// Set header
	headers := []string{"Core value", "Core value description", "Appreciation description", "Sender first name", "Sender last name", "Sender designation", "Receiver first name", "Receiver last name", "Receiver designation", "Total rewards", "Total reward points"}
	for colIndex, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+colIndex))
		f.SetCellValue(sheetName, cell, header)
	}

	// Add data to the sheet
	for rowIndex, app := range appreciations {
		row := rowIndex + 2 // Starting from row 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), app.CoreValueName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), app.CoreValueDesc)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), app.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), app.SenderFirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), app.SenderLastName)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), app.SenderDesignation)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), app.ReceiverFirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), app.ReceiverLastName)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), app.ReceiverDesignation)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), app.TotalRewards)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), app.TotalRewardPoints)
	}

	// Set the active sheet
	f.SetActiveSheet(index)

	// Save the Excel file temporarily
	tempFileName = "report.xlsx"
	if err = f.SaveAs(tempFileName); err != nil {
		logger.Errorf(ctx, "Failed to save file: %v", err)
		return
	}

	return
}

func (us *service) ReportedAppreciationReport(ctx context.Context, appreciations []dto.ReportedAppreciation) (tempFileName string, err error) {

	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet
	sheetName := "ReportedAppreciations"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		logger.Errorf(ctx, "err in generating newsheet, err: %v", err)
		return
	}

	// Set header
	headers := []string{"Core value", "Core value description", "Appreciation description", "Sender first name", "Sender last name", "Sender designation", "Receiver first name", "Receiver last name", "Receiver designation", "Reporting Comment", "Reported by first name", "Reported by last name", "Reported at", "Moderator comment", "Moderator first name", "Moderator last name", "Status"}
	for colIndex, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+colIndex))
		f.SetCellValue(sheetName, cell, header)
	}

	// Add data to the sheet
	for rowIndex, app := range appreciations {
		row := rowIndex + 2 // Starting from row 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), app.CoreValueName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), app.CoreValueDesc)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), app.AppreciationDesc)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), app.SenderFirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), app.SenderLastName)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), app.SenderDesignation)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), app.ReceiverFirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), app.ReceiverLastName)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), app.ReceiverDesignation)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), app.ReportingComment)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), app.ReportedByFirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), app.ReportedByLastName)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), app.ReportedAt)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), app.ModeratorComment)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), app.ModeratedByFirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), app.ModeratedByLastName)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), app.Status)
	}

	// Set the active sheet
	f.SetActiveSheet(index)

	// Save the Excel file temporarily
	tempFileName = "reportedAppreciations.xlsx"
	if err = f.SaveAs(tempFileName); err != nil {
		logger.Errorf(ctx, "Failed to save file: %v", err)
		return
	}

	return
}
