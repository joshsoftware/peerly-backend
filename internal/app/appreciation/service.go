package appreciation

import (
	"context"
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/app/email"
	"github.com/joshsoftware/peerly-backend/internal/app/notification"
	user "github.com/joshsoftware/peerly-backend/internal/app/users"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
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
	UpdateAppreciation(ctx context.Context, orgTimezone string) (bool, error)
}

func NewService(appreciationRepo repository.AppreciationStorer, coreValuesRepo repository.CoreValueStorer, userRepo repository.UserStorer) Service {
	return &service{
		appreciationRepo: appreciationRepo,
		corevaluesRespo:  coreValuesRepo,
		userRepo:         userRepo,
	}
}

func (apprSvc *service) CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error) {

	logger.Debug(ctx, "svc: CreateAppreciation: appreciation: ", appreciation)
	//add quarter
	appreciation.Quarter = utils.GetQuarter()
	logger.Debug(ctx, "appreciationService CreateAppreciation: appreciation: ", appreciation)

	//add sender
	data := ctx.Value(constants.UserId)
	sender, ok := data.(int64)
	if !ok {
		logger.Error(ctx, "appreciationService err in parsing userid from token")
		return dto.Appreciation{}, apperrors.InternalServer
	}

	//check is receiver present in database
	chk, err := apprSvc.appreciationRepo.IsUserPresent(ctx, nil, appreciation.Receiver)
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return dto.Appreciation{}, err
	}
	if !chk {
		logger.Errorf(ctx, "appreciationService User not found (user_id): %v", appreciation.Receiver)
		return dto.Appreciation{}, apperrors.UserNotFound
	}
	appreciation.Sender = sender

	logger.Debug(ctx, "appreciationService appreciation: ", appreciation)
	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		logger.Errorf(ctx, "servive: err: %v", err)
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
		logger.Errorf(ctx, "appreciationService err: %v", err)
		return dto.Appreciation{}, err
	}

	// check self appreciation
	if appreciation.Receiver == sender {
		logger.Errorf(ctx, "appreciationService Self Appreciation Error: %v \n Userid: %d", apperrors.SelfAppreciationError, appreciation.Receiver)
		return dto.Appreciation{}, apperrors.SelfAppreciationError
	}

	appr, err := apprSvc.appreciationRepo.CreateAppreciation(ctx, tx, appreciation)
	if err != nil {
		logger.Errorf(ctx, "appreciationService err: %v", err)
		return dto.Appreciation{}, err
	}

	res := mapAppreciationDBToDTO(appr)
	apprInfo, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, tx, int32(res.ID))
	if err != nil {
		logger.Errorf(ctx, "appreciationService err: %v", err)
		return res, nil
	}

	logger.Debug(ctx, "appreciationService createAppreciation result: ", res)
	quaterTimeStamp := user.GetQuarterStartUnixTime()

	reqGetUserById := dto.GetUserByIdReq{
		UserId:          sender,
		QuaterTimeStamp: quaterTimeStamp,
	}
	senderInfo, err := apprSvc.userRepo.GetUserById(ctx, reqGetUserById)
	if err != nil {
		logger.Info(ctx, "appreciationService error in getting create appreciation sender info")
	}

	reqGetUserById.UserId = appreciation.Receiver
	receiverInfo, err := apprSvc.userRepo.GetUserById(ctx, reqGetUserById)
	if err != nil {
		logger.Info(ctx, "appreciationService error in getting create appreciation sender info")
	}
	err = sendAppreciationEmail(apprInfo, senderInfo.Email, receiverInfo.Email)
	apprSvc.sendAppreciationNotificationToReceiver(ctx, apprInfo)
	apprSvc.sendAppreciationNotificationToAll(ctx, apprInfo)
	return res, nil
}

