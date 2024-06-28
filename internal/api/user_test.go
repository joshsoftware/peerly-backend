package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/app/users/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

func TestLoginUser(t *testing.T) {
	userSvc := mocks.NewService(t)
	listCoreValuesHandler := loginUser(userSvc)

	tests := []struct {
		name               string
		authToken          string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:      "Success for login",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{
						JwtToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.qdKwgFLwHeNg8PaYFEjLT7g4sk0DGdoSHt-wZ7eq5LQ",
						UserId:   6,
					},
				}, nil).Once()
				mockSvc.On("GetIntranetUserData", mock.Anything, mock.Anything).Return(dto.IntranetUserData{
					Id:    1,
					Email: "sharyu@josh.com",
					PublicProfile: dto.PublicProfile{
						ProfileImgUrl: "image url",
						FirstName:     "Sharyu",
						LastName:      "Marwadi",
					},
					EmpolyeeDetail: dto.EmpolyeeDetail{
						EmployeeId: "26",
						Designation: dto.Designation{
							Name: "Manager",
						},
						Grade: "J2",
					},
				}, nil).Once()
				mockSvc.On("LoginUser", mock.Anything, mock.Anything).Return(dto.LoginUserResp{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid auth token",
			authToken:          "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:      "Validation api faliure",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{},
				}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:      "Validation api faliure",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{},
				}, apperrors.JSONParsingErrorResp).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:      "Validation api faliure",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{},
				}, apperrors.IntranetValidationFailed).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Intranet get user api faliure",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{
						JwtToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.qdKwgFLwHeNg8PaYFEjLT7g4sk0DGdoSHt-wZ7eq5LQ",
						UserId:   6,
					},
				}, nil).Once()
				mockSvc.On("GetIntranetUserData", mock.Anything, mock.Anything).Return(dto.IntranetUserData{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:      "Intranet get user api faliure",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{
						JwtToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.qdKwgFLwHeNg8PaYFEjLT7g4sk0DGdoSHt-wZ7eq5LQ",
						UserId:   6,
					},
				}, nil).Once()
				mockSvc.On("GetIntranetUserData", mock.Anything, mock.Anything).Return(dto.IntranetUserData{}, apperrors.JSONParsingErrorResp).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:      "Faliure for login",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{
						JwtToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.qdKwgFLwHeNg8PaYFEjLT7g4sk0DGdoSHt-wZ7eq5LQ",
						UserId:   6,
					},
				}, nil).Once()
				mockSvc.On("GetIntranetUserData", mock.Anything, mock.Anything).Return(dto.IntranetUserData{
					Id:    1,
					Email: "sharyu@josh.com",
					PublicProfile: dto.PublicProfile{
						ProfileImgUrl: "image url",
						FirstName:     "Sharyu",
						LastName:      "Marwadi",
					},
					EmpolyeeDetail: dto.EmpolyeeDetail{
						EmployeeId: "26",
						Designation: dto.Designation{
							Name: "Manager",
						},
						Grade: "J2",
					},
				}, nil).Once()
				mockSvc.On("LoginUser", mock.Anything, mock.Anything).Return(dto.LoginUserResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("POST", "/user/login", bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}
			req.Header.Set("Authorization", test.authToken)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(listCoreValuesHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}