package validations

import (
	"testing"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
)
func TestOrgValidate(t *testing.T) {
	tests := []struct {
		name             string
		org              dto.Organization
		expectedValid    bool
		expectedResponse map[string]dto.ErrorResponse
	}{
		{
			name: "Valid Organization",
			org: dto.Organization{
				Name:                      "Valid Organization",
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "UTC",
			},
			expectedValid:    true,
			expectedResponse: nil,
		},
		{
			name: "Missing Name",
			org: dto.Organization{
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"name": "Can't be blank",
						},
					},
				},
			},
		},
		{
			name: "Invalid Email",
			org: dto.Organization{
				Name:                      "Invalid Email Org",
				ContactEmail:              "invalid-email",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"email": "Please enter a valid email",
						},
					},
				},
			},
		},
		{
			name: "Invalid Domain Name",
			org: dto.Organization{
				Name:                      "Invalid Domain Org",
				ContactEmail:              "valid@example.com",
				DomainName:                "invalid_domain",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"domain_name": "Please enter valid domain",
						},
					},
				},
			},
		},
		{
			name: "Past Subscription Valid Upto Date",
			org: dto.Organization{
				Name:                      "Past Subscription Org",
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(-24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"subscription_valid_upto": "Please enter subscription valid upto date",
						},
					},
				},
			},
		},
		{
			name: "Zero Hi5 Limit",
			org: dto.Organization{
				Name:                      "Zero Hi5 Limit Org",
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  0,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"hi5_limit": "Please enter hi5 limit greater than 0",
						},
					},
				},
			},
		},
		{
			name: "Invalid Hi5 Quota Renewal Frequency",
			org: dto.Organization{
				Name:                      "Invalid Hi5 Quota Org",
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "invalid_frequency",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"hi5_quota_renewal_frequency": "Please enter valid hi5 renewal frequency",
						},
					},
				},
			},
		},
		{
			name: "Invalid Timezone",
			org: dto.Organization{
				Name:                      "Invalid Timezone Org",
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5Limit:                  5,
				Hi5QuotaRenewalFrequency:  "month",
				Timezone:                  "",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code:          "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"timezone": "Please enter valid timezone",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorResponse, valid := OrgValidate(tt.org)

			assert.Equal(t, tt.expectedValid, valid)

			if !tt.expectedValid {
				// Check if the expected errors are present in the actual errors
				for key, expectedError := range tt.expectedResponse {
					actualError, exists := errorResponse[key]
					assert.True(t, exists, "Expected error key %v not found", key)
					assert.Equal(t, expectedError, actualError)
				}
			} else {
				assert.Nil(t, errorResponse)
			}
		})
	}
}

func TestOrgUpdateValidate(t *testing.T) {
	tests := []struct {
		name             string
		org              dto.Organization
		expectedValid    bool
		expectedResponse map[string]dto.ErrorResponse
	}{
		{
			name: "Valid Organization Update",
			org: dto.Organization{
				ID:                        1,
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5QuotaRenewalFrequency:  "week",
				Timezone:                  "UTC",
			},
			expectedValid:    true,
			expectedResponse: nil,
		},
		{
			name: "Invalid ID",
			org: dto.Organization{
				ID:                        0,
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5QuotaRenewalFrequency:  "week",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"id": "Please enter valid id",
						},
					},
				},
			},
		},
		{
			name: "Invalid Email",
			org: dto.Organization{
				ID:                        1,
				ContactEmail:              "invalid-email",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5QuotaRenewalFrequency:  "week",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"email": "Please enter a valid email",
						},
					},
				},
			},
		},
		{
			name: "Invalid Domain Name",
			org: dto.Organization{
				ID:                        1,
				ContactEmail:              "valid@example.com",
				DomainName:                "invalid_domain",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5QuotaRenewalFrequency:  "week",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"domain_name": "Please enter valid domain",
						},
					},
				},
			},
		},
		{
			name: "Past Subscription Valid Upto Date",
			org: dto.Organization{
				ID:                        1,
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(-24 * time.Hour),
				Hi5QuotaRenewalFrequency:  "week",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"subscription_valid_upto": "Please enter subscription valid upto date",
						},
					},
				},
			},
		},
		{
			name: "Invalid Hi5 Quota Renewal Frequency",
			org: dto.Organization{
				ID:                        1,
				ContactEmail:              "valid@example.com",
				DomainName:                "example.com",
				SubscriptionValidUpto:     time.Now().Add(24 * time.Hour),
				Hi5QuotaRenewalFrequency:  "invalid_frequency",
				Timezone:                  "UTC",
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"hi5_quota_renewal_frequency": "Please enter valid hi5 renewal frequency",
						},
					},
				},
			},
		},
		{
			name: "Invalid Timezone",
			org: dto.Organization{
				ID:                        1,
				Timezone:                  "",
			},
			expectedValid: true,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid organization data"},
						Fields: map[string]string{
							"timezone": "Please enter valid timezone",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorResponse, valid := OrgUpdateValidate(tt.org)

			assert.Equal(t, tt.expectedValid, valid)

			if !tt.expectedValid {
				for key, expectedError := range tt.expectedResponse {
					actualError, exists := errorResponse[key]
					assert.True(t, exists, "Expected error key %v not found", key)
					assert.Equal(t, expectedError, actualError)
				}
			} else {
				assert.Nil(t, errorResponse)
			}
		})
	}
}

func TestOTPInfoValidate(t *testing.T) {
	tests := []struct {
		name             string
		otp              dto.OTP
		expectedValid    bool
		expectedResponse map[string]dto.ErrorResponse
	}{
		{
			name: "Valid OTP Info",
			otp: dto.OTP{
				OTPCode: "123456",
				OrgId:   1,
			},
			expectedValid:    true,
			expectedResponse: nil,
		},
		{
			name: "Invalid OTP Code Length",
			otp: dto.OTP{
				OTPCode: "12345",
				OrgId:   1,
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid otp data"},
						Fields: map[string]string{
							"otp_code": "enter 6 digit valid otp code",
						},
					},
				},
			},
		},
		{
			name: "Invalid Org ID",
			otp: dto.OTP{
				OTPCode: "123456",
				OrgId:   0,
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid otp data"},
						Fields: map[string]string{
							"id": "Please enter valid organization id",
						},
					},
				},
			},
		},
		{
			name: "Invalid OTP Code Length and Org ID",
			otp: dto.OTP{
				OTPCode: "12345",
				OrgId:   0,
			},
			expectedValid: false,
			expectedResponse: map[string]dto.ErrorResponse{
				"error": {
					Error: dto.ErrorObject{
						Code: "invalid_data",
						MessageObject: dto.MessageObject{Message: "Please provide valid otp data"},
						Fields: map[string]string{
							"otp_code": "enter 6 digit valid otp code",
							"id":       "Please enter valid organization id",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorResponse, valid := OTPInfoValidate(tt.otp)

			assert.Equal(t, tt.expectedValid, valid)

			if !tt.expectedValid {
				for key, expectedError := range tt.expectedResponse {
					actualError, exists := errorResponse[key]
					assert.True(t, exists, "Expected error key %v not found", key)
					assert.Equal(t, expectedError, actualError)
				}
			} else {
				assert.Nil(t, errorResponse)
			}
		})
	}
}