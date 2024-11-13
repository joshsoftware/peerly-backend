package reportappreciations

import (
	"context"
	"fmt"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/app/email"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/joshsoftware/peerly-backend/internal/repository"
)

type service struct {
	reportAppreciationRepo repository.ReportAppreciationStorer
	userRepo               repository.UserStorer
	appreciationRepo       repository.AppreciationStorer
}

type Service interface {
	ReportAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error)
	ListReportedAppreciations(ctx context.Context) (dto.ListReportedAppreciationsResponse, error)
	GetReportedAppreciationByAppreciationID(ctx context.Context, appreciationID int64) (dto.ReportedAppreciation, error)
	DeleteAppreciation(ctx context.Context, reqData dto.ModerationReq) (err error)
	ResolveAppreciation(ctx context.Context, reqData dto.ModerationReq) (err error)
}

func NewService(reportAppreciationRepo repository.ReportAppreciationStorer, userRepo repository.UserStorer, appreciationRepo repository.AppreciationStorer) Service {
	return &service{
		reportAppreciationRepo: reportAppreciationRepo,
		userRepo:               userRepo,
		appreciationRepo:       appreciationRepo,
	}
}

func (rs *service) ReportAppreciation(ctx context.Context, reqData dto.ReportAppreciationReq) (resp dto.ReportAppricaitionResp, err error) {

	reporterId := ctx.Value(constants.UserId)
	fmt.Printf("reporterId: %T", reporterId)
	data, ok := reporterId.(int64)
	if !ok {
		logger.Error(ctx, "Error in typecasting reporter id")
		err = apperrors.InternalServerError
		return
	}
	reqData.ReportedBy = data

	doesAppreciationExist, err := rs.reportAppreciationRepo.CheckAppreciation(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	if !doesAppreciationExist {
		err = apperrors.InvalidId
		return
	}

	isDupliate, err := rs.reportAppreciationRepo.CheckDuplicateReport(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	if isDupliate {
		err = apperrors.RepeatedReport
		return
	}

	usersData, err := rs.reportAppreciationRepo.GetSenderAndReceiver(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}
	if usersData.Sender == reqData.ReportedBy || usersData.Receiver == reqData.ReportedBy {
		err = apperrors.CannotReportOwnAppreciation
		return
	}

	resp, err = rs.reportAppreciationRepo.ReportAppreciation(ctx, reqData)
	if err != nil {
		err = apperrors.InternalServerError
		return
	}

	quaterTimeStamp := GetQuarterStartUnixTime()

	reqGetUserById := dto.GetUserByIdReq{
		UserId:          data,
		QuaterTimeStamp: quaterTimeStamp,
	}
	senderInfo, err := rs.userRepo.GetUserById(ctx, reqGetUserById)
	if err != nil {
		return
	}

	apprInfo, err := rs.appreciationRepo.GetAppreciationById(ctx, nil, int32(reqData.AppreciationId))
	if err != nil {
		return
	}
	err = sendReportEmail(senderInfo.Email,
		senderInfo.FirstName,
		senderInfo.LastName,
		apprInfo.SenderFirstName,
		apprInfo.SenderLastName,
		apprInfo.ReceiverFirstName,
		apprInfo.ReceiverLastName,
		reqData.ReportingComment)

	return
}

func (rs *service) ListReportedAppreciations(ctx context.Context) (dto.ListReportedAppreciationsResponse, error) {

	var resp dto.ListReportedAppreciationsResponse

	var appreciationList []dto.ReportedAppreciation

	appreciations, err := rs.reportAppreciationRepo.ListReportedAppreciations(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}

	for _, appreciation := range appreciations {

		senderDataReq := dto.GetUserByIdReq{
			UserId:          appreciation.Sender,
			QuaterTimeStamp: GetQuarterStartUnixTime(),
		}

		sender, err := rs.userRepo.GetUserById(ctx, senderDataReq)
		if err != nil {
			return resp, err
		}

		receiverDataReq := dto.GetUserByIdReq{
			UserId:          appreciation.Receiver,
			QuaterTimeStamp: GetQuarterStartUnixTime(),
		}

		receiver, err := rs.userRepo.GetUserById(ctx, receiverDataReq)
		if err != nil {
			return resp, err
		}

		reporterDataReq := dto.GetUserByIdReq{
			UserId:          appreciation.ReportedBy,
			QuaterTimeStamp: GetQuarterStartUnixTime(),
		}

		reporter, err := rs.userRepo.GetUserById(ctx, reporterDataReq)
		if err != nil {
			return resp, err
		}

		moderatorDataReq := dto.GetUserByIdReq{
			UserId:          appreciation.ModeratedBy.Int64,
			QuaterTimeStamp: GetQuarterStartUnixTime(),
		}

		var moderator dto.GetUserByIdResp
		if appreciation.ModeratedBy.Valid {
			moderator, err = rs.userRepo.GetUserById(ctx, moderatorDataReq)
			if err != nil {
				return resp, err
			}
		}

		svcApp := mapDbAppreciationsToSvcAppreciations(appreciation, sender, receiver, reporter, moderator)

		appreciationList = append(appreciationList, svcApp)

	}
	resp.Appreciations = appreciationList
	return resp, err
}

func (rs *service) GetReportedAppreciationByAppreciationID(ctx context.Context, appreciationID int64) (dto.ReportedAppreciation, error) {

	var resp dto.ReportedAppreciation

	// var appreciationList []dto.ReportedAppreciation

	appreciation, err := rs.reportAppreciationRepo.GetReportedAppreciationByAppreciationID(ctx, appreciationID)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return resp, err
	}

	senderDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.Sender,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	sender, err := rs.userRepo.GetUserById(ctx, senderDataReq)
	if err != nil {
		return resp, err
	}

	receiverDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.Receiver,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	receiver, err := rs.userRepo.GetUserById(ctx, receiverDataReq)
	if err != nil {
		return resp, err
	}

	reporterDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.ReportedBy,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	reporter, err := rs.userRepo.GetUserById(ctx, reporterDataReq)
	if err != nil {
		return resp, err
	}

	moderatorDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.ModeratedBy.Int64,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	var moderator dto.GetUserByIdResp
	if appreciation.ModeratedBy.Valid {
		moderator, err = rs.userRepo.GetUserById(ctx, moderatorDataReq)
		if err != nil {
			return resp, err
		}
	}

	resp = mapDbAppreciationsToSvcAppreciations(appreciation, sender, receiver, reporter, moderator)

	return resp, err
}
func (rs *service) DeleteAppreciation(ctx context.Context, reqData dto.ModerationReq) (err error) {
	moderatorId := ctx.Value(constants.UserId)
	fmt.Printf("moderatorId: %T", moderatorId)
	data, ok := moderatorId.(int64)
	if !ok {
		logger.Error(ctx, "Error in typecasting moderator id")
		err = apperrors.InternalServerError
		return
	}
	reqData.ModeratedBy = data

	appreciation, err := rs.reportAppreciationRepo.GetResolution(ctx, reqData.ResolutionId)
	if err != nil {
		return
	}

	reqData.AppreciationId = appreciation.Appreciation_id
	err = rs.reportAppreciationRepo.DeleteAppreciation(ctx, reqData)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	senderDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.Sender,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	sender, err := rs.userRepo.GetUserById(ctx, senderDataReq)
	if err != nil {
		return
	}

	receiverDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.Receiver,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	receiver, err := rs.userRepo.GetUserById(ctx, receiverDataReq)
	if err != nil {
		return
	}

	reporterDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.ReportedBy,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	reporter, err := rs.userRepo.GetUserById(ctx, reporterDataReq)
	if err != nil {
		return
	}

	seconds := appreciation.CreatedAt / 1000
	nanoseconds := (appreciation.CreatedAt % 1000) * 1e6

	tm := time.Unix(seconds, nanoseconds)
	formattedDate := tm.Format("02/01/2006") // date format: dd/mm/yyyy

	templateData := dto.DeleteAppreciationMail{
		ModeratorComment: reqData.ModeratorComment,
		AppreciationBy:   sender.FirstName + " " + sender.LastName,
		AppreciationTo:   receiver.FirstName + " " + receiver.LastName,
		ReportingComment: appreciation.ReportingComment,
		AppreciationDesc: appreciation.AppreciationDesc,
		Date:             formattedDate,
		Icon:             config.PeerlyBaseUrl() + constants.CheckIconLogo,
	}

	fmt.Println("Reporter mail: ", reporter.Email)
	err = sendDeleteEmail(reporter.Email, sender.Email, receiver.Email, templateData)

	return
}

func GetQuarterStartUnixTime() int64 {
	// Example function to get the Unix timestamp of the start of the quarter
	now := time.Now()
	quarterStart := time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.UTC)
	return quarterStart.Unix() * 1000 // convert to milliseconds
}

