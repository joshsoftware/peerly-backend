package notification

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"google.golang.org/api/option"
)

type NotificationService interface {
	SendNotificationToNotificationToken(notificationToken string) (err error)
	SendNotificationToTopic(topic string) (err error)
	// SendNotificationToNotificationTokens(userId int64)
}

type Message struct {
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	ImageURL string `json:"image,omitempty"`
}

func (notificationSvc *Message) SendNotificationToNotificationToken(notificationToken string) (err error) {

	// Path to your service account key file
	serviceAccountKey := "serviceAccountKey.json"

	// Initialize the Firebase app
	opt := option.WithCredentialsFile(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf(context.Background(), "Error initializing app: %v", err)
		err = apperrors.InternalServerError
		return
	}

	// Obtain a messaging client from the Firebase app
	client, err := app.Messaging(context.Background())
	if err != nil {
		logger.Errorf(context.Background(), "Error getting Messaging client: %v", err)
		err = apperrors.InternalServerError
		return
	}

	logger.Debug(context.Background(), " notificationSvc: ", notificationSvc)
	// Create a message to send
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notificationSvc.Title,
			Body:  notificationSvc.Body,
		},
		Token: notificationToken,
	}

	// Send the message
	response, err := client.Send(context.Background(), message)
	logger.Debug(context.Background(), " response: ", response)
	logger.Debug(context.Background(), " err: ", err)
	if err != nil {
		logger.Errorf(context.Background(), "Error sending message: %v", err)
		err = apperrors.InternalServerError
		return
	}
	logger.Infof(context.Background(), "Successfully sent message: %v", response)
	return
}

func (notificationSvc *Message) SendNotificationToTopic(topic string) (err error) {

	// Path to your service account key file
	serviceAccountKey := "serviceAccountKey.json"

	// Initialize the Firebase app
	opt := option.WithCredentialsFile(serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Errorf(context.Background(), "error initializing app: %v", err)
		err = apperrors.InternalServerError
		return
	}

	// Obtain a messaging client from the Firebase app
	client, err := app.Messaging(context.Background())
	if err != nil {
		logger.Errorf(context.Background(), "error getting Messaging client: %v", err)
		err = apperrors.InternalServerError
		return
	}

	logger.Debug(context.Background(), " notificationSvc: ", notificationSvc)
	// Create a message to send
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: notificationSvc.Title,
			Body:  notificationSvc.Body,
		},
		Topic: topic,
	}

	// Send the message
	response, err := client.Send(context.Background(), message)
	if err != nil {
		logger.Errorf(context.Background(), "error sending message: %v", err)
		err = apperrors.InternalServerError
		return
	}

	logger.Infof(context.Background(), "Successfully sent message: %v", response)

	return
}

// func (notificationSvc *Message) SendNotificationToNotificationTokens(userId int64) {

// 	// Path to your service account key file
// 	serviceAccountKey := "serviceAccountKey.json"

// 	// Initialize the Firebase app
// 	opt := option.WithCredentialsFile(serviceAccountKey)
// 	app, err := firebase.NewApp(context.Background(), nil, opt)
// 	if err != nil {
// 		logger.Errorf("Error initializing app: %v", err)
// 		return
// 	}

// 	// Obtain a messaging client from the Firebase app
// 	client, err := app.Messaging(context.Background())
// 	if err != nil {
// 		logger.Errorf("Error getting Messaging client: %v", err)
// 		return
// 	}
// 	// Create a message to send
// 	message := &messaging.Message{
// 		Notification: &messaging.Notification{
// 			Title: notificationSvc.Title,
// 			Body:  notificationSvc.Body,
// 		},
// 		Token: notificationToken,
// 	}

// 	// Send the message
// 	response, err := client.Send(context.Background(), message)
// 	if err != nil {
// 		logger.Errorf("Error sending message: %v", err)
// 		return
// 	}
// 	logger.Infof("Successfully sent message: %v", response)
// }
