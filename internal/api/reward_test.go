package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/app/appreciation/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGiveRewardHandlerHandler(t *testing.T) {
	appreciationSvc := new(mocks.Service)
	handler := createAppreciationHandler(appreciationSvc)

	tests := []struct {
		name               string
		input              dto.Reward
		mockSetup          func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "success",
			input: dto.Reward{Point: 1},
			mockSetup: func(mockSvc *mocks.Service) {
				mockSvc.On("GiveReward",mock.Anything,mock.Anything).Return(dto.Reward{Id: 1, AppreciationId: 1, SenderId: 1, Point: 1},nil)
			},
			expectedStatusCode: http.StatusOK,
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