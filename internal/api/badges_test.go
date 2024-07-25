package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/badges/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListBadgesHandler(t *testing.T) {
	badgeSvc := mocks.NewService(t)
	listBadgesHandler := listBadgesHandler(badgeSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for list badges",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListBadges", mock.Anything).Return([]dto.Badge{
					{
						ID:           1,
						Name:         "Gold",
						RewardPoints: 1000,
					},
					{
						ID:           2,
						Name:         "Platinum",
						RewardPoints: 2000,
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Error from badge service",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListBadges", mock.Anything).Return([]dto.Badge{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(badgeSvc)

			req, err := http.NewRequest(http.MethodGet, "/badges", nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(listBadgesHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetBadgeHandler(t *testing.T) {
	badgeSvc := mocks.NewService(t)
	getBadgeHandler := getBadgeHandler(badgeSvc)

	tests := []struct {
		name               string
		badgeID            string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:    "Success for get badge",
			badgeID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetBadge", mock.Anything, int8(1)).Return(dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "Error from badge service",
			badgeID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetBadge", mock.Anything, int8(1)).Return(dto.Badge{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(badgeSvc)

			req := httptest.NewRequest(http.MethodGet, "/badges/"+tt.badgeID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.badgeID})
			rr := httptest.NewRecorder()

			getBadgeHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			badgeSvc.AssertExpectations(t)
		})
	}
}

func TestDeleteBadgeHandler(t *testing.T) {
	badgeSvc := mocks.NewService(t)
	deleteBadgeHandler := deleteBadgeHandler(badgeSvc)

	tests := []struct {
		name               string
		badgeID            string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:    "Success for delete badge",
			badgeID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteBadge", mock.Anything, int8(1)).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "Error in deleting badge",
			badgeID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteBadge", mock.Anything, int8(1)).Return(apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(badgeSvc)

			req := httptest.NewRequest(http.MethodDelete, "/badges/"+tt.badgeID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.badgeID})
			rr := httptest.NewRecorder()

			deleteBadgeHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			badgeSvc.AssertExpectations(t)
		})
	}
}

func TestUpdateBadgeHandler(t *testing.T) {
	badgeSvc := mocks.NewService(t)
	updateBadgeHandler := updateBadgeHandler(badgeSvc)

	tests := []struct {
		name               string
		badgeID            string
		badgeUpdatedInfo   dto.Badge
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:    "Success for update badge",
			badgeID: "1",
			badgeUpdatedInfo: dto.Badge{
				Name:         "Gold",
				RewardPoints: 1000,
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateBadge", mock.Anything, dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:    "Error in updating badge",
			badgeID: "1",
			badgeUpdatedInfo: dto.Badge{
				Name:         "Gold",
				RewardPoints: 1000,
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateBadge", mock.Anything, dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(dto.Badge{}, apperrors.BadgeNotFound).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(badgeSvc)

			reqBody, _ := json.Marshal(tt.badgeUpdatedInfo)
			req := httptest.NewRequest(http.MethodPatch, "/badges", bytes.NewReader(reqBody))
			req = mux.SetURLVars(req, map[string]string{"id": tt.badgeID})
			rr := httptest.NewRecorder()

			updateBadgeHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			badgeSvc.AssertExpectations(t)
		})
	}
}

func TestCreateBadgeHandler(t *testing.T) {
	badgeSvc := mocks.NewService(t)
	createBadgeHandler := createBadgeHandler(badgeSvc)

	tests := []struct {
		name               string
		badgeID            string
		badgeUpdatedInfo   dto.Badge
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for create badge",
			badgeUpdatedInfo: dto.Badge{
				Name:         "Gold",
				RewardPoints: 1000,
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateBadge", mock.Anything, dto.Badge{
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Error in creating badge",
			badgeUpdatedInfo: dto.Badge{
				Name:         "Gold",
				RewardPoints: 1000,
			},
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateBadge", mock.Anything, dto.Badge{
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(dto.Badge{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(badgeSvc)

			reqBody, _ := json.Marshal(tt.badgeUpdatedInfo)
			req := httptest.NewRequest(http.MethodPost, "/badges", bytes.NewReader(reqBody))
			rr := httptest.NewRecorder()

			createBadgeHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			badgeSvc.AssertExpectations(t)
		})
	}
}
