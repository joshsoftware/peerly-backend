package validation

import (
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
)


func TestOrgValidate(t *testing.T) {
	tests := []struct {
		name           string
		organization   dto.OrganizationConfig
		expectedErrors map[string]string
		expectedValid  bool
	}{
		{
			name: "Valid organization config",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "UTC",
			},
			expectedErrors: nil,
			expectedValid:  true,
		},
		{
			name: "Invalid RewardMultiplier",
			organization: dto.OrganizationConfig{
				RewardMultiplier: 0,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			expectedErrors: map[string]string{
				"reward_multiplier": "Please enter reward multiplier greater than 0",
			},
			expectedValid: false,
		},
		{
			name: "Invalid RewardQuotaRenewalFrequency",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            10,
				RewardQuotaRenewalFrequency: 0,
				Timezone:                    "UTC",
			},
			expectedErrors: map[string]string{
				"reward_quota_renewal_frequency": "Please enter valid reward renewal frequency",
			},
			expectedValid: false,
		},
		{
			name: "Invalid Timezone",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "Invalid Timezone",
			},
			expectedErrors: map[string]string{
				"timezone": "Please enter valid timezone",
			},
			expectedValid: false,
		},
		{
			name: "Multiple Invalid Fields",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            0,
				RewardQuotaRenewalFrequency: 0,
				Timezone:                    "Invalid/Timezone",
			},
			expectedErrors: map[string]string{
				"reward_multiplier":                     "Please enter reward multiplier greater than 0",
				"reward_quota_renewal_frequency": "Please enter valid reward renewal frequency",
				"timezone":                      "Please enter valid timezone",
			},
			expectedValid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errorResponse, valid := OrgValidate(test.organization)

			assert.Equal(t, test.expectedValid, valid)

			if !test.expectedValid {
				// Check if the expected errors are present in the actual errors
				assert.NotNil(t, errorResponse["error"], "Expected 'error' key in error response")

				errorObject := errorResponse["error"].Error.(dto.ErrorObject)
				assert.Equal(t, "invalid_data", errorObject.Code)
				assert.Equal(t, "Please provide valid organization data", errorObject.Message)

				for key, expectedMessage := range test.expectedErrors {
					actualMessage, exists := errorObject.Fields[key]
					assert.True(t, exists, "Expected error key %v not found", key)
					assert.Equal(t, expectedMessage, actualMessage)
				}
			} else {
				assert.Nil(t, errorResponse)
			}
		})
	}
}

func TestOrgUpdateValidate(t *testing.T) {
	tests := []struct {
		name           string
		organization   dto.OrganizationConfig
		expectedErrors map[string]string
		expectedValid  bool
	}{
		{
			name: "Valid organization config",
			organization: dto.OrganizationConfig{
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                     "UTC",
			},
			expectedErrors: nil,
			expectedValid:  true,
		},
		{
			name: "Invalid Timezone",
			organization: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                     "Invalid/Timezone",
			},
			expectedErrors: map[string]string{
				"timezone": "Please enter valid timezone",
			},
			expectedValid: false,
		},

	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errorResponse, valid := OrgUpdateValidate(test.organization)

			assert.Equal(t, test.expectedValid, valid)

			if !test.expectedValid {
				// Check if the expected errors are present in the actual errors
				assert.NotNil(t, errorResponse["error"], "Expected 'error' key in error response")

				errorObject := errorResponse["error"].Error.(dto.ErrorObject)
				assert.Equal(t, "invalid_data", errorObject.Code)
				assert.Equal(t, "Please provide valid organization data", errorObject.Message)

				for key, expectedMessage := range test.expectedErrors {
					actualMessage, exists := errorObject.Fields[key]
					assert.True(t, exists, "Expected error key %v not found", key)
					assert.Equal(t, expectedMessage, actualMessage)
				}
			} else {
				assert.Nil(t, errorResponse)
			}
		})
	}
}