func mapDbAppreciationsToSvcAppreciations(dbApp repository.ListReportedAppreciations, sender dto.GetUserByIdResp, receiver dto.GetUserByIdResp, reporter dto.GetUserByIdResp, moderator dto.GetUserByIdResp) (svcApp dto.ReportedAppreciation) {
	svcApp.Id = dbApp.Id
	svcApp.Appreciation_id = dbApp.Appreciation_id
	svcApp.AppreciationDesc = dbApp.AppreciationDesc
	svcApp.TotalRewardPoints = dbApp.TotalRewardPoints
	svcApp.Quarter = dbApp.TotalRewardPoints
	svcApp.CoreValueName = dbApp.CoreValueName
	svcApp.CoreValueDesc = dbApp.CoreValueDesc
	svcApp.SenderFirstName = sender.FirstName
	svcApp.SenderLastName = sender.LastName
	svcApp.SenderImgUrl = sender.ProfileImgUrl
	svcApp.SenderDesignation = sender.Designation
	svcApp.ReceiverFirstName = receiver.FirstName
	svcApp.ReceiverLastName = receiver.LastName
	svcApp.ReceiverImgUrl = receiver.ProfileImgUrl
	svcApp.ReceiverDesignation = receiver.Designation
	svcApp.CreatedAt = dbApp.CreatedAt
	svcApp.ReportingComment = dbApp.ReportingComment
	svcApp.ReportedByFirstName = reporter.FirstName
	svcApp.ReportedByLastName = reporter.LastName
	svcApp.ReportedAt = dbApp.ReportedAt
	svcApp.IsValid = dbApp.IsValid
	if (moderator != dto.GetUserByIdResp{}) {
		svcApp.ModeratedAt = dbApp.ModeratedAt.Int64
		svcApp.ModeratedByFirstName = moderator.FirstName
		svcApp.ModeratedByLastName = moderator.LastName
		svcApp.ModeratorComment = dbApp.ModeratorComment.String
	}
	svcApp.Status = dbApp.Status
	return
}

