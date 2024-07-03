package organizationConfig

import (
	"context"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOrganizationConfig(t *testing.T) {
	orgRepo := mocks.NewOrganizationStorer(t)
	service := NewService(orgRepo)

	tests := []struct {
		name           string
		context        context.Context
		setup          func(orgMock *mocks.OrganizationStorer)
		expectedResult dto.OrganizationConfig
		expectedError  error
	}{
		{
			name:    "Successful retrieval of organization config",
			context: context.WithValue(context.Background(), "userId", int64(1)),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{
					ID:                           1,
					RewardMultiplier:             200,
					RewardQuotaRenewalFrequency: 12,
					Timezone:                    "ACT",
					CreatedAt:                   1719918501194,
					CreatedBy:                   7,
					UpdatedAt:                   1719920402224,
					UpdatedBy:                   7,
				}, nil).Once()
			},
			expectedResult: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             200,
				RewardQuotaRenewalFrequency: 12,
				Timezone:                    "ACT",
				CreatedAt:                   1719918501194,
				CreatedBy:                   7,
				UpdatedAt:                   1719920402224,
				UpdatedBy:                   7,
			},
			expectedError: nil,
		},
		{
			name:    "Error while retrieving organization config",
			context: context.WithValue(context.Background(), "userId", int64(1)),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{}, apperrors.InternalServer).Once()
			},
			expectedResult: dto.OrganizationConfig{},
			expectedError:  apperrors.InternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(orgRepo)
			result, err := service.GetOrganizationConfig(tt.context)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			orgRepo.AssertExpectations(t)
		})
	}
}

func TestCreateOrganizationConfig(t *testing.T) {
	orgRepo := mocks.NewOrganizationStorer(t)
	orgSvc := NewService(orgRepo)

	tests := []struct {
		name              string
		context           context.Context
		organizationInput dto.OrganizationConfig
		setup func(orgMock *mocks.OrganizationStorer)
		expectedResult    dto.OrganizationConfig
		expectedError     error
	}{
		{
			name: "Successful organization config creation",
			context: context.WithValue(context.Background(), "userId", 1),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{}, apperrors.OrganizationNotFound).Once()
				orgMock.On("CreateOrganizationConfig", mock.Anything,mock.Anything).Return(repository.OrganizationConfig{
					ID:                           1,
					RewardMultiplier:             200,
					RewardQuotaRenewalFrequency: 12,
					Timezone:                    "ACT",
					CreatedAt:                   1719918501194,
					CreatedBy:                   7,
					UpdatedAt:                   1719920402224,
					UpdatedBy:                   7,
				}, nil).Once()

			},
			organizationInput: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			expectedResult: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             200,
				RewardQuotaRenewalFrequency: 12,
				Timezone:                    "ACT",
				CreatedAt:                   1719918501194,
				CreatedBy:                   7,
				UpdatedAt:                   1719920402224,
				UpdatedBy:                   7,
			},
			expectedError: nil,
		},
		{
			name: "Organization config already present",
			context: context.WithValue(context.Background(), "userId", 1),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{
					ID:                           1,
					RewardMultiplier:             200,
					RewardQuotaRenewalFrequency: 12,
					Timezone:                    "ACT",
					CreatedAt:                   1719918501194,
					CreatedBy:                   7,
					UpdatedAt:                   1719920402224,
					UpdatedBy:                   7,
				}, nil).Once()
			},
			organizationInput: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			expectedResult:  dto.OrganizationConfig{},
			expectedError:   apperrors.OrganizationConfigAlreadyPresent,
		},
		{
			name: "Error while creating organization config",
			context: context.WithValue(context.Background(), "userId", 1),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{}, apperrors.OrganizationNotFound).Once()
				orgMock.On("CreateOrganizationConfig", mock.Anything,mock.Anything).Return(repository.OrganizationConfig{}, apperrors.InernalServer).Once()
			},
			organizationInput: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			expectedResult:  dto.OrganizationConfig{},
			expectedError:   apperrors.InternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(orgRepo)
			// Call the service method
			result, err := orgSvc.CreateOrganizationConfig(tt.context, tt.organizationInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Assert that the mocks were called as expected
			orgRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateOrganizationConfig(t *testing.T) {
	orgRepo := mocks.NewOrganizationStorer(t)
	orgSvc := NewService(orgRepo)

	tests := []struct {
		name              string
		context           context.Context
		organizationInput dto.OrganizationConfig
		setup             func(orgMock *mocks.OrganizationStorer)
		expectedResult    dto.OrganizationConfig
		expectedError     error
	}{
		{
			name: "Successful organization config update",
			context: context.WithValue(context.Background(), "userId", 1),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{
					ID:                           1,
					RewardMultiplier:             200,
					RewardQuotaRenewalFrequency: 12,
					Timezone:                    "ACT",
					CreatedAt:                   1719918501194,
					CreatedBy:                   7,
					UpdatedAt:                   1719920402224,
					UpdatedBy:                   7,
				}, nil).Once()
				orgMock.On("UpdateOrganizationCofig", mock.Anything, mock.Anything).Return(repository.OrganizationConfig{
					ID:                           1,
					RewardMultiplier:             10,
					RewardQuotaRenewalFrequency: 5,
					Timezone:                    "UTC",
					CreatedAt:                   1719918501194,
					CreatedBy:                   7,
					UpdatedAt:                   1719920402224,
					UpdatedBy:                   1, // Updated with user ID
				}, nil).Once()
			},
			organizationInput: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "UTC",
			},
			expectedResult: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "UTC",
				CreatedAt:                   1719918501194,
				CreatedBy:                   7,
				UpdatedAt:                   1719920402224,
				UpdatedBy:                   1,
			},
			expectedError: nil,
		},
		{
			name: "Organization config not found",
			context: context.WithValue(context.Background(), "userId", 1),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{}, apperrors.OrganizationNotFound).Once()
			},
			organizationInput: dto.OrganizationConfig{
				ID: 1,
			},
			expectedResult: dto.OrganizationConfig{},
			expectedError:  apperrors.OrganizationNotFound,
		},
		{
			name: "Error while updating organization config",
			context: context.WithValue(context.Background(), "userId", 1),
			setup: func(orgMock *mocks.OrganizationStorer) {
				orgMock.On("GetOrganizationConfig", mock.Anything).Return(repository.OrganizationConfig{
					ID:                           1,
					RewardMultiplier:             200,
					RewardQuotaRenewalFrequency: 12,
					Timezone:                    "ACT",
					CreatedAt:                   1719918501194,
					CreatedBy:                   7,
					UpdatedAt:                   1719920402224,
					UpdatedBy:                   7,
				}, nil).Once()
				orgMock.On("UpdateOrganizationCofig", mock.Anything, mock.Anything).Return(repository.OrganizationConfig{}, apperrors.InternalServer).Once()
			},
			organizationInput: dto.OrganizationConfig{
				ID:                           1,
				RewardMultiplier:             10,
				RewardQuotaRenewalFrequency: 5,
				Timezone:                    "UTC",
			},
			expectedResult: dto.OrganizationConfig{},
			expectedError:  apperrors.InternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(orgRepo)
			// Call the service method
			result, err := orgSvc.UpdateOrganizationConfig(tt.context, tt.organizationInput)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			// Assert that the mocks were called as expected
			orgRepo.AssertExpectations(t)
		})
	}
}
