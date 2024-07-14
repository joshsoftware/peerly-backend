package validation

import (
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/stretchr/testify/assert"
)


func TestOrgValidate(t *testing.T) {
	tests := []struct {
		name           string
		organization   dto.OrganizationConfig
		expectedErrors error
	}{
		{
			name: "Valid organization config",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "UTC",
			},
			expectedErrors: nil,
		},
		{
			name: "Invalid RewardMultiplier",
			organization: dto.OrganizationConfig{
				RewardMultiplier: 0,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			expectedErrors: apperrors.InvalidRewardMultiplier,
		},
		{
			name: "Invalid RewardQuotaRenewalFrequency",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            10,
				RewardQuotaRenewalFrequency: 0,
				Timezone:                    "UTC",
			},
			expectedErrors: apperrors.InvalidRewardQuotaRenewalFrequency,
		},
		{
			name: "Invalid Timezone",
			organization: dto.OrganizationConfig{
				RewardMultiplier:            10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "Invalid Timezone",
			},
			expectedErrors: apperrors.InvalidTimezone,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := OrgValidate(test.organization)

			assert.Equal(t, test.expectedErrors,err)

		})
	}
}

func TestOrgUpdateValidate(t *testing.T) {
	tests := []struct {
		name           string
		organization   dto.OrganizationConfig
		expectedErrors error
	}{
		{
			name: "Valid organization config",
			organization: dto.OrganizationConfig{
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                     "UTC",
			},
			expectedErrors: nil,
		},
		{
			name: "Invalid Timezone",
			organization: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                     "Invalid/Timezone",
			},
			expectedErrors: apperrors.InvalidTimezone,
		},

	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := OrgUpdateValidate(test.organization)

			assert.Equal(t, test.expectedErrors,err)
		})
	}
}
