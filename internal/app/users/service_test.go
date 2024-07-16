package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginUser(t *testing.T) {
	config.Load()
	userRepo := mocks.NewUserStorer(t)
	service := NewService(userRepo)

	tests := []struct {
		name            string
		context         context.Context
		u               dto.IntranetUserData
		setup           func(userMock *mocks.UserStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success for login for existing user",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{
					Id:                 1,
					EmployeeId:         "26",
					FirstName:          "sharyu",
					LastName:           "marwadi",
					Email:              "sharyu@josh.com",
					ProfileImgUrl:      "image url",
					RoleId:             1,
					RewardQuotaBalance: 10,
					Designation:        "Intern",
					GradeId:            1,
					Grade:              "J12",
					CreatedAt:          0,
				}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Success for register user",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("GetRoleByName", mock.Anything, mock.Anything).Return(1, nil).Once()
				userMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(dto.GetUserResp{
					Id:                 1,
					EmployeeId:         "26",
					FirstName:          "sharyu",
					LastName:           "marwadi",
					Email:              "sharyu@josh.com",
					ProfileImgUrl:      "image url",
					RoleId:             1,
					RewardQuotaBalance: 10,
					Designation:        "Intern",
					GradeId:            1,
					Grade:              "J12",
					CreatedAt:          0,
				}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Sync data success",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{
					Id:                 1,
					EmployeeId:         "26",
					FirstName:          "sharyu",
					LastName:           "marwadi",
					Email:              "sharyu@josh.com",
					ProfileImgUrl:      "image url",
					RoleId:             1,
					RewardQuotaBalance: 10,
					Designation:        "Manager",
					GradeId:            1,
					Grade:              "J12",
					CreatedAt:          0,
				}, nil).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("SyncData", mock.Anything, mock.Anything).Return(nil).Once()
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{
					Id:                 1,
					EmployeeId:         "26",
					FirstName:          "sharyu",
					LastName:           "marwadi",
					Email:              "sharyu@josh.com",
					ProfileImgUrl:      "image url",
					RoleId:             1,
					RewardQuotaBalance: 10,
					Designation:        "Intern",
					GradeId:            1,
					Grade:              "J12",
					CreatedAt:          0,
				}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:    "GetGradeByName Faliure",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{}, apperrors.GradeNotFound).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "GetRoleByName Faliure",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("GetRoleByName", mock.Anything, mock.Anything).Return(1, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Create user faliure",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("GetRoleByName", mock.Anything, mock.Anything).Return(1, nil).Once()
				userMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(dto.GetUserResp{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Sync data faliure",
			context: context.Background(),
			u: dto.IntranetUserData{
				Id:    1,
				Email: "sharyu@josh.com",
				PublicProfile: dto.PublicProfile{
					ProfileImgUrl: "image url",
					FirstName:     "sharyu",
					LastName:      "marwadi",
				},
				EmpolyeeDetail: dto.EmployeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{
					Id:                 1,
					EmployeeId:         "26",
					FirstName:          "sharyu",
					LastName:           "marwadi",
					Email:              "sharyu@josh.com",
					ProfileImgUrl:      "image url",
					RoleId:             1,
					RewardQuotaBalance: 10,
					Designation:        "Manager",
					GradeId:            1,
					Grade:              "J12",
					CreatedAt:          0,
				}, nil).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("SyncData", mock.Anything, mock.Anything).Return(apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userRepo)

			// test service
			_, err := service.LoginUser(test.context, test.u)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestGetUserList(t *testing.T) {
	userRepo := mocks.NewUserStorer(t)
	service := NewService(userRepo)

	tests := []struct {
		name            string
		context         context.Context
		reqData         dto.UserListReq
		setup           func(userMock *mocks.UserStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success for get user list",
			context: context.Background(),
			reqData: dto.UserListReq{},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserList", mock.Anything, mock.Anything).Return([]dto.GetUserListResp{}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Faliure for get user list",
			context: context.Background(),
			reqData: dto.UserListReq{},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserList", mock.Anything, mock.Anything).Return([]dto.GetUserListResp{}, apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userRepo)

			// test service
			_, err := service.GetUserList(test.context, test.reqData)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}

}

func TestUpdateRewardQuota(t *testing.T) {
	userRepo := mocks.NewUserStorer(t)
	service := NewService(userRepo)

	tests := []struct {
		name          string
		context       context.Context
		setup         func(userMock *mocks.UserStorer)
		expectedError error
	}{
		{
			name:    "success",
			context: context.Background(),
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("UpdateRewardQuota", mock.Anything, nil).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:    "failure",
			context: context.Background(),
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("UpdateRewardQuota", mock.Anything, nil).Return(apperrors.InternalServer)
			},
			expectedError: apperrors.InternalServer,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userRepo)

			// test service
			err := service.UpdateRewardQuota(test.context)

			assert.Equal(t, test.expectedError, err)

		})
	}
}

func TestGetActiveUserList(t *testing.T) {
	userRepo := mocks.NewUserStorer(t)
	service := NewService(userRepo)

	tests := []struct {
		name          string
		context       context.Context
		setup         func(userMock *mocks.UserStorer)
		expectedResp  []dto.ActiveUser
		expectedError error
	}{
		{
			name:    "success",
			context: context.Background(),
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetActiveUserList", mock.Anything, mock.Anything).Return([]repository.ActiveUser{
					{
						ID:                 55,
						FirstName:          "Deepak",
						LastName:           "Kumar",
						ProfileImageURL:    sql.NullString{String: "", Valid: false},
						BadgeName:          sql.NullString{String: "", Valid: false},
						AppreciationPoints: 0,
					},
					{
						ID:                 58,
						FirstName:          "Dominic",
						LastName:           "Lopes",
						ProfileImageURL:    sql.NullString{String: "", Valid: false},
						BadgeName:          sql.NullString{String: "Gold", Valid: true},
						AppreciationPoints: 5000,
					},
				}, nil).Once()
			},
			expectedResp: []dto.ActiveUser{
				{
					ID:                 55,
					FirstName:          "Deepak",
					LastName:           "Kumar",
					ProfileImageURL:    "",
					BadgeName:          "",
					AppreciationPoints: 0,
				},
				{
					ID:                 58,
					FirstName:          "Dominic",
					LastName:           "Lopes",
					ProfileImageURL:    "",
					BadgeName:          "Gold",
					AppreciationPoints: 5000,
				},
			},
			expectedError: nil,
		},
		{
			name:    "failure",
			context: context.Background(),
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetActiveUserList", mock.Anything, mock.Anything).Return([]repository.ActiveUser{}, apperrors.InternalServer).Once()
			},
			expectedResp:  []dto.ActiveUser{},
			expectedError: apperrors.InternalServer,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userRepo)

			// test service
			resp, err := service.GetActiveUserList(test.context)

			if err != nil {
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.Equal(t, test.expectedResp, resp)
			}

		})
	}
}

