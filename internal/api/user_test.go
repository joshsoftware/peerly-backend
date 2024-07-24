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
	loginUserHandler := loginUser(userSvc)

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
					EmpolyeeDetail: dto.EmployeeDetail{
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
			expectedStatusCode: http.StatusUnauthorized,
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
					EmpolyeeDetail: dto.EmployeeDetail{
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
		{
			name:      "Error for invalid first name",
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
						LastName:      "Marwadi",
					},
					EmpolyeeDetail: dto.EmployeeDetail{
						EmployeeId: "26",
						Designation: dto.Designation{
							Name: "Manager",
						},
						Grade: "J2",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Error for invalid last name",
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
					},
					EmpolyeeDetail: dto.EmployeeDetail{
						EmployeeId: "26",
						Designation: dto.Designation{
							Name: "Manager",
						},
						Grade: "J2",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Error for invalid designation",
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
					EmpolyeeDetail: dto.EmployeeDetail{
						EmployeeId: "26",
						Grade:      "J2",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Error for invalid email",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ValidatePeerly", mock.Anything, mock.Anything).Return(dto.ValidateResp{
					Data: dto.IntranetValidateApiData{
						JwtToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.qdKwgFLwHeNg8PaYFEjLT7g4sk0DGdoSHt-wZ7eq5LQ",
						UserId:   6,
					},
				}, nil).Once()
				mockSvc.On("GetIntranetUserData", mock.Anything, mock.Anything).Return(dto.IntranetUserData{
					Id: 1,
					PublicProfile: dto.PublicProfile{
						ProfileImgUrl: "image url",
						FirstName:     "Sharyu",
						LastName:      "Marwadi",
					},
					EmpolyeeDetail: dto.EmployeeDetail{
						EmployeeId: "26",
						Designation: dto.Designation{
							Name: "Manager",
						},
						Grade: "J2",
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Error for invalid grade",
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
					EmpolyeeDetail: dto.EmployeeDetail{
						EmployeeId: "26",
						Designation: dto.Designation{
							Name: "Manager",
						},
					},
				}, nil).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}
			req.Header.Set("Intranet-Auth", test.authToken)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(loginUserHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetUserHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	getUsersHandler := getUsersHandler(userSvc)

	tests := []struct {
		name               string
		authToken          string
		page               string
		per_page           string
		paramName          string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:      "Success for get user list",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			page:      "1",
			per_page:  "10",
			paramName: "sharyu%20marwadi",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetUserList", mock.Anything, mock.Anything).Return(dto.UserListWithMetadata{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "Internal server error",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			page:      "1",
			per_page:  "10",
			paramName: "sharyu%20marwadi",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetUserList", mock.Anything, mock.Anything).Return(dto.UserListWithMetadata{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "Page param not found",
			authToken:          "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			page:               "",
			per_page:           "10",
			paramName:          "sharyu%20marwadi",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("GET", fmt.Sprintf("/users/all?name=%s&page=%s&per_page=%s", test.name, test.page, test.per_page), bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}
			if test.page == "" {
				req, err = http.NewRequest("GET", fmt.Sprintf("/users/all?name=%s&per_page=%s", test.name, test.per_page), bytes.NewBuffer([]byte("")))
				if err != nil {
					t.Fatal(err)
					return
				}
			}

			req.Header.Set("Authorization", test.authToken)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getUsersHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetActiveUserListHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	getActiveUserListHandler := getActiveUserListHandler(userSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for get active user list",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetActiveUserList", mock.Anything, mock.Anything).Return([]dto.ActiveUser{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Failure",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetActiveUserList", mock.Anything, mock.Anything).Return([]dto.ActiveUser{}, apperrors.InternalServer).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("GET", "/users/activeusers", bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getActiveUserListHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetUserByIdHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	getUserById := getUserByIdHandler(userSvc)

	tests := []struct {
		name               string
		authToken          string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:      "Success for get user by id",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetUserById", mock.Anything).Return(dto.GetUserByIdResp{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "Faliure for get user by id",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetUserById", mock.Anything).Return(dto.GetUserByIdResp{}, apperrors.InvalidId).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Faliure for get user by id",
			authToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2fQ.XaYo0qdBCdDh1-nEeuUSdTbtp0enWFIySKnw-oQpTBg",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetUserById", mock.Anything).Return(dto.GetUserByIdResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("GET", "/user_profile", bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}

			req.Header.Set("Authorization", test.authToken)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getUserById)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetTop10UsersHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	getTop10Users := getTop10UserHandler(userSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for get top 10 users",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetTop10Users", mock.Anything).Return([]dto.Top10User{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Faliure for get top 10 users",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetTop10Users", mock.Anything).Return([]dto.Top10User{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("GET", "/users/top10", bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getTop10Users)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
