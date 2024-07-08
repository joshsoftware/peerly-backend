package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/app/appreciation/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/gorilla/mux"
)

func TestCreateAppreciationHandler(t *testing.T) {
	appreciationSvc := new(mocks.Service)
	handler := createAppreciationHandler(appreciationSvc)

	tests := []struct {
		name               string
		input              dto.Appreciation
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "successful creation",
			input: dto.Appreciation{
				Description: "Great job!",
				CoreValueID: 5,
				Receiver:    2,
			},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAppreciation", mock.Anything, mock.Anything).Return(dto.Appreciation{}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "invalid JSON input",
			input: dto.Appreciation{
				Description: "Great job!",
				CoreValueID: -1, 
			},
			mockSetup: func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			input: dto.Appreciation{
				CoreValueID: 5,
				Description: "Great job!",
				Receiver:    2,
			},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAppreciation", mock.Anything, mock.Anything).Return(dto.Appreciation{}, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(appreciationSvc)

			reqBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/appreciation", bytes.NewReader(reqBody))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			appreciationSvc.AssertExpectations(t)
		})
	}
}


func TestGetAppreciationByIdHandler(t *testing.T) {
	appreciationSvc := new(mocks.Service)
	handler := getAppreciationByIdHandler(appreciationSvc)

	tests := []struct {
		name               string
		id                 string
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "successful retrieval",
			id:   "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciationById", mock.Anything, 1).Return(dto.ResponseAppreciation{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "service error",
			id:   "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciationById", mock.Anything, 1).Return(dto.ResponseAppreciation{}, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "appreciation not found",
			id:   "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciationById", mock.Anything, 1).Return(dto.ResponseAppreciation{}, apperrors.AppreciationNotFound).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(appreciationSvc)

			req := httptest.NewRequest(http.MethodGet, "/appreciation/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			appreciationSvc.AssertExpectations(t)
		})
	}
}

func TestGetAppreciationsHandler(t *testing.T) {
	appreciationSvc := new(mocks.Service)
	handler := getAppreciationsHandler(appreciationSvc)

	tests := []struct {
		name               string
		queryParams        map[string]string
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "successful retrieval",
			queryParams: map[string]string{
				"name":      "John Doe",
				"sort_order": "asc",
				"page": "1",
				"limit": "5",
			},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciation", mock.Anything, dto.AppreciationFilter{
					Name:      "John Doe",
					SortOrder: "asc",
					Page:      1,
					Limit:     5,
				}).Return(dto.GetAppreciationResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "service error",
			queryParams:        map[string]string{},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciation", mock.Anything, dto.AppreciationFilter{
					Name:      "",
					SortOrder: "",
					Page:      1,
					Limit:     10, // Default limit in case not provided
				}).Return(dto.GetAppreciationResponse{}, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "invalid pagination parameters",
			queryParams: map[string]string{
				"page":  "invalid",
				"limit": "invalid",
			},
			mockSetup: func(mockSvc *mocks.Service) {
				// No service call expected
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(appreciationSvc)

			req := httptest.NewRequest(http.MethodGet, "/appreciations", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			appreciationSvc.AssertExpectations(t)
		})
	}
}


func TestValidateAppreciationHandler(t *testing.T) {
	appreciationSvc := new(mocks.Service)
	handler := validateAppreciationHandler(appreciationSvc)

	tests := []struct {
		name               string
		appreciationID     string
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:           "successful validation",
			appreciationID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidateAppreciation", mock.Anything, false, 1).Return(true, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "invalid appreciation ID",
			appreciationID: "invalid",
			mockSetup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:           "service error",
			appreciationID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidateAppreciation", mock.Anything, false, 1).Return(false, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "validation failed",
			appreciationID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidateAppreciation", mock.Anything, false, 1).Return(false, nil).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(appreciationSvc)

			req := httptest.NewRequest(http.MethodPost, "/appreciations/validate/"+tt.appreciationID, nil)
			rr := httptest.NewRecorder()

			// Use mux router to test path variables
			router := mux.NewRouter()
			router.HandleFunc("/appreciations/validate/{id}", handler)
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			appreciationSvc.AssertExpectations(t)
		})
	}
}