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
		organizationId  string
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:           "Success for list corevalues",
			context:        context.Background(),
			organizationId: "1",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("ListCoreValues", mock.Anything, mock.Anything).Return([]dto.ListCoreValuesResp{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:           "Invalid organisation id",
			context:        context.Background(),
			organizationId: "0",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(apperrors.InvalidOrgId).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Error in list corevalues",
			context:        context.Background(),
			organizationId: "1",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("ListCoreValues", mock.Anything, mock.Anything).Return([]dto.ListCoreValuesResp{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.ListCoreValues(test.context, test.organizationId)

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
		organizationId  string
		coreValueId     string
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:           "Success for get corevalue",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:           "Invalid organisation",
			context:        context.Background(),
			organizationId: "0",
			coreValueId:    "1",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(apperrors.InvalidOrgId).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Invalid corevalue",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "0",
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, apperrors.InvalidCoreValueData).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.GetCoreValue(test.context, test.organizationId, test.coreValueId)

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
		organizationId  string
		userId          int64
		coreValue       dto.CreateCoreValueReq
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:           "Success for create corevalue",
			context:        context.Background(),
			organizationId: "1",
			userId:         1,
			coreValue: dto.CreateCoreValueReq{
				Text:         "CoreValue",
				Description:  "core value desc",
				ParentID:     nil,
				ThumbnailURL: "thumbnail url string",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
				coreValueMock.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:           "Invalid organisation",
			context:        context.Background(),
			organizationId: "0",
			userId:         1,
			coreValue: dto.CreateCoreValueReq{
				Text:         "CoreValue",
				Description:  "core value desc",
				ParentID:     nil,
				ThumbnailURL: "thumbnail url string",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(apperrors.InvalidOrgId).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Error while checking organisation",
			context:        context.Background(),
			organizationId: "0",
			userId:         1,
			coreValue: dto.CreateCoreValueReq{
				Text:         "CoreValue",
				Description:  "core value desc",
				ParentID:     nil,
				ThumbnailURL: "thumbnail url string",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything, mock.Anything).Return(true, apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Repeated core value",
			context:        context.Background(),
			organizationId: "1",
			userId:         1,
			coreValue: dto.CreateCoreValueReq{
				Text:         "CoreValue",
				Description:  "core value desc",
				ParentID:     nil,
				ThumbnailURL: "thumbnail url string",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Error while creating core value",
			context:        context.Background(),
			organizationId: "1",
			userId:         1,
			coreValue: dto.CreateCoreValueReq{
				Text:         "CoreValue",
				Description:  "core value desc",
				ParentID:     nil,
				ThumbnailURL: "thumbnail url string",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("CheckUniqueCoreVal", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
				coreValueMock.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.CreateCoreValue(test.context, test.organizationId, test.userId, test.coreValue)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestDeleteCoreValue(t *testing.T) {
	coreValueRepo := mocks.NewCoreValueStorer(t)
	service := NewService(coreValueRepo)

	tests := []struct {
		name            string
		context         context.Context
		organizationId  string
		coreValueId     string
		userId          int64
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:           "Success for delete corevalue",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			userId:         1,
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, nil).Once()
				coreValueMock.On("DeleteCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:           "Invalid organisation",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			userId:         1,
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(apperrors.InvalidOrgId).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Invalid corevalue id",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			userId:         1,
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, apperrors.InvalidCoreValueData).Once()
			},
			isErrorExpected: true,
		},
		{
			name:           "Error in deleting core value",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			userId:         1,
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, nil).Once()
				coreValueMock.On("DeleteCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			err := service.DeleteCoreValue(test.context, test.organizationId, test.coreValueId, test.userId)

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
		organizationId  string
		coreValueId     string
		reqData         dto.UpdateQueryRequest
		setup           func(coreValueMock *mocks.CoreValueStorer)
		isErrorExpected bool
	}{
		{
			name:           "Success for delete corevalue",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			reqData: dto.UpdateQueryRequest{
				Text:         "Updated core value",
				Description:  "updated description",
				ThumbnailUrl: "updated thumbnail url",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, nil).Once()
				coreValueMock.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.UpdateCoreValuesResp{}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:           "Invalid organisation",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			reqData: dto.UpdateQueryRequest{
				Text:         "Updated core value",
				Description:  "updated description",
				ThumbnailUrl: "updated thumbnail url",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(apperrors.InvalidOrgId).Once()

			},
			isErrorExpected: true,
		},
		{
			name:           "Invalid corevalue id",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			reqData: dto.UpdateQueryRequest{
				Text:         "Updated core value",
				Description:  "updated description",
				ThumbnailUrl: "updated thumbnail url",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, apperrors.InvalidCoreValueData).Once()
			},
			isErrorExpected: true,
		},
		{
			name:           "Error in updating core value",
			context:        context.Background(),
			organizationId: "1",
			coreValueId:    "1",
			reqData: dto.UpdateQueryRequest{
				Text:         "Updated core value",
				Description:  "updated description",
				ThumbnailUrl: "updated thumbnail url",
			},
			setup: func(coreValueMock *mocks.CoreValueStorer) {
				coreValueMock.On("CheckOrganisation", mock.Anything, mock.Anything).Return(nil).Once()
				coreValueMock.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, nil).Once()
				coreValueMock.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.UpdateCoreValuesResp{}, apperrors.InternalServerError).Once()

			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueRepo)

			// test service
			_, err := service.UpdateCoreValue(test.context, test.organizationId, test.coreValueId, test.reqData)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}
