package appreciation

import (
	"context"
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/app/email"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"

	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
)

type service struct {
	appreciationRepo repository.AppreciationStorer
	corevaluesRespo  repository.CoreValueStorer
	userRepo         repository.UserStorer
}

// Service contains all
type Service interface {
	CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error)
	GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error)
	ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error)
	DeleteAppreciation(ctx context.Context, apprId int32) error
	UpdateAppreciation(ctx context.Context) (bool, error)
	sendAppreciationNotificationToReceiver(ctx context.Context, appr repository.AppreciationResponse)
	sendAppreciationNotificationToAll(ctx context.Context, appr repository.AppreciationResponse)
	sendEmailForBadgeAllocation(userBadgeDetails []repository.UserBadgeDetails)
}

func NewService(appreciationRepo repository.AppreciationStorer, coreValuesRepo repository.CoreValueStorer, userRepo repository.UserStorer) Service {
	return &service{
		appreciationRepo: appreciationRepo,
		corevaluesRespo:  coreValuesRepo,
		userRepo:         userRepo,
	}
}

func (apprSvc *service) CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error) {

	//add quarter
	appreciation.Quarter = utils.GetQuarter()
	logger.Debug(ctx,"service: CreateAppreciation: appreciation: ",appreciation)

	//add sender
	data := ctx.Value(constants.UserId)
	sender, ok := data.(int64)
	if !ok {
		logger.Error(ctx,"service: err in parsing userid from token")
		return dto.Appreciation{}, apperrors.InternalServer
	}

	//check is receiver present in database
	chk, err := apprSvc.appreciationRepo.IsUserPresent(ctx, nil, appreciation.Receiver)
	if err != nil {
		logger.Errorf(ctx,"err: %v", err)
		return dto.Appreciation{}, err
	}
	if !chk {
		logger.Errorf(ctx,"User not found (user_id): %v",appreciation.Receiver)
		return dto.Appreciation{}, apperrors.UserNotFound
	}
	appreciation.Sender = sender

	logger.Debug(ctx,"service: appreciation: ",appreciation)
	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		logger.Errorf(ctx,"err: %v",err)
		return dto.Appreciation{}, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Infof(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := apprSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Infof(ctx, "error in handle transaction, err: %s", txErr.Error())
			return
		}
	}()

	//check is corevalue present in database
	_, err = apprSvc.corevaluesRespo.GetCoreValue(ctx, int64(appreciation.CoreValueID))
	if err != nil {
		logger.Errorf(ctx,"service: err: %v", err)
		return dto.Appreciation{}, err
	}

	// check self appreciation
	if appreciation.Receiver == sender {
		logger.Errorf(ctx,"Self Appreciation Error: %v \n Userid: %d",apperrors.SelfAppreciationError,appreciation.Receiver)
		return dto.Appreciation{}, apperrors.SelfAppreciationError
	}

	appr, err := apprSvc.appreciationRepo.CreateAppreciation(ctx, tx, appreciation)
	if err != nil {
		logger.Errorf(ctx,"err: %v", err)
		return dto.Appreciation{}, err
	}

	res := mapAppreciationDBToDTO(appr)
	apprInfo, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, tx, int32(res.ID))
	if err != nil {
		logger.Errorf(ctx,"err: %v", err)
		return res, nil
	}

	logger.Debug(ctx,"service: createAppreciation result: ",res)
	quaterTimeStamp := user.GetQuarterStartUnixTime()

	reqGetUserById := dto.GetUserByIdReq{
		UserId:          sender,
		QuaterTimeStamp: quaterTimeStamp,
	}
	senderInfo,err := apprSvc.userRepo.GetUserById(ctx,reqGetUserById)
	if err != nil{
		logger.Info(ctx,"error in getting create appreciation sender info")
		return res, nil
	}

	reqGetUserById.UserId = appreciation.Receiver
	receiverInfo,err := apprSvc.userRepo.GetUserById(ctx,reqGetUserById)
	if err != nil{
		logger.Info(ctx,"error in getting create appreciation sender info")
		return res, nil
	}
	sendAppreciationEmail(apprInfo,senderInfo.Email,receiverInfo.Email)
	apprSvc.sendAppreciationNotificationToReceiver(ctx, apprInfo)
	apprSvc.sendAppreciationNotificationToAll(ctx, apprInfo)
	return res, nil
}

func (apprSvc *service) GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error) {

	logger.Debug(ctx,"service: appreciationId: ",appreciationId)
	resAppr, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, nil, appreciationId)
	if err != nil {
		logger.Errorf(ctx,"err: %v", err)
		return dto.AppreciationResponse{}, err
	}
	logger.Debug(ctx,"service: resAppr: ",resAppr)
	return mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(resAppr), nil
}

func (apprSvc *service) ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error) {

	logger.Debug(ctx,"service: filter: ",filter)
	infos, pagination, err := apprSvc.appreciationRepo.ListAppreciations(ctx, nil, filter)
	if err != nil {
		logger.Errorf(ctx,"err: %v", err)
		return dto.ListAppreciationsResponse{}, err
	}
	logger.Debug(ctx,"service: infos: ",infos," pagination: ",pagination)

	responses := make([]dto.AppreciationResponse, 0)
	for _, info := range infos {
		responses = append(responses, mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info))
	}
	paginationResp := dtoPagination(pagination)
	logger.Debug(ctx," apprSvc: ",responses," paginationResp: ",paginationResp)
	return dto.ListAppreciationsResponse{Appreciations: responses, MetaData: paginationResp}, nil
}

func (apprSvc *service) DeleteAppreciation(ctx context.Context, apprId int32) error {
	logger.Debug(ctx,"service: apprId: ",apprId)
	return apprSvc.appreciationRepo.DeleteAppreciation(ctx, nil, apprId)
}

