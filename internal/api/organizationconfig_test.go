package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/app/organizationConfig/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

func TestGetOrganizationConfigHandler(t *testing.T) {
	orgSvc := mocks.NewService(t)
	getOrganizationConfigHandler := getOrganizationConfigHandler(orgSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success fetching organization",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetOrganizationConfig", mock.Anything).Return(dto.OrganizationConfig{
					ID: 1,
					RewardMultiplier: 1,
					RewardQuotaRenewalFrequency: 1,
					Timezone: "UTC",
					CreatedAt:1,
					CreatedBy:1721631405219,
					UpdatedAt:1,
					UpdatedBy:1721631405219,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Error fetching organization",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetOrganizationConfig", mock.Anything).Return(dto.OrganizationConfig{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)

			req, err := http.NewRequest(http.MethodGet, "/organization/config", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			getOrganizationConfigHandler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}


func TestCreateOrganizationConfigHandler(t *testing.T) {
	orgSvc := mocks.NewService(t)
	handler := createOrganizationConfigHandler(orgSvc)

	tests := []struct {
		name               string
		requestBody        dto.OrganizationConfig
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Successful organization config creation",
			requestBody: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateOrganizationConfig", mock.Anything, mock.Anything).Return(dto.OrganizationConfig{
					ID: 1,
					RewardMultiplier: 200,
					RewardQuotaRenewalFrequency: 12,
					Timezone: "ACT",
					CreatedAt: 1719918501194,
					CreatedBy: 7,
					UpdatedAt: 1719920402224,
					UpdatedBy: 7,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "JSON decoding error",
			requestBody: dto.OrganizationConfig{
			},
			setup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Validation error",
			requestBody: dto.OrganizationConfig{
				RewardMultiplier: 0,
				RewardQuotaRenewalFrequency: 0,
				Timezone: "ABCD",
			},
			setup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error creating organization config",
			requestBody: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateOrganizationConfig", mock.Anything, mock.Anything).Return(dto.OrganizationConfig{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)
			reqBody, err := json.Marshal(test.requestBody)
			if err != nil {
				t.Fatal("Failed to marshal request body")
			}

			req, err := http.NewRequest(http.MethodPost, "/create_organization", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}


func TestUpdateOrganizationConfigHandler(t *testing.T) {
	orgSvc := mocks.NewService(t)
	handler := updateOrganizationConfigHandler(orgSvc)

	tests := []struct {
		name               string
		requestBody        dto.OrganizationConfig
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Successful organization config update",
			requestBody: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateOrganizationConfig", mock.Anything, mock.Anything).Return(dto.OrganizationConfig{
					ID: 1,
					RewardMultiplier: 10,
					RewardQuotaRenewalFrequency: 5,
					Timezone: "UTC",
					CreatedAt: 1719918501194,
					CreatedBy: 7,
					UpdatedAt: 1719920402224,
					UpdatedBy: 7,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Validation error",
			requestBody: dto.OrganizationConfig{
				RewardMultiplier: 0,
				RewardQuotaRenewalFrequency: 0,
				Timezone: "ABCD",
			},
			setup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error updating organization config",
			requestBody: dto.OrganizationConfig{
				RewardMultiplier: 10,
				RewardQuotaRenewalFrequency: 5,
				Timezone: "UTC",
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateOrganizationConfig", mock.Anything, mock.Anything).Return(dto.OrganizationConfig{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(orgSvc)
			reqBody, err := json.Marshal(test.requestBody)
			if err != nil {
				t.Fatal("Failed to marshal request body")
			}

			req, err := http.NewRequest(http.MethodPut, "/update_organization", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
