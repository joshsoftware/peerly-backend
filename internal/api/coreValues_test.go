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
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for list corevalues",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything).Return([]dto.CoreValue{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Error in vars string to int conversion",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything).Return([]dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Error in vars ListCoreValues db functions",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListCoreValues", mock.Anything).Return([]dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("GET", "/core_values", bytes.NewBuffer([]byte("")))
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
		coreValueId        int
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success for get corevalue",
			coreValueId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Error in vars string to int conversion",
			coreValueId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:        "Error in GetCoreValue db function",
			coreValueId: 1,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetCoreValue", mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InvalidCoreValueData).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("GET", fmt.Sprintf("/core_values/%d", test.coreValueId), bytes.NewBuffer([]byte("")))
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
		userId             int
		coreValue          string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:   "Success for create corevalue",
			userId: 1,
			coreValue: `{
				"name": "corevalue3",
				"description": "desc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:   "Name missing in json request",
			userId: 1,
			coreValue: `{
				"description": "desc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.TextFieldBlank).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Description missing in json request",
			userId: 1,
			coreValue: `{
				"name": "corevalue3"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.DescFieldBlank).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:   "Wrong parent id",
			userId: 1,
			coreValue: `{
				"name": "corevalue3",
				"description": "desc",
				"parent_core_value_id": 0
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InvalidParentValue).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:   "Error in CreateCoreValues db function",
			userId: 1,
			coreValue: `{
				"name": "corevalue3",
				"description": "desc",
				"parent_core_value_id": 1
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("POST", "/core_values", bytes.NewBuffer([]byte(test.coreValue)))
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

func TestUpdateCoreValueHandler(t *testing.T) {
	coreValueSvc := mocks.NewService(t)
	updateCoreValueHandler := updateCoreValueHandler(coreValueSvc)

	tests := []struct {
		name               string
		coreValueId        int
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:        "Success for update corevalue",
			coreValueId: 1,
			input: `{
				"name": "corevalue3",
				"description": "desc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Wrong corevalue id",
			coreValueId: 1,
			input: `{
				"name": "corevalue3",
				"description": "desc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InvalidCoreValueData).Once()
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:        "Error in UpdateCoreValue db function",
			coreValueId: 1,
			input: `{
				"name": "corevalue3",
				"description": "desc"
			}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateCoreValue", mock.Anything, mock.Anything, mock.Anything).Return(dto.CoreValue{}, apperrors.InternalServerError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(coreValueSvc)

			req, err := http.NewRequest("PUT", fmt.Sprintf("/core_values/%d", test.coreValueId), bytes.NewBuffer([]byte(test.input)))
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