func (apprSvc *service) UpdateAppreciation(ctx context.Context) (bool, error) {

	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		logger.Errorf(ctx,"error in begin transaction: %v",err)
		return false, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Infof(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := apprSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Infof(ctx, "error in handle transaction, err: %s", txErr.Error())
			return
		}
	}()

	_, err = apprSvc.appreciationRepo.UpdateAppreciationTotalRewardsOfYesterday(ctx, tx)

	if err != nil {
		logger.Errorf(ctx,"err: %v", err)
		return false, err
	}
	logger.Debug(ctx,"service: UpdateAppreciationTotalRewardsOfYesterday completed")

	userBadgeDetails, err := apprSvc.appreciationRepo.UpdateUserBadgesBasedOnTotalRewards(ctx, tx)
	if err != nil {
		logger.Error(ctx,"err: ", err.Error())
		return false, err
	}
	logger.Debug(ctx,"service: UpdateAppreciation: ",userBadgeDetails)
	apprSvc.sendEmailForBadgeAllocation(userBadgeDetails)
	return true, nil
}

func sendAppreciationEmail(emailData repository.AppreciationResponse,senderEmail string,receiverEmail string) error {
	// Plain text content
	plainTextContent := "Samnit " + "123456"

	templateData := struct {
		SenderName    string
		ReceiverName string
		Description   string
		CoreValueName string
	}{
		SenderName:    fmt.Sprint(emailData.SenderFirstName, " ", emailData.SenderLastName),
		ReceiverName: fmt.Sprint(emailData.ReceiverFirstName, " ", emailData.ReceiverLastName),
		Description:   emailData.Description,
		CoreValueName: emailData.CoreValueName,
	}

	logger.Infof(context.Background(),"appreciation sender email: %v :receiver email: %v  ",senderEmail,receiverEmail)
	mailReq := email.NewMail([]string{receiverEmail}, []string{senderEmail}, []string{}, fmt.Sprintf("%s %s appreciated %s %s",emailData.SenderFirstName,emailData.SenderLastName,emailData.ReceiverFirstName,emailData.ReceiverLastName))
	err := mailReq.ParseTemplate("./internal/app/email/templates/createAppreciation.html", templateData)
	if err != nil {
		logger.Errorf(context.Background(),"err in creating html file : %v", err)
		return err
	}
	err = mailReq.Send(plainTextContent)
	if err != nil {
		logger.Errorf(context.Background(),"err: %v", err)
		return err
	}
	return nil
}

func (apprSvc *service) sendAppreciationNotificationToReceiver(ctx context.Context, appr repository.AppreciationResponse) {

	logger.Debug(ctx,"service: apprResponse: ",appr)
	notificationTokens, err := apprSvc.userRepo.ListDeviceTokensByUserID(ctx, appr.ReceiverID)
	if err != nil {
		logger.Errorf(ctx,"err in getting device tokens: %v", err)
		return
	}

	logger.Debug(ctx,"apprSvc: notificationTokens: ",notificationTokens)
	msg := notification.Message{
		Title: "Received Appreciation",
		Body:  fmt.Sprintf("You've been appreciated by %s %s! Well done and keep up the JOSH!", appr.SenderFirstName, appr.SenderLastName),
	}

	logger.Infof(ctx,"message: %v",msg)
	for _, notificationToken := range notificationTokens {
		msg.SendNotificationToNotificationToken(notificationToken)
	}

}

func (apprSvc *service) sendAppreciationNotificationToAll(ctx context.Context, appr repository.AppreciationResponse) {

	logger.Debug(ctx," apprSvc: appr: ",appr)
	msg := notification.Message{
		Title: "Appreciation",
		Body:  fmt.Sprintf(" %s %s appreciated %s %s", appr.SenderFirstName, appr.SenderLastName, appr.ReceiverFirstName, appr.ReceiverLastName),
	}
	logger.Infof(ctx,"message: %v",msg)
	msg.SendNotificationToTopic("peerly")
}

func (apprSvc *service) sendEmailForBadgeAllocation(userBadgeDetails []repository.UserBadgeDetails) {

	logger.Debug(context.Background(),"service: user Badge Details: ", userBadgeDetails)
	for _, userBadgeDetail := range userBadgeDetails {

		// Determine the BadgeImageUrl based on the BadgeName
		var badgeImageUrl string
		switch userBadgeDetail.BadgeName.String {
		case "Bronze":
			badgeImageUrl = "bronzeBadge"
		case "Silver":
			badgeImageUrl = "silverBadge"
		case "Gold":
			badgeImageUrl = "goldBadge"
		case "Platinum":
			badgeImageUrl = "platinumBadge"
		}

		// repository.UserBadgeDetails
		templateData := struct {
			EmployeeName       string
			BadgeName          string
			BadgeImageName     string
			AppreciationPoints int32
		}{
			EmployeeName:       fmt.Sprint(userBadgeDetail.FirstName, " ", userBadgeDetail.LastName),
			BadgeName:          userBadgeDetail.BadgeName.String,
			BadgeImageName:     badgeImageUrl,
			AppreciationPoints: userBadgeDetail.BadgePoints,
		}
		logger.Info(context.Background(),"service: badge data: ", templateData)
		mailReq := email.NewMail([]string{userBadgeDetail.Email}, []string{}, []string{}, "Received an badge")
		err := mailReq.ParseTemplate("./internal/app/email/templates/badge.html", templateData)
		if err != nil {
			logger.Errorf(context.Background(),"service: err in creating html file : %v", err)
			return
		}
		err = mailReq.Send("badge allocation")
		if err != nil {
			logger.Errorf(context.Background(),"service: err: %v", err)
			return
		}
		logger.Infof(context.Background(),"service: mail request: %v",mailReq)
	}
}
