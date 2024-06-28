package user

import (
	"context"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestLoginUser(t *testing.T) {
	config.Load("application")
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
				EmpolyeeDetail: dto.EmpolyeeDetail{
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
				EmpolyeeDetail: dto.EmpolyeeDetail{
					EmployeeId: "26",
					Designation: dto.Designation{
						Name: "Intern",
					},
					Grade: "J12",
				},
			},
			setup: func(userMock *mocks.UserStorer) {
				userMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(dto.GetUserResp{}, apperrors.UserNotFound).Once()
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(1, nil).Once()
				userMock.On("GetRewardOuotaDefault", mock.Anything).Return(10, nil).Once()
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
				EmpolyeeDetail: dto.EmpolyeeDetail{
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
				userMock.On("GetGradeByName", mock.Anything, mock.Anything).Return(1, nil).Once()
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
				}, apperrors.UserNotFound).Once()
			},
			isErrorExpected: false,
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
