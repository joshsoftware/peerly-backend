package api

import (
	"bytes"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/app/coreValues/mocks"
	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

func TestListCoreValuesHandler(t *testing.T) {
	coreValueSvc := mocks.NewService(t)
	listCoreValuesHandler := listCoreValuesHandler(coreValueSvc)

	tests := []struct {
		name               string
		organisationId     int
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:           "Success for list corevalues",
			organisationId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything, mock.Anything).Return([]dto.ListCoreValuesResp{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "Wrong organisation id for list corevalues",
			organisationId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything, mock.Anything).Return([]dto.ListCoreValuesResp{}, apperrors.InvalidOrgId).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Error in vars string to int conversion",
			organisationId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything, mock.Anything).Return([]dto.ListCoreValuesResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "Error in vars ListCoreValues db functions",
			organisationId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything, mock.Anything).Return([]dto.ListCoreValuesResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("GET", fmt.Sprintf("/organisations/%d/core_values", test.organisationId), bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}

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

func TestGetCoreValueHandler(t *testing.T) {
	coreValueSvc := mocks.NewService(t)
	getCoreValueHandler := getCoreValueHandler(coreValueSvc)

	tests := []struct {
		name               string
		organisationId     int
		coreValueId        int
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:           "Success for get corevalue",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "Wrong organisation id for list corevalues",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, apperrors.InvalidOrgId).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Error in vars string to int conversion",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "Error in GetCoreValue db function",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.GetCoreValueResp{}, apperrors.InvalidCoreValueData).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("GET", fmt.Sprintf("/organisations/%d/core_values/%d", test.organisationId, test.coreValueId), bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getCoreValueHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestCreateCoreValueHandler(t *testing.T) {
	coreValueSvc := mocks.NewService(t)
	createCoreValueHandler := createCoreValueHandler(coreValueSvc)

	tests := []struct {
		name               string
		organisationId     int
		userId             int
		coreValue          string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:           "Success for create corevalue",
			organisationId: 1,
			userId:         1,
			coreValue: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:           "Text missing in json request",
			organisationId: 1,
			userId:         1,
			coreValue: `{
				"description": "desc",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, apperrors.TextFieldBlank).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Description missing in json request",
			organisationId: 1,
			userId:         1,
			coreValue: `{
				"text": "corevalue3",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, apperrors.DescFieldBlank).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Wrong parent id",
			organisationId: 1,
			userId:         1,
			coreValue: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc",
				"parent_id": 0
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, apperrors.InvalidParentValue).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Error in CreateCoreValues db function",
			organisationId: 1,
			userId:         1,
			coreValue: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc",
				"parent_id": 1
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "Wrong organisation id",
			organisationId: 0,
			userId:         1,
			coreValue: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc",
				"parent_id": 1
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.CreateCoreValueResp{}, apperrors.InvalidOrgId).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("POST", fmt.Sprintf("/organisations/%d/core_values", test.organisationId), bytes.NewBuffer([]byte(test.coreValue)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createCoreValueHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestDeleteCoreValueHandler(t *testing.T) {
	coreValueSvc := mocks.NewService(t)
	deleteCoreValueHandler := deleteCoreValueHandler(coreValueSvc)

	tests := []struct {
		name               string
		organisationId     int
		coreValueId        int
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:           "Success for delete corevalue",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "Wrong organisation id",
			organisationId: 0,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(apperrors.InvalidOrgId).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Wrong corevalue id",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(apperrors.InvalidCoreValueData).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Error in DeleteCoreValue db function",
			organisationId: 1,
			coreValueId:    1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/organisations/%d/core_values/%d", test.organisationId, test.coreValueId), bytes.NewBuffer([]byte("")))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(deleteCoreValueHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestUpdateCoreValueHandler(t *testing.T) {
	coreValueSvc := mocks.NewService(t)
	updateCoreValueHandler := updateCoreValueHandler(coreValueSvc)

	tests := []struct {
		name               string
		organisationId     int
		coreValueId        int
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:           "Success for update corevalue",
			organisationId: 1,
			coreValueId:    1,
			input: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.UpdateCoreValuesResp{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "Wrong organisation id",
			organisationId: 1,
			coreValueId:    1,
			input: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.UpdateCoreValuesResp{}, apperrors.InvalidOrgId).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Wrong corevalue id",
			organisationId: 1,
			coreValueId:    1,
			input: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.UpdateCoreValuesResp{}, apperrors.InvalidCoreValueData).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:           "Error in UpdateCoreValue db function",
			organisationId: 1,
			coreValueId:    1,
			input: `{
				"text": "corevalue3",
				"description": "desc",
				"thumbnail_url": "abc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dto.UpdateCoreValuesResp{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("PUT", fmt.Sprintf("/organisations/%d/core_values/%d", test.organisationId, test.coreValueId), bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateCoreValueHandler)
			handler.ServeHTTP(rr, req)

			fmt.Println("Error")

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
