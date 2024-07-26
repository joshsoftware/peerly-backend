package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/appreciation/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				mockSvc.On("CreateAppreciation", mock.Anything, dto.Appreciation{
					Description: "Great job!",
					CoreValueID: 5,
					Receiver:    2,
				}).Return(dto.Appreciation{
					ID:                1,
					Description:       "Great job!",
					CoreValueID:       5,
					TotalRewardPoints: 0,
					Quarter:           2,
					Receiver:          2,
					CreatedAt:         1721631405219,
					UpdatedAt:         1721631405219,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "invalid JSON input",
			input: dto.Appreciation{
				Description: "Great job!",
				CoreValueID: -1,
			},
			mockSetup:          func(mockSvc *mocks.Service) {},
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
				mockSvc.On("CreateAppreciation", mock.Anything, dto.Appreciation{
					CoreValueID: 5,
					Description: "Great job!",
					Receiver:    2,
				}).Return(dto.Appreciation{}, apperrors.InternalServer).Once()
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
	handler := getAppreciationByIDHandler(appreciationSvc)

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
				mockSvc.On("GetAppreciationById", mock.Anything, int32(1)).Return(dto.AppreciationResponse{
					ID:                  1,
					Description:         "Great job!",
					CoreValueName:       "Trust",
					CoreValueDesc:       "We foster trust by being transparent,reliable, and accountable in all our actions.",
					TotalRewardPoints:   0,
					Quarter:             2,
					SenderFirstName:     "John",
					SenderLastName:      "Doe",
					SenderImageURL:      "example.com",
					SenderDesignation:   "Software Engineer",
					ReceiverFirstName:   "Rohit",
					ReceiverLastName:    "Patil",
					ReceiverImageURL:    "example.com",
					ReceiverDesignation: "senior software engineer",
					TotalRewards:        5,
					GivenRewardPoint:    2,
					CreatedAt:           1721631405219,
					UpdatedAt:           1721631405219,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "service error",
			id:   "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciationById", mock.Anything, int32(1)).Return(dto.AppreciationResponse{}, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "appreciation not found",
			id:   "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetAppreciationById", mock.Anything, int32(1)).Return(dto.AppreciationResponse{}, apperrors.AppreciationNotFound).Once()
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

func TestListAppreciationsHandler(t *testing.T) {
	appreciationSvc := new(mocks.Service)
	handler := listAppreciationsHandler(appreciationSvc)

	tests := []struct {
		name               string
		queryParams        map[string]string
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "successful retrieval",
			queryParams: map[string]string{
				"name":       "John Doe",
				"sort_order": "asc",
				"page":       "1",
				"page_size":  "5",
			},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAppreciations", mock.Anything, dto.AppreciationFilter{
					Name:      "John Doe",
					SortOrder: "asc",
					Self:      false,
					Page:      1,
					Limit:     5,
				}).Return(dto.ListAppreciationsResponse{
					Appreciations: []dto.AppreciationResponse{
						{
							ID:                  1,
							Description:         "Great job!",
							CoreValueName:       "Trust",
							CoreValueDesc:       "We foster trust by being transparent,reliable, and accountable in all our actions.",
							TotalRewardPoints:   0,
							Quarter:             2,
							SenderFirstName:     "John",
							SenderLastName:      "Doe",
							SenderImageURL:      "example.com",
							SenderDesignation:   "Software Engineer",
							ReceiverFirstName:   "Rohit",
							ReceiverLastName:    "Patil",
							ReceiverImageURL:    "example.com",
							ReceiverDesignation: "senior software engineer",
							TotalRewards:        5,
							GivenRewardPoint:    2,
							CreatedAt:           1721631405219,
							UpdatedAt:           1721631405219,
						},
					},
					MetaData: dto.Pagination{
						CurrentPage:  1,
						TotalPage:    2,
						PageSize:     5,
						TotalRecords: 12,
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "service error",
			queryParams: map[string]string{
				"name":       "John Doe",
				"sort_order": "asc",
				"page":       "1",
				"page_size":  "10",
			},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListAppreciations", mock.Anything, dto.AppreciationFilter{
					Name:      "John Doe",
					SortOrder: "asc",
					Self:      false,
					Page:      1,
					Limit:     10,
				}).Return(dto.ListAppreciationsResponse{}, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
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
	handler := deleteAppreciationHandler(appreciationSvc)

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
				mockSvc.On("DeleteAppreciation", mock.Anything, int32(1)).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "invalid appreciation ID",
			appreciationID: "abcd",
			mockSetup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:           "service error",
			appreciationID: "1",
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteAppreciation", mock.Anything, int32(1)).Return(apperrors.InternalServer).Once()
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
