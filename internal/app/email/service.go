package email

import (
	"context"
	"fmt"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	logger "github.com/joshsoftware/peerly-backend/internal/pkg/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// MailService represents the interface for our mail service.
type MailService interface {
	Send() error
}

// Mail represents a email request
type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	cc      []string
	bcc     []string
}

func (ms *Mail) Send() error {

	logger.Info(context.Background(), " Mail: ", ms)

	sendGridAPIKey := config.ReadEnvString("SENDGRID_API_KEY")
	if sendGridAPIKey == "" {
		logger.Error(context.Background(), "SENDGRID_API_KEY environment variable is not set")
		return fmt.Errorf("sendgrid API key not configured")
	}

	err := ms.validateMail()
	if err != nil {
		logger.Errorf(context.Background(),"err: error in mail validation: %v",err)
		return err
	}
	fromEmail := mail.NewEmail("Peerly", ms.from)
	content := mail.NewContent("text/html", ms.body)

	// create new *SGMailV3
	m := mail.NewV3Mail()
	m.SetFrom(fromEmail)
	m.AddContent(content)

	personalization := mail.NewPersonalization()

	for _, email := range ms.to {
		toEmail := mail.NewEmail("to", email)
		personalization.AddTos(toEmail)
	}

	for _, email := range ms.cc {
		ccEmail := mail.NewEmail("cc", email)
		personalization.AddCCs(ccEmail)
	}

	for _, email := range ms.bcc {
		bccEmail := mail.NewEmail("bcc", email)
		personalization.AddBCCs(bccEmail)
	}

	personalization.Subject = ms.subject
	m.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(sendGridAPIKey)

	response, err := client.Send(m)
	if err != nil {
		logger.Error(context.Background(), "unable to send mail", "error", err)
		return err
	}

	logger.Infof(context.Background(),"email sent successfully to %v ",ms.to)
	logger.Debug(context.Background(), "Email response: %v ",response)
	return nil
}

func (ms *Mail) validateMail() error {
    if len(ms.to) == 0 {
        return apperrors.InvalidTos
    }

    if len(ms.from) == 0 {
        return apperrors.InvalidFrom
    }

    if len(ms.body) == 0 {
        return apperrors.InvalidBody
    }

    if len(ms.subject) == 0 {
        return apperrors.InvalidSub
    }

    return nil
}

// NewMail returns a new mail request.
func NewMail(to []string, cc []string, bcc []string, subject string, body string) MailService {
	return &Mail{
		from:    config.ReadEnvString("SENDER_EMAIL"),
		to:      to,
		cc:      cc,
		bcc:     bcc,
		subject: subject,
		body:    body,
	}
}