func TestGetUserById(t *testing.T) {
	userRepo := mocks.NewUserStorer(t)
	service := NewService(userRepo)

	tests := []struct {
		name            string
		context         context.Context
		userId          int64
		setup           func(userMock *mocks.UserStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success for get user list",
			context: context.Background(),
			userId:  2,
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserById", mock.Anything, mock.Anything).Return(dto.GetUserByIdResp{}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "GetUserById db function failed",
			context: context.Background(),
			userId:  0,
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserById", mock.Anything, mock.Anything).Return(dto.GetUserByIdResp{}, apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
		{
			name:    "Invalid Id",
			context: context.Background(),
			userId:  0,
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserById", mock.Anything, mock.Anything).Return(dto.GetUserByIdResp{}, apperrors.InvalidId).Once()

			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, constants.UserId, int64(test.userId))
			test.setup(userRepo)

			// test service
			_, err := service.GetUserById(ctx)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}

}

func TestGetTop10Users(t *testing.T) {
	userRepo := mocks.NewUserStorer(t)
	service := NewService(userRepo)

	tests := []struct {
		name            string
		context         context.Context
		setup           func(userMock *mocks.UserStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success for get top 10 users",
			context: context.Background(),
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetTop10Users", mock.Anything).Return([]repository.Top10Users{}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Faliure for get top 10 users",
			context: context.Background(),
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetTop10Users", mock.Anything).Return([]repository.Top10Users{}, apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.setup(userRepo)

			// test service
			_, err := service.GetTop10Users(test.context)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}

}
