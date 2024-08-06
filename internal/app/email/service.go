package email

import (
	// "fmt"
	"bytes"
	"fmt"
	"html/template"

	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	logger "github.com/sirupsen/logrus"
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

	senderEmail := config.ReadEnvString("SENDER_EMAIL")
	if senderEmail == "" {
		logger.Error("SENDER_EMAIL environment variable is not set")
		return fmt.Errorf("sender email not configured")
	}

	sendGridAPIKey := config.ReadEnvString("SENDGRID_API_KEY")
	if sendGridAPIKey == "" {
		logger.Error("SENDGRID_API_KEY environment variable is not set")
		return fmt.Errorf("sendgrid API key not configured")
	}

	fromEmail := mail.NewEmail("Peerly", senderEmail)
	toEmail := mail.NewEmail("Example User", ms.to[0])
	// cc1 := mail.NewEmail("Example CC", "sakshimokashi23@gmail.com")
	content := mail.NewContent("text/html", ms.body)
	if toEmail == nil {
		logger.Error("Recipient email is invalid")
		return fmt.Errorf("recipient email is invalid")
	}

	// create new *SGMailV3
	m := mail.NewV3Mail()
	m.SetFrom(fromEmail)
	m.AddContent(content)
	
	personalization := mail.NewPersonalization()
	personalization.AddTos(toEmail)
	// personalization.AddCCs(cc1)
	personalization.Subject = ms.subject
	m.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(sendGridAPIKey)

	response, err := client.Send(m)
	if err != nil {
		logger.Error("unable to send mail", "error", err)
		return err
	}

	logger.Info("Email sent successfully!")
	logger.Infof("Response status code: %v", response.StatusCode)
	logger.Infof("Response body: %v", response.Body)
	logger.Infof("Response headers: %v", response.Headers)
	return nil
}
func (r *Mail) ParseTemplate(templateFileName string, data interface{}) error {
	fmt.Println("--------------------->")
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Println("--------------------------->", err.Error())
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Println("--------------------------->", err.Error())
		return err
	}
	r.body = buf.String()
	fmt.Println("--------------------->")
	fmt.Println(r.body)
	fmt.Println("--------------------->")
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
