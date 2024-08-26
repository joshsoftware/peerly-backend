package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// MailService represents the interface for our mail service.
type MailService interface {
	Send(plainTextContent string) error
	ParseTemplate(templateFileName string, data interface{}) error
}

// Mail represents a email request
type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	CC      []string
	BCC     []string
}

func (ms *Mail) Send(plainTextContent string) error {

	logger.Info(context.Background()," Mail: ",ms)
	senderEmail := config.ReadEnvString("SENDER_EMAIL")
	if senderEmail == "" {
		logger.Error(context.Background(),"SENDER_EMAIL environment variable is not set")
		return fmt.Errorf("sender email not configured")
	}

	sendGridAPIKey := config.ReadEnvString("SENDGRID_API_KEY")
	if sendGridAPIKey == "" {
		logger.Error(context.Background(),"SENDGRID_API_KEY environment variable is not set")
		return fmt.Errorf("sendgrid API key not configured")
	}
	logger.Info(context.Background(),"from_------------->, ",senderEmail)

	fromEmail := mail.NewEmail("Peerly", senderEmail)
	content := mail.NewContent("text/html", ms.body)

	// create new *SGMailV3
	m := mail.NewV3Mail()
	m.SetFrom(fromEmail)
	m.AddContent(content)

	personalization := mail.NewPersonalization()

	for _,email := range ms.to{
		toEmail := mail.NewEmail("to", email)
		personalization.AddTos(toEmail)
	}

	for _,email := range ms.CC{
		ccEmail := mail.NewEmail("cc", email)
		personalization.AddCCs(ccEmail)
	}

	for _,email := range ms.BCC{
		bccEmail := mail.NewEmail("bcc", email)
		personalization.AddBCCs(bccEmail)
	}
	
	personalization.Subject = ms.subject
	m.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(sendGridAPIKey)

	response, err := client.Send(m)
	if err != nil {
		logger.Error(context.Background(),"unable to send mail", "error", err)
		return err
	}

	logger.Info(context.Background(),"Email response: ")
	logger.Infof(context.Background(),"Response status code: %v", response.StatusCode)
	logger.Infof(context.Background(),"Response body: %v", response.Body)
	logger.Infof(context.Background(),"Response headers: %v", response.Headers)
	return nil
}
func (r *Mail) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	logger.Info(context.Background(),"--------------------->")
	logger.Info(context.Background(),r.body)
	logger.Info(context.Background(),"--------------------->")
	return nil
}

// NewMail returns a new mail request.
func NewMail(to []string, cc []string, bcc []string, subject string) MailService {
	return &Mail{
		from:    config.ReadEnvString("SENDER_EMAIL"),
		to:      to,
		CC:      cc,
		BCC:     bcc,
		subject: subject,
	}
}
