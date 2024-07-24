package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/pkg/testConfig"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestLoginUser(t *testing.T) {
	testConfig.Load()
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{
					Id:         1,
					EmployeeId: "26",
					FirstName:  "sharyu",
					LastName:   "marwadi",
					Email:      "sharyu@josh.com",
					ProfileImageURL: sql.NullString{
						Valid:  true,
						String: "image url",
					},
					RoleID:              1,
					RewardsQuotaBalance: 10,
					Designation:         "Intern",
					GradeId:             1,
					CreatedAt:           0,
				}, apperrors.InternalServerError).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("GetRewardMultiplier", mock.Anything).Return(int64(10), nil).Once()
				userMock.On("GetRoleByName", mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				userMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(repository.User{
					Id:         1,
					EmployeeId: "26",
					FirstName:  "sharyu",
					LastName:   "marwadi",
					Email:      "sharyu@josh.com",
					ProfileImageURL: sql.NullString{
						Valid:  true,
						String: "image url",
					},
					RoleID:              1,
					RewardsQuotaBalance: 10,
					Designation:         "Intern",
					GradeId:             1,
					CreatedAt:           0,
					Status:              1,
					SoftDelete:          false,
					SoftDeleteBy: sql.NullInt64{
						Valid: false,
						Int64: 0,
					},
					SoftDeleteOn: sql.NullTime{
						Valid: false,
						Time:  time.Now(),
					},
				}, nil).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{
					Id:         1,
					EmployeeId: "26",
					FirstName:  "sharyu",
					LastName:   "marwadi",
					Email:      "sharyu@josh.com",
					ProfileImageURL: sql.NullString{
						Valid:  true,
						String: "image url",
					},
					RoleID:              1,
					RewardsQuotaBalance: 10,
					Designation:         "Manager",
					GradeId:             1,
					CreatedAt:           0,
				}, nil).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("SyncData", mock.Anything, mock.Anything).Return(nil).Once()
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{
					Id:         1,
					EmployeeId: "26",
					FirstName:  "sharyu",
					LastName:   "marwadi",
					Email:      "sharyu@josh.com",
					ProfileImageURL: sql.NullString{
						Valid:  true,
						String: "image url",
					},
					RoleID:              1,
					RewardsQuotaBalance: 10,
					Designation:         "Intern",
					GradeId:             1,
					CreatedAt:           0,
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{}, apperrors.UserNotFound).Once()
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("GetRewardMultiplier", mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				userMock.On("GetRoleByName", mock.Anything, mock.Anything).Return(int64(1), apperrors.InternalServerError).Once()
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(repository.Grade{
					Id:     1,
					Name:   "J12",
					Points: 100,
				}, nil).Once()
				userMock.On("GetRewardMultiplier", mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				userMock.On("GetRoleByName", mock.Anything, mock.Anything).Return(int64(1), nil).Once()
				userMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(repository.User{}, apperrors.InternalServerError).Once()
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
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(repository.User{
					Id:         1,
					EmployeeId: "26",
					FirstName:  "sharyu",
					LastName:   "marwadi",
					Email:      "sharyu@josh.com",
					ProfileImageURL: sql.NullString{
						Valid:  true,
						String: "image url",
					},
					RoleID:              1,
					RewardsQuotaBalance: 10,
					Designation:         "Manager",
					GradeId:             1,
					CreatedAt:           0,
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
				userMock.On("GetTotalUserCount", mock.Anything, mock.Anything).Return(int64(280), nil).Once()
				userMock.On("ListUsers", mock.Anything, mock.Anything).Return([]repository.User{}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Faliure for get user count",
			context: context.Background(),
			reqData: dto.UserListReq{},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetTotalUserCount", mock.Anything, mock.Anything).Return(int64(0), apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Faliure for get user count",
			context: context.Background(),
			reqData: dto.UserListReq{},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetTotalUserCount", mock.Anything, mock.Anything).Return(int64(280), nil).Once()
				userMock.On("ListUsers", mock.Anything, mock.Anything).Return([]repository.User{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userRepo)

			// test service
			_, err := service.ListUsers(test.context, test.reqData)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}

}
