package email

import (
	// "testing"

	// "github.com/stretchr/testify/assert"
)

// import "github.com/sendgrid/sendgrid-go"

// func TestSendMail(t *testing.T) {
// 	tests := []struct {
// 		name            string
// 		mockResponse    *sendgrid.Response
// 		mockError       error
// 		expectedError   error
// 		expectedSuccess bool
// 	}{
// 		{
// 			name: "Success",
// 			mockResponse: &sendgrid.Response{
// 				StatusCode: 202,
// 				Body:       "Accepted",
// 				Headers:    map[string][]string{},
// 			},
// 			mockError:       nil,
// 			expectedError:   nil,
// 			expectedSuccess: true,
// 		},
// 		{
// 			name: "Send Error",
// 			mockResponse:    nil,
// 			mockError:       errors.New("unable to send mail"),
// 			expectedError:   errors.New("unable to send mail"),
// 			expectedSuccess: false,
// 		},
// 		{
// 			name: "Invalid API Key",
// 			mockResponse: &sendgrid.Response{
// 				StatusCode: 401,
// 				Body:       "Unauthorized",
// 				Headers:    map[string][]string{},
// 			},
// 			mockError:       nil,
// 			expectedError:   errors.New("SendGrid API error: Unauthorized"),
// 			expectedSuccess: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockClient := new(MockSendGridClient)
// 			mockClient.On("Send", mock.Anything).Return(tt.mockResponse, tt.mockError).Once()

// 			mail := &Mail{
// 				from: "sender@example.com",
// 				to:   []string{"recipient@example.com"},
// 				data: &MailData{
// 					OTPCode: "123456",
// 				},
// 			}

// 			sendgrid.NewSendClient = func(apiKey string) sendgrid.Client {
// 				return mockClient
// 			}

// 			err := mail.SendMail()

// 			if tt.expectedSuccess {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.Error(t, err)
// 				assert.Equal(t, tt.expectedError, err)
// 			}

// 			mockClient.AssertExpectations(t)
// 		})
// 	}
// }