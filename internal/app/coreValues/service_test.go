package corevalues

import (
	"context"

	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestListCoreValues(t *testing.T) {
	coreValueRepo := mocks.NewCoreValueStorer(t)
	service := NewService(coreValueRepo)

	tests := []struct {
		name            string
		context         context.Context
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success for list corevalues",
			context: context.Background(),
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("ListCoreValues", mock.Anything).Return([]dto.CoreValue{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:    "Error in list corevalues",
			context: context.Background(),
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("ListCoreValues", mock.Anything).Return([]dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.ListCoreValues(test.context)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestGetCoreValue(t *testing.T) {
	coreValueRepo := mocks.NewCoreValueStorer(t)
	service := NewService(coreValueRepo)

	tests := []struct {
		name            string
		context         context.Context
		coreValueId     string
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:        "Success for get corevalue",
			context:     context.Background(),
			coreValueId: "1",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:        "Invalid corevalue",
			context:     context.Background(),
			coreValueId: "0",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InvalidCoreValueData).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.GetCoreValue(test.context, test.coreValueId)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateCoreValue(t *testing.T) {
	coreValueRepo := mocks.NewCoreValueStorer(t)
	service := NewService(coreValueRepo)

	tests := []struct {
		name            string
		context         context.Context
		userId          int64
		coreValue       dto.CreateCoreValueReq
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success for create corevalue",
			context: context.Background(),
			userId:  1,
			coreValue: dto.CreateCoreValueReq{
				Name:              "CoreValue",
				Description:       "core value desc",
				ParentCoreValueID: nil,
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything).Return(true, nil).Once()
				coreValueMock.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:    "Repeated core value",
			context: context.Background(),
			userId:  1,
			coreValue: dto.CreateCoreValueReq{
				Name:              "CoreValue",
				Description:       "core value desc",
				ParentCoreValueID: nil,
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything).Return(false, nil).Once()

			},
			isErrorExpected: true,
		},
		{
			name:    "Error while creating core value",
			context: context.Background(),
			userId:  1,
			coreValue: dto.CreateCoreValueReq{
				Name:              "CoreValue",
				Description:       "core value desc",
				ParentCoreValueID: nil,
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything).Return(true, nil).Once()
				coreValueMock.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.CreateCoreValue(test.context, test.coreValue)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestUpdateCoreValue(t *testing.T) {
	coreValueRepo := mocks.NewCoreValueStorer(t)
	service := NewService(coreValueRepo)

	tests := []struct {
		name            string
		context         context.Context
		coreValueId     string
		reqData         dto.UpdateQueryRequest
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:        "Success for update corevalue",
			context:     context.Background(),
			coreValueId: "1",
			reqData: dto.UpdateQueryRequest{
				Name:        "Updated core value",
				Description: "updated description",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything).Return(true, nil).Once()
				coreValueMock.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:        "Invalid corevalue id",
			context:     context.Background(),
			coreValueId: "1",
			reqData: dto.UpdateQueryRequest{
				Name:        "Updated core value",
				Description: "updated description",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InvalidCoreValueData).Once()
			},
			isErrorExpected: true,
		},
		{
			name:        "Error in updating core value",
			context:     context.Background(),
			coreValueId: "1",
			reqData: dto.UpdateQueryRequest{
				Name:        "Updated core value",
				Description: "updated description",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything).Return(true, nil).Once()
				coreValueMock.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.UpdateCoreValue(test.context, test.coreValueId, test.reqData)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}