func sendReportEmail(senderEmail string, senderFirstName string, senderLastName string, apprSenderFirstName string, apprSenderLastName string, apprReceiverFirstName string, apprReceiverLastName string, reportingComment string) error {

	templateData := struct {
		SenderName               string
		ReportingComment         string
		AppreciationSenderName   string
		AppreciationReceiverName string
	}{
		SenderName:               fmt.Sprint(senderFirstName, " ", senderLastName),
		ReportingComment:         reportingComment,
		AppreciationSenderName:   fmt.Sprint(apprSenderFirstName, " ", apprSenderLastName),
		AppreciationReceiverName: fmt.Sprint(apprReceiverFirstName, " ", apprReceiverLastName),
	}

	ctx := context.Background()
	logger.Info(ctx, "report sender email: ---------> ", senderEmail)
	mailReq := email.NewMail([]string{senderEmail}, []string{"dl_peerly.support@joshsoftware.com"}, []string{}, "🙏 Thanks for Your Feedback! We’re On It! 🔧")
	mailReq.ParseTemplate("./internal/app/email/templates/reportAppreciation.html", templateData)
	err := mailReq.Send()
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return err
	}
	return nil
}

func (rs *service) ResolveAppreciation(ctx context.Context, reqData dto.ModerationReq) (err error) {
	moderatorId := ctx.Value(constants.UserId)
	fmt.Printf("moderatorId: %T", moderatorId)
	data, ok := moderatorId.(int64)
	if !ok {
		logger.Error(ctx, "Error in typecasting moderator id")
		err = apperrors.InternalServerError
		return
	}
	reqData.ModeratedBy = data
	appreciation, err := rs.reportAppreciationRepo.GetResolution(ctx, reqData.ResolutionId)
	if err != nil {
		return
	}

	reqData.AppreciationId = appreciation.Appreciation_id
	err = rs.reportAppreciationRepo.ResolveAppreciation(ctx, reqData)
	if err != nil {
		logger.Error(ctx, err.Error())
		err = apperrors.InternalServerError
		return
	}

	senderDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.Sender,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	sender, err := rs.userRepo.GetUserById(ctx, senderDataReq)
	if err != nil {
		return
	}

	receiverDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.Receiver,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	receiver, err := rs.userRepo.GetUserById(ctx, receiverDataReq)
	if err != nil {
		return
	}

	reporterDataReq := dto.GetUserByIdReq{
		UserId:          appreciation.ReportedBy,
		QuaterTimeStamp: GetQuarterStartUnixTime(),
	}

	reporter, err := rs.userRepo.GetUserById(ctx, reporterDataReq)
	if err != nil {
		return
	}

	templateData := dto.ResolveAppreciationMail{
		ModeratorComment: reqData.ModeratorComment,
		AppreciationBy:   sender.FirstName + " " + sender.LastName,
		AppreciationTo:   receiver.FirstName + " " + receiver.LastName,
		ReportingComment: appreciation.ReportingComment,
		AppreciationDesc: appreciation.AppreciationDesc,
		Icon:             config.PeerlyBaseUrl() + constants.CheckIconLogo,
	}

	fmt.Println("Reporter mail: ", reporter.Email)
	err = sendResolveEmail(reporter.Email, templateData)
	return
}