func (apprSvc *service) GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error) {

	logger.Debug(ctx, "appreciationService appreciationId: ", appreciationId)
	resAppr, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, nil, appreciationId)
	if err != nil {
		logger.Errorf(ctx, "appreciationService err: %v", err)
		return dto.AppreciationResponse{}, err
	}
	logger.Debug(ctx, "appreciationService resAppr: ", resAppr)
	return mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(resAppr), nil
}

func (apprSvc *service) ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error) {

	logger.Debug(ctx, "appreciationService filter: ", filter)
	infos, pagination, err := apprSvc.appreciationRepo.ListAppreciations(ctx, nil, filter)
	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return dto.ListAppreciationsResponse{}, err
	}
	logger.Debug(ctx, "appreciationService infos: ", infos, " pagination: ", pagination)

	responses := make([]dto.AppreciationResponse, 0)
	for _, info := range infos {
		responses = append(responses, mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info))
	}
	paginationResp := dtoPagination(pagination)
	logger.Debug(ctx, " appreciationService ", responses, " paginationResp: ", paginationResp)
	return dto.ListAppreciationsResponse{Appreciations: responses, MetaData: paginationResp}, nil
}

func (apprSvc *service) DeleteAppreciation(ctx context.Context, apprId int32) error {
	logger.Debug(ctx, "appreciationService apprId: ", apprId)
	return apprSvc.appreciationRepo.DeleteAppreciation(ctx, nil, apprId)
}

func (apprSvc *service) UpdateAppreciation(ctx context.Context, orgTimezone string) (bool, error) {

	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		logger.Errorf(ctx, "appreciationService error in begin transaction: %v", err)
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
			logger.Infof(ctx, "appreciationService error in handle transaction, err: %s", txErr.Error())
			return
		}
	}()

	_, err = apprSvc.appreciationRepo.UpdateAppreciationTotalRewardsOfYesterday(ctx, tx, orgTimezone)

	if err != nil {
		logger.Errorf(ctx, "err: %v", err)
		return false, err
	}
	logger.Debug(ctx, "appreciationService UpdateAppreciationTotalRewardsOfYesterday completed")

	userBadgeDetails, err := apprSvc.appreciationRepo.UpdateUserBadgesBasedOnTotalRewards(ctx, tx)
	if err != nil {
		logger.Error(ctx, "appreciationService err: ", err.Error())
		return false, err
	}
	logger.Debug(ctx, "appreciationService UpdateAppreciation: ", userBadgeDetails)
	apprSvc.sendEmailForBadgeAllocation(userBadgeDetails)
	return true, nil
}

func sendAppreciationEmail(emailData repository.AppreciationResponse, senderEmail string, receiverEmail string) error {

	templateData := struct {
		SenderName               string
		ReceiverName             string
		Description              string
		CoreValueName            string
		CoreValueBackgroundColor string
	}{
		SenderName:               fmt.Sprint(emailData.SenderFirstName, " ", emailData.SenderLastName),
		ReceiverName:             fmt.Sprint(emailData.ReceiverFirstName, " ", emailData.ReceiverLastName),
		Description:              emailData.Description,
		CoreValueName:            emailData.CoreValueName,
		CoreValueBackgroundColor: utils.GetCoreValueBackgroundColor(emailData.CoreValueName),
	}

	logger.Infof(context.Background(), "appreciation sender email: %v :receiver email: %v  ", senderEmail, receiverEmail)
	mailReq := email.NewMail([]string{receiverEmail}, []string{}, []string{}, fmt.Sprintf("Kudos! You've Been Praised by %s %s! 🎉 ", emailData.SenderFirstName, emailData.SenderLastName))
	err := mailReq.ParseTemplate("./internal/app/email/templates/receiverAppreciation.html", templateData)
	if err != nil {
		logger.Errorf(context.Background(), "err in creating html file : %v", err)
		return err
	}
	err = mailReq.Send()
	if err != nil {
		logger.Errorf(context.Background(), "appreciationService err: %v", err)
		return err
	}
	mailReq = email.NewMail([]string{senderEmail}, []string{}, []string{}, fmt.Sprintf("Your appreciation to %s %s has been sent! 🙌", emailData.ReceiverFirstName, emailData.ReceiverLastName))
	err = mailReq.ParseTemplate("./internal/app/email/templates/senderAppreciation.html", templateData)
	if err != nil {
		logger.Errorf(context.Background(), "appreciationService err: %v", err)
		return err
	}
	err = mailReq.Send()
	if err != nil {
		logger.Errorf(context.Background(), "appreciationService err: %v", err)
		return err
	}
	return nil
}

