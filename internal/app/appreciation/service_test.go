package appreciation

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/constants"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAppreciation(t *testing.T) {
	appreciationRepo := mocks.NewAppreciationStorer(t)
	corevalueRepo := mocks.NewCoreValueStorer(t)
	userRepo := mocks.NewUserStorer(t)
	service := NewService(appreciationRepo, corevalueRepo, userRepo)

	tests := []struct {
		name            string
		context         context.Context
		appreciation    dto.Appreciation
		setup           func(apprMock *mocks.AppreciationStorer, coreValueRepo *mocks.CoreValueStorer)
		isErrorExpected bool
		expectedResult  dto.Appreciation
		expectedError   error
	}{
		{
			name:    "successful appreciation creation",
			context: context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciation: dto.Appreciation{
				CoreValueID: 1,
				Receiver:    2,
			},
			setup: func(apprMock *mocks.AppreciationStorer, coreValueRepo *mocks.CoreValueStorer) {
				tx := &sql.Tx{}
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(2)).Return(true, nil).Once()
				apprMock.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				coreValueRepo.On("GetCoreValue", mock.Anything, int64(1)).Return(repository.CoreValue{
					ID:                1,
					Name:              "Trust",
					Description:       "We foster trust by being transparent,reliable, and accountable in all our actions",
					ParentCoreValueID: sql.NullInt64{Int64: int64(0), Valid: true},
				}, nil).Once()
				apprMock.On("CreateAppreciation", mock.Anything, tx, mock.Anything).Return(repository.Appreciation{ID: 1}, nil).Once()
				apprMock.On("HandleTransaction", mock.Anything, tx, true).Return(nil).Once()
			},
			isErrorExpected: false,
			expectedResult:  dto.Appreciation{ID: 1},
			expectedError:   nil,
		},
		{
			name:    "core value not found",
			context: context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciation: dto.Appreciation{
				CoreValueID: 1,
				Receiver:    2,
			},
			setup: func(apprMock *mocks.AppreciationStorer, coreValueRepo *mocks.CoreValueStorer) {
				tx := &sql.Tx{}
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(2)).Return(true, nil).Once()
				apprMock.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				coreValueRepo.On("GetCoreValue", mock.Anything, int64(1)).Return(repository.CoreValue{}, apperrors.InvalidCoreValueData).Once()
				apprMock.On("HandleTransaction", mock.Anything, tx, false).Return(apperrors.InvalidCoreValueData).Once()
			},
			isErrorExpected: true,
			expectedResult:  dto.Appreciation{},
			expectedError:   apperrors.InvalidCoreValueData,
		},
		{
			name:    "receiver not found",
			context: context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciation: dto.Appreciation{
				CoreValueID: 1,
				Description: "Great teamwork!",
				Receiver:    2,
			},
			setup: func(apprMock *mocks.AppreciationStorer, coreValueRepo *mocks.CoreValueStorer) {
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(2)).Return(false, apperrors.UserNotFound).Once() // Ensure correct transaction context
			},
			isErrorExpected: true,
			expectedResult:  dto.Appreciation{},
			expectedError:   apperrors.UserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(appreciationRepo, corevalueRepo)

			result, err := service.CreateAppreciation(tt.context, tt.appreciation)

			if tt.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			appreciationRepo.AssertExpectations(t)
			corevalueRepo.AssertExpectations(t)
		})
	}
}

func TestGetAppreciationById(t *testing.T) {
	appreciationRepo := mocks.NewAppreciationStorer(t)
	coreVaueRepo := mocks.NewCoreValueStorer(t)
	userRepo := mocks.NewUserStorer(t)
	service := NewService(appreciationRepo, coreVaueRepo, userRepo)

	tests := []struct {
		name            string
		context         context.Context
		appreciationId  int32
		setup           func(apprMock *mocks.AppreciationStorer)
		isErrorExpected bool
		expectedResult  dto.AppreciationResponse
		expectedError   error
	}{
		{
			name:           "successful get appreciation by id",
			context:        context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciationId: 1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, mock.Anything, int32(1)).Return(repository.AppreciationResponse{
					ID:                  1,
					CoreValueName:       "Integrity",
					Description:         "Great work",
					IsValid:             true,
					TotalRewards:        100,
					Quarter:             2,
					SenderFirstName:     "John",
					SenderLastName:      "Doe",
					SenderImageURL:      sql.NullString{String: "image_url", Valid: true},
					SenderDesignation:   "Manager",
					ReceiverFirstName:   "Jane",
					ReceiverLastName:    "Smith",
					ReceiverImageURL:    sql.NullString{String: "image_url", Valid: true},
					ReceiverDesignation: "Developer",
					CreatedAt:           1620000000,
					UpdatedAt:           1620000000,
				}, nil).Once()
			},
			isErrorExpected: false,
			expectedResult: dto.AppreciationResponse{
				ID:                  1,
				CoreValueName:       "Integrity",
				Description:         "Great work",
				TotalRewards:        100,
				Quarter:             2,
				SenderFirstName:     "John",
				SenderLastName:      "Doe",
				SenderImageURL:      "image_url",
				SenderDesignation:   "Manager",
				ReceiverFirstName:   "Jane",
				ReceiverLastName:    "Smith",
				ReceiverImageURL:    "image_url",
				ReceiverDesignation: "Developer",
				CreatedAt:           1620000000,
				UpdatedAt:           1620000000,
			},
			expectedError: nil,
		},
		{
			name:           "appreciation not found",
			context:        context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciationId: 1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, mock.Anything, int32(1)).Return(repository.AppreciationResponse{}, apperrors.AppreciationNotFound).Once()
			},
			isErrorExpected: true,
			expectedResult:  dto.AppreciationResponse{},
			expectedError:   apperrors.AppreciationNotFound,
		},
		{
			name:           "database error",
			context:        context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciationId: 1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, mock.Anything, int32(1)).Return(repository.AppreciationResponse{}, errors.New("database error"))
			},
			isErrorExpected: true,
			expectedResult:  dto.AppreciationResponse{},
			expectedError:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(appreciationRepo)

			result, err := service.GetAppreciationById(tt.context, tt.appreciationId)

			if tt.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			appreciationRepo.AssertExpectations(t)
		})
	}
}

func TestValidateAppreciation(t *testing.T) {
	appreciationRepo := mocks.NewAppreciationStorer(t)
	coreVaueRepo := mocks.NewCoreValueStorer(t)
	userRepo := mocks.NewUserStorer(t)
	service := NewService(appreciationRepo, coreVaueRepo, userRepo)

	tests := []struct {
		name            string
		context         context.Context
		isValid         bool
		apprId          int32
		setup           func(apprMock *mocks.AppreciationStorer)
		isErrorExpected bool
		expectedResult  bool
		expectedError   error
	}{
		{
			name:    "successful validation",
			context: context.Background(),
			isValid: true,
			apprId:  1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("DeleteAppreciation", mock.Anything, nil, int32(1)).Return(nil).Once()
			},
			isErrorExpected: false,
			expectedResult:  true,
			expectedError:   nil,
		},
		{
			name:    "validation failed",
			context: context.Background(),
			isValid: false,
			apprId:  1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("DeleteAppreciation", mock.Anything, nil, int32(1)).Return(apperrors.InternalServer).Once()
			},
			isErrorExpected: true,
			expectedResult:  false,
			expectedError:   apperrors.InternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(appreciationRepo)

			err := service.DeleteAppreciation(tt.context, tt.apprId)

			if tt.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			}

			appreciationRepo.AssertExpectations(t)
		})
	}
}