func sendDeleteEmail(reporterEmail string, senderEmail string, receiverEmail string, templateData dto.DeleteAppreciationMail) error {

	ctx := context.Background()
	logger.Info(ctx, "reporter email: ---------> ", reporterEmail)
	mailReq := email.NewMail([]string{reporterEmail}, []string{}, []string{}, "Results of reported appreciation")
	mailReq.ParseTemplate("./internal/app/email/templates/deleteAppreciation.html", templateData)
	err := mailReq.Send()
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return err
	}

	logger.Info(ctx, "sender email: ---------> ", senderEmail)
	mailReq = email.NewMail([]string{senderEmail}, []string{}, []string{}, "Results of reported appreciation")
	mailReq.ParseTemplate("./internal/app/email/templates/senderDeleteEmail.html", templateData)
	err = mailReq.Send()
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return err
	}

	logger.Info(ctx, "receiver email: ---------> ", receiverEmail)
	mailReq = email.NewMail([]string{receiverEmail}, []string{}, []string{}, "Results of reported appreciation")
	mailReq.ParseTemplate("./internal/app/email/templates/receiverDeleteEmail.html", templateData)
	err = mailReq.Send()
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return err
	}

	return nil
}

func sendResolveEmail(senderEmail string, templateData dto.ResolveAppreciationMail) error {

	ctx := context.Background()
	logger.Info(ctx, "report sender email: ---------> ", senderEmail)
	mailReq := email.NewMail([]string{senderEmail}, []string{}, []string{}, "Results of reported appreciation")
	mailReq.ParseTemplate("./internal/app/email/templates/resolveAppreciation.html", templateData)
	err := mailReq.Send()
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return err
	}
	return nil
}
