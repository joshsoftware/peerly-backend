package appreciation

import (
	"context"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"fmt"

	// "io/ioutil"

	"github.com/joshsoftware/peerly-backend/internal/app/email"

	// "github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/utils"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	logger "github.com/sirupsen/logrus"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type service struct {
	appreciationRepo repository.AppreciationStorer
	corevaluesRespo  repository.CoreValueStorer
}

// Service contains all
type Service interface {
CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error)
GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error)
ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error)
DeleteAppreciation(ctx context.Context, apprId int32) error
	UpdateAppreciation(ctx context.Context) (bool, error)
	sendAppreciationEmail(to string, sub string, maildata string) error
}

func NewService(appreciationRepo repository.AppreciationStorer, coreValuesRepo repository.CoreValueStorer) Service {
	return &service{
		appreciationRepo: appreciationRepo,
		corevaluesRespo:  coreValuesRepo,
	}
}

func (apprSvc *service) CreateAppreciation(ctx context.Context, appreciation dto.Appreciation) (dto.Appreciation, error) {

	//add quarter
	appreciation.Quarter = utils.GetQuarter()

	//add sender
	data := ctx.Value(constants.UserId)
	sender, ok := data.(int64)
	if !ok {
		logger.Error("err in parsing userid from token")
		return dto.Appreciation{}, apperrors.InternalServer
	}

	//check is receiver present in database
	chk, err := apprSvc.appreciationRepo.IsUserPresent(ctx, nil, appreciation.Receiver)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.Appreciation{}, err
	}
	if !chk {
		return dto.Appreciation{}, apperrors.UserNotFound
	}
	appreciation.Sender = sender

	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		return dto.Appreciation{}, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Info(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := apprSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Info(ctx, "error in creating transaction, err: %s", txErr.Error())
			return
		}
	}()

	//check is corevalue present in database
	_, err = apprSvc.corevaluesRespo.GetCoreValue(ctx, int64(appreciation.CoreValueID))
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.Appreciation{}, err
	}

	// check self appreciation
	if appreciation.Receiver == sender {
		return dto.Appreciation{}, apperrors.SelfAppreciationError
	}

	appr, err := apprSvc.appreciationRepo.CreateAppreciation(ctx, tx, appreciation)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.Appreciation{}, err
	}

	// to := []string{"samnitpatil9882@gmail.com"}
	// sub := "OTP Verification"
	// mailData := "samnit patil sp email"
	// mailReq := email.NewMail(config.ReadEnvString("SENDER_EMAIL"), to, sub, mailData)
	// err = mailReq.SendOTPMail()
	// if err != nil{
	// 	logger.Errorf("email err: %v",err)
	// }
	// fmt.Println("----------------------> successfull")

	getAppr,err := apprSvc.appreciationRepo.GetAppreciationById(ctx,tx,int32(appr.ID))
	if err != nil {
		logger.Errorf("get by id in create  err: %v", err)
	}
	apprSvc.sendNotification(getAppr)
	err = apprSvc.sendAppreciationEmail("samnitpatil9882@gmail.com", "Received appreciation", "samnit patil peerly")
	if err != nil {
		logger.Errorf("email err: %v", err)
	}

	return mapAppreciationDBToDTO(appr), nil
}

func (apprSvc *service) GetAppreciationById(ctx context.Context, appreciationId int32) (dto.AppreciationResponse, error) {

	resAppr, err := apprSvc.appreciationRepo.GetAppreciationById(ctx, nil, appreciationId)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.AppreciationResponse{}, err
	}

	return mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(resAppr), nil
}

func (apprSvc *service) ListAppreciations(ctx context.Context, filter dto.AppreciationFilter) (dto.ListAppreciationsResponse, error) {

	infos, pagination, err := apprSvc.appreciationRepo.ListAppreciations(ctx, nil, filter)
	if err != nil {
		logger.Errorf("err: %v", err)
		return dto.ListAppreciationsResponse{}, err
	}

	responses := make([]dto.AppreciationResponse, 0)
	for _, info := range infos {
		responses = append(responses, mapRepoGetAppreciationInfoToDTOGetAppreciationInfo(info))
	}
	paginationResp := dtoPagination(pagination)
	return dto.ListAppreciationsResponse{Appreciations: responses, MetaData: paginationResp}, nil
}

func (apprSvc *service) DeleteAppreciation(ctx context.Context, apprId int32) error {
	return apprSvc.appreciationRepo.DeleteAppreciation(ctx, nil, apprId)
}

