package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/peerly-backend/internal/app/reward/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGiveRewardHandlerHandler(t *testing.T) {
	rewardSvc := new(mocks.Service)
	handler := giveRewardHandler(rewardSvc)

	tests := []struct {
		name               string
		id                 string
		input              dto.Reward
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:  "success",
			id:    "1",
			input: dto.Reward{Point: 2},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GiveReward", mock.Anything, mock.Anything).Return(dto.Reward{Id: 1, AppreciationId: 1, SenderId: 1, Point: 1}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Error decoding request data",
			id:   "1",
			mockSetup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Invalid reward point",
			id:    "1",
			input: dto.Reward{Point: 10},
			mockSetup: func(mockSvc *mocks.Service) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "give reward failure",
			id:    "1",
			input: dto.Reward{Point: 2},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GiveReward", mock.Anything, mock.Anything).Return(dto.Reward{}, apperrors.InternalServer)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(rewardSvc)

			reqBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodGet, "/reward/"+tt.id, bytes.NewReader(reqBody))
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			rewardSvc.AssertExpectations(t)
		})
	}
}