func (apprSvc *service) sendAppreciationNotificationToReceiver(ctx context.Context, appr repository.AppreciationResponse) {

	logger.Debug(ctx, "appreciationService apprResponse: ", appr)
	notificationTokens, err := apprSvc.userRepo.ListDeviceTokensByUserID(ctx, appr.ReceiverID)
	if err != nil {
		logger.Errorf(ctx, "appreciationService err in getting device tokens: %v", err)
		return
	}

	logger.Debug(ctx, "appreciationService notificationTokens: ", notificationTokens)
	msg := notification.Message{
		Title: "Appreciation incoming!",
		Body:  fmt.Sprintf("You've been appreciated by %s %s! Well done and keep up the JOSH!", appr.SenderFirstName, appr.SenderLastName),
	}

	logger.Infof(ctx, "appreciationService message: %v", msg)
	for _, notificationToken := range notificationTokens {
		msg.SendNotificationToNotificationToken(notificationToken)
	}

}

func (apprSvc *service) sendAppreciationNotificationToAll(ctx context.Context, appr repository.AppreciationResponse) {

	logger.Debug(ctx, " appreciationService appr: ", appr)
	msg := notification.Message{
		Title: "Appreciation",
		Body:  fmt.Sprintf(" %s %s has received an appreciation", appr.ReceiverFirstName, appr.ReceiverLastName),
	}
	logger.Infof(ctx, "appreciationService message: %v", msg)
	msg.SendNotificationToTopic("peerly")
}

func (apprSvc *service) sendEmailForBadgeAllocation(userBadgeDetails []repository.UserBadgeDetails) {

	logger.Debug(context.Background(), "appreciationService user Badge Details: ", userBadgeDetails)
	for _, userBadgeDetail := range userBadgeDetails {

		// Determine the BadgeImageUrl based on the BadgeName
		var badgeImageUrl string
		switch userBadgeDetail.BadgeName.String {
		case "Bronze":
			badgeImageUrl = fmt.Sprint(config.PeerlyBaseUrl() + constants.BronzeBadgeIconImagePath)
		case "Silver":
			badgeImageUrl = fmt.Sprint(config.PeerlyBaseUrl() + constants.SilverBadgeIconImagePath)
		case "Gold":
			badgeImageUrl = fmt.Sprint(config.PeerlyBaseUrl() + constants.GoldBadgeIconImagePath)
		case "Platinum":
			badgeImageUrl = fmt.Sprint(config.PeerlyBaseUrl() + constants.PlatinumIconImagePath)
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
		logger.Info(context.Background(), "appreciationService badge data: ", templateData)
		mailReq := email.NewMail([]string{userBadgeDetail.Email}, []string{}, []string{}, fmt.Sprintf("You've Bagged the %s for Crushing %d Points! 🏆", userBadgeDetail.BadgeName.String, userBadgeDetail.BadgePoints))
		err := mailReq.ParseTemplate("./internal/app/email/templates/badge.html", templateData)
		if err != nil {
			logger.Errorf(context.Background(), "appreciationService err in creating html file : %v", err)
			return
		}
		err = mailReq.Send()
		if err != nil {
			logger.Errorf(context.Background(), "appreciationService err: %v", err)
			return
		}
		logger.Infof(context.Background(), "appreciationService mail request: %v", mailReq)
	}
}