func (apprSvc *service) UpdateAppreciation(ctx context.Context) (bool, error) {

	//initializing database transaction
	tx, err := apprSvc.appreciationRepo.BeginTx(ctx)

	if err != nil {
		return false, err
	}

	defer func() {
		rvr := recover()
		defer func() {
			if rvr != nil {
				logger.Info(ctx, "Transaction aborted because of panic: %v, Propagating panic further", rvr)
				panic(rvr)
			}
		}()

		txErr := apprSvc.appreciationRepo.HandleTransaction(ctx, tx, err == nil && rvr == nil)
		if txErr != nil {
			err = txErr
			logger.Info(ctx, "error in creating transaction, err: %s", txErr.Error())
			return
		}
	}()

	_, err = apprSvc.appreciationRepo.UpdateAppreciationTotalRewardsOfYesterday(ctx, tx)

	if err != nil {
		logger.Error("err: ", err.Error())
		return false, err
	}

	_, err = apprSvc.appreciationRepo.UpdateUserBadgesBasedOnTotalRewards(ctx, tx)

	if err != nil {
		logger.Error("err: ", err.Error())
		return false, err
	}

	return true, nil
}
func (apprSvc *service) sendAppreciationEmail(to string, sub string, maildata string) error {
	// Plain text content
	plainTextContent := "Samnit " + "123456"

	templateData := struct {
		ReceiverName string
		SenderName   string
		Quarter      int8
		Message      string
		YourName     string
		YourPosition string
	}{
		ReceiverName: "samnit patil",
		SenderName:   "samir patil",
		Quarter:      2,
		Message:      "Good TeamWork",
		YourName:     "peerly",
		YourPosition: "hr",
	}
	// ReceiverName := "samnit patil"
	// SenderName := "samir patil"
	// Quarter := 2
	// Message := "Good TeamWork"
	// YourName := "peerly"
	// YourPosition := "hr"
	// CompanyName := "josh software"
	// htmlContent := `
    // <!DOCTYPE html>
    // <html>
    // <head>
    //     <meta charset="UTF-8">
    //     <style>
    //         body {
    //             font-family: Arial, sans-serif;
    //             color: #333333;
    //             line-height: 1.6;
    //         }
    //         .content {
    //             padding: 20px;
    //             background-color: #f9f9f9;
    //             border: 1px solid #dddddd;
    //             border-radius: 5px;
    //         }
    //         .header {
    //             font-size: 18px;
    //             font-weight: bold;
    //             margin-bottom: 20px;
    //         }
    //         .footer {
    //             margin-top: 30px;
    //             font-size: 14px;
    //             color: #555555;
    //         }
    //     </style>
    // </head>
    // <body>
    //     <div class="content">
    //         <p class="header">Dear ` + ReceiverName + `,</p>
    //         <p>We're pleased to inform you that you've received an appreciation for your work!</p>
    //         <p><strong>Appreciation Details:</strong></p>
    //         <ul>
    //             <li><strong>From:</strong> ` + SenderName + `</li>
    //             <li><strong>Message:</strong> ` + Message + `</li>
    //             <li><strong>Quarter:</strong> Q` + fmt.Sprintf("%d", Quarter) + `</li>
    //         </ul>
    //         <p>Thank you for your valuable contributions to the team.</p>
    //         <p class="footer">
    //             Best regards,<br>
    //             ` + YourName + `, ` + YourPosition + `<br>
    //             ` + CompanyName + `
    //         </p>
    //     </div>
    // </body>
    // </html>
    // `

	mailReq := email.NewMail([]string{to}, []string{"samnitpatil@gmail.com"}, []string{"samirpatil9882@gmail.com"}, sub)
	mailReq.ParseTemplate("./internal/app/email/templates/createAppreciation.html", templateData)
	err := mailReq.Send(plainTextContent)
	if err != nil {
		logger.Errorf("err: %v", err)
		return err
	}

	return nil
}

func (apprSvc *service) sendNotification(appr repository.AppreciationResponse) {
	// Path to your service account key file
	serviceAccountKey := "serviceAccountKey.json"

	// Initialize the Firebase app
	opt := option.WithCredentialsFile(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf("Error initializing app: %v", err)
		return
	}

	// Obtain a messaging client from the Firebase app
	client, err := app.Messaging(context.Background())
	if err != nil {
		logger.Errorf("Error getting Messaging client: %v", err)
		return
	}

	fmt.Println("---------------------------->")
	fmt.Println("appr info: ",appr)
	fmt.Printf("%s %s appreciated %s %s", appr.SenderFirstName, appr.SenderLastName, appr.ReceiverFirstName, appr.ReceiverLastName)

	// Create a message to send
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Peerly",
			Body:  fmt.Sprintf("%s %s appreciated %s %s", appr.SenderFirstName, appr.SenderLastName, appr.ReceiverFirstName, appr.ReceiverLastName),
		},
		Topic: "appreciation",
		// Token: "dUwh4EYPT1SGkNW-0zh8gn:APA91bE-rLoeDPuzMORjQtqyMAGbaB75yK0g0f-BGNF2qnq343d4Iih1Kh2lNeZdeILF9oXA0RPT3foaDT8CFG1FcZcTbOOj6JoPUdfnBPFzb4kVWppUHl354JzHduATpEx-pweb40pa", // Replace with a secure method to retrieve the token
	}

	// Send the message
	response, err := client.Send(context.Background(), message)
	if err != nil {
		logger.Errorf("Error sending message: %v", err)
		return
	}

	// Print the response
	logger.Infof("Successfully sent message: %v", response)
}

