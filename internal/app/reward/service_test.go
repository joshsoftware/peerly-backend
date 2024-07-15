package reward

import (
	"context"
	// "database/sql"
	// "errors"
	"testing"

	// "github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestGiveReward(t *testing.T) {
	tests := []struct {
		name            string
		ctx             context.Context
		rewardReq       dto.Reward
		setup           func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer)
		isErrorExpected bool
		expectedResult  dto.Reward
		expectedError   error
	}{
		{
			name: "Success",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 2, ReceiverId: 3}, nil)
				rwrdMock.On("UserHasRewardQuota", mock.Anything, nil, int64(1), int64(10)).Return(true, nil)
				rwrdMock.On("IsUserRewardForAppreciationPresent", mock.Anything, nil, int64(1), int64(1)).Return(false, nil)
				rwrdMock.On("BeginTx", mock.Anything).Return(nil, nil)
				rwrdMock.On("GiveReward", mock.Anything, mock.Anything, mock.Anything).Return(repository.Reward{Id: 1, AppreciationId: 1, SenderId: 1, Point: 10}, nil)
				rwrdMock.On("DeduceRewardQuotaOfUser", mock.Anything, mock.Anything, int64(1), 10).Return(true, nil)
				apprMock.On("HandleTransaction", mock.Anything, mock.Anything, true).Return(nil)
			},
			isErrorExpected: false,
			expectedResult:  dto.Reward{Id: 1, AppreciationId: 1, SenderId: 1, Point: 10},
			expectedError:   nil,
		},
		{
			name: "Error in parsing userid from token",
			ctx:  context.Background(),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup:           func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.InternalServer,
		},
		{
			name: "Self appreciation reward error",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 1, ReceiverId: 2}, nil)
			},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.SelfAppreciationRewardError,
		},
		{
			name: "Self reward error",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 2, ReceiverId: 1}, nil)
			},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.SelfRewardError,
		},
		{
			name: "Insufficient reward quota",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 2, ReceiverId: 3}, nil)
				rwrdMock.On("UserHasRewardQuota", mock.Anything, nil, int64(1), int64(10)).Return(false, nil)
			},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.RewardQuotaIsNotSufficient,
		},
		{
			name: "Reward already present",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 2, ReceiverId: 3}, nil)
				rwrdMock.On("UserHasRewardQuota", mock.Anything, nil, int64(1), int64(10)).Return(true, nil)
				rwrdMock.On("IsUserRewardForAppreciationPresent", mock.Anything, nil, int64(1), int64(1)).Return(true, nil)
			},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.RewardAlreadyPresent,
		},
		{
			name: "Database transaction error",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 2, ReceiverId: 3}, nil)
				rwrdMock.On("UserHasRewardQuota", mock.Anything, nil, int64(1), int64(10)).Return(true, nil)
				rwrdMock.On("IsUserRewardForAppreciationPresent", mock.Anything, nil, int64(1), int64(1)).Return(false, nil)
				rwrdMock.On("BeginTx", mock.Anything).Return(nil, apperrors.InternalServer)
			},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.InternalServer,
		},
		{
			name: "Deduce reward quota failure",
			ctx:  context.WithValue(context.Background(), constants.UserId, int64(1)),
			rewardReq: dto.Reward{
				AppreciationId: 1,
				Point:          10,
			},
			setup: func(rwrdMock *mocks.RewardStorer, apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, nil, 1).Return(repository.AppreciationInfo{SenderId: 2, ReceiverId: 3}, nil)
				rwrdMock.On("UserHasRewardQuota", mock.Anything, nil, int64(1), int64(10)).Return(true, nil)
				rwrdMock.On("IsUserRewardForAppreciationPresent", mock.Anything, nil, int64(1), int64(1)).Return(false, nil)
				rwrdMock.On("BeginTx", mock.Anything).Return(nil, nil)
				rwrdMock.On("GiveReward", mock.Anything, mock.Anything, mock.Anything).Return(repository.Reward{Id: 1, AppreciationId: 1, SenderId: 1, Point: 10}, nil)
				rwrdMock.On("DeduceRewardQuotaOfUser", mock.Anything, mock.Anything, int64(1), 10).Return(false, apperrors.RewardQuotaIsNotSufficient)
				apprMock.On("HandleTransaction", mock.Anything, mock.Anything, false).Return(apperrors.RewardQuotaIsNotSufficient)
			},
			isErrorExpected: true,
			expectedResult:  dto.Reward{},
			expectedError:   apperrors.RewardQuotaIsNotSufficient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rwrdMock := &mocks.RewardStorer{}
			apprMock := &mocks.AppreciationStorer{}

			if test.setup != nil {
				test.setup(rwrdMock, apprMock)
			}

			service := &service{
				rewardRepo:        rwrdMock,
				appreciationRepo:  apprMock,
			}

			result, err := service.GiveReward(test.ctx, test.rewardReq)

			if test.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}

			rwrdMock.AssertExpectations(t)
			apprMock.AssertExpectations(t)
		})
	}
}
