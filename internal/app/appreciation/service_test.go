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
	service := NewService(appreciationRepo, corevalueRepo)

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
				apprMock.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(1)).Return(true, nil).Once()
				coreValueRepo.On("GetCoreValue", mock.Anything, int64(1)).Return(dto.GetCoreValueResp{}, nil).Once()
				apprMock.On("IsUserPresent", mock.Anything, tx, int64(2)).Return(true, nil).Once()
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
				apprMock.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(1)).Return(true, nil).Once()
				coreValueRepo.On("GetCoreValue", mock.Anything, int64(1)).Return(dto.GetCoreValueResp{}, apperrors.InvalidCoreValueData).Once()
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
				tx := &sql.Tx{}
				apprMock.On("BeginTx", mock.Anything).Return(tx, nil).Once()
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(1)).Return(true, nil).Once() // Mock sender presence check
				coreValueRepo.On("GetCoreValue", mock.Anything, int64(1)).Return(dto.GetCoreValueResp{ID: 1}, nil).Once()
				apprMock.On("IsUserPresent", mock.Anything, tx, int64(2)).Return(false, apperrors.UserNotFound).Once() // Ensure correct transaction context
				apprMock.On("HandleTransaction", mock.Anything, tx, false).Return(nil).Once()
			},
			isErrorExpected: true,
			expectedResult:  dto.Appreciation{},
			expectedError:   apperrors.UserNotFound,
		},

		{
			name:    "transaction failure",
			context: context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciation: dto.Appreciation{
				CoreValueID: 1,
				Receiver:    2,
			},
			setup: func(apprMock *mocks.AppreciationStorer, coreValueRepo *mocks.CoreValueStorer) {
				apprMock.On("IsUserPresent", mock.Anything, nil, int64(1)).Return(true, nil).Once()
				// apprMock.On("IsUserPresent", mock.Anything, nil, int64(1)).Return(true, nil).Once()
				apprMock.On("BeginTx", mock.Anything).Return(nil, apperrors.InternalServer).Once()
			},
			isErrorExpected: true,
			expectedResult:  dto.Appreciation{},
			expectedError:   apperrors.InternalServer,
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
	service := NewService(appreciationRepo, nil)

	tests := []struct {
		name            string
		context         context.Context
		appreciationId  int
		setup           func(apprMock *mocks.AppreciationStorer)
		isErrorExpected bool
		expectedResult  dto.ResponseAppreciation
		expectedError   error
	}{
		{
			name:           "successful get appreciation by id",
			context:        context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciationId: 1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, mock.Anything, 1).Return(repository.AppreciationInfo{
					ID:                  1,
					CoreValueName:       "Integrity",
					Description:         "Great work",
					IsValid:             true,
					TotalRewards:        100,
					Quarter:             "Q1",
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
			expectedResult: dto.ResponseAppreciation{
				ID:                  1,
				CoreValueName:       "Integrity",
				Description:         "Great work",
				TotalRewards:        100,
				Quarter:             "Q1",
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
				apprMock.On("GetAppreciationById", mock.Anything, mock.Anything, 1).Return(repository.AppreciationInfo{}, apperrors.AppreciationNotFound).Once()
			},
			isErrorExpected: true,
			expectedResult:  dto.ResponseAppreciation{},
			expectedError:   apperrors.AppreciationNotFound,
		},
		{
			name:           "database error",
			context:        context.WithValue(context.Background(), constants.UserId, int64(1)),
			appreciationId: 1,
			setup: func(apprMock *mocks.AppreciationStorer) {
				apprMock.On("GetAppreciationById", mock.Anything, mock.Anything, 1).Return(repository.AppreciationInfo{}, errors.New("database error"))
			},
			isErrorExpected: true,
			expectedResult:  dto.ResponseAppreciation{},
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
	service := NewService(appreciationRepo, nil)

	tests := []struct {
		name            string
		context         context.Context
		isValid         bool
		apprId          int
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
				apprMock.On("ValidateAppreciation", mock.Anything, nil, true, 1).Return(true, nil).Once()
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
				apprMock.On("ValidateAppreciation", mock.Anything, nil, false, 1).Return(false, apperrors.InternalServer).Once()
			},
			isErrorExpected: true,
			expectedResult:  false,
			expectedError:   apperrors.InternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(appreciationRepo)

			result, err := service.ValidateAppreciation(tt.context, tt.isValid, tt.apprId)

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
