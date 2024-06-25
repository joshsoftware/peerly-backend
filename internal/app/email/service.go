package email

import (
	// "fmt"
	"fmt"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	logger "github.com/sirupsen/logrus"
)

// MailService represents the interface for our mail service.
type MailService interface {
	// CreateMail(mailReq *Mail) []byte
	SendMail(mailReq *Mail) error
}

// MailData represents the data to be sent to the template of the mail.
type MailData struct {
	OTPCode     string
}

// Mail represents a email request
type Mail struct {
	from    string
	to      []string
	subject string
	data    *MailData
}

// func (ms *Mail) CreateMail() ([]byte) {
//     m := mail.NewV3Mail()

//     from := mail.NewEmail("OTP Verification", "samirpatil9882@gmail.com")
//     m.SetFrom(from)

//     // Since you want to send to only one email, you can take the first one from ms.to
//     toEmail := ms.to[0]
//     to := mail.NewEmail("user", toEmail)

//     p := mail.NewPersonalization()
//     p.AddTos(to)

//     // Add custom substitutions (if needed)
//     p.SetSubstitution("Username", ms.data.ReferenceId)
//     p.SetSubstitution("Code", ms.data.OTPCode)

//     m.AddPersonalizations(p)

//     // Create content for the email
//     content := mail.NewContent("text/plain", "Hello, here is your OTP code: "+ms.data.OTPCode)
//     m.AddContent(content)

//     // Generate request body
//     requestBody := mail.GetRequestBody(m)
//     return requestBody
// }

// SendMail creates a sendgrid mail from the given mail request and sends it.
func (ms *Mail) SendMail() error {
	from := mail.NewEmail("Organization", ms.from)
	subject := "Peerly: OTP verification"
	to := mail.NewEmail("Example User", ms.to[0])

	// Plain text content
	plainTextContent := "Verify your account with the OTP: " + ms.data.OTPCode 

	// HTML content with inline styling
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>Email Verification</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					line-height: 1.6;
					background-color: #f0f0f0;
					padding: 20px;
				}
				.content {
					background-color: #ffffff;
					padding: 20px;
					border-radius: 5px;
					box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
				}
				.otp-code {
					font-size: 18px;
					font-weight: bold;
					color: #333333;
				}
			</style>
		</head>
		<body>
			<div class="content">
				<p><strong>Verify your account with the OTP:</strong></p>
				<p><span class="otp-label">OTP:</span> <span class="otp-code">` + ms.data.OTPCode + `</span></p>
			</div>
		</body>
		</html>
	`

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(config.ReadEnvString("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		logger.Error("unable to send mail", "error", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	fmt.Println("Response status code:", response.StatusCode)
	fmt.Println("Response body:", response.Body)
	fmt.Println("Response headers:", response.Headers)

	return nil
}

// NewMail returns a new mail request.
func NewMail(from string, to []string, subject string, data *MailData) *Mail {
	return &Mail{
		from:    from,
		to:      to,
		subject: subject,
		data:    data,
	}
}
