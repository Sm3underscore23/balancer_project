package api

import (
	"balancer/internal/model"
	"balancer/internal/service"
	mock_service "balancer/internal/service/mocks"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateLimits(t *testing.T) {
	validInput := fmt.Sprintf("{ \"client_id\": \"%s\", \"capacity\": %v, \"rate_per_sec\": %v }", "123.4.5.6", 10, 1)

	testTable := []struct {
		name           string
		input          string
		mockBehavior   func(c *gomock.Controller) service.LimitsManagerService
		expectedStatus int
		expectedData   string
	}{
		{
			name:  "OK",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().CreateClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "123.4.5.6",
					Capacity:   10,
					RatePerSec: 1,
				}).Return(nil)
				return mock
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:  "invalid input data",
			input: "{ \"invalid\": \"1234\" }",
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				return mock_service.NewMockLimitsManagerService(c)
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrInvalidInput),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "user already exists",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().CreateClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "123.4.5.6",
					Capacity:   10,
					RatePerSec: 1,
				}).Return(model.ErrClientAlreadyExists)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrClientAlreadyExists),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "db error",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().CreateClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "123.4.5.6",
					Capacity:   10,
					RatePerSec: 1,
				}).Return(model.ErrDb)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrDb),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/limits/create", strings.NewReader(testCase.input))

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := NewProxyHandler(nil, nil, nil, testCase.mockBehavior(controller))

			handler.CreateLimits(w, req)
			log.Printf(testCase.expectedData, w.Body.String())
			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
			require.Equal(t, testCase.expectedData, w.Body.String())
		})
	}
}

func TestGetLimits(t *testing.T) {
	validInput := fmt.Sprintf("{ \"client_id\": \"%s\"}", "123.4.5.6")

	testTable := []struct {
		name           string
		input          string
		mockBehavior   func(c *gomock.Controller) service.LimitsManagerService
		expectedStatus int
		expectedData   string
	}{
		{
			name:  "OK",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().GetClientLimits(
					gomock.Any(),
					"123.4.5.6").
					Return(model.ClientLimits{
						ClientId:   "123.4.5.6",
						Capacity:   10,
						RatePerSec: 1,
					},
						nil)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"client_id\":\"%s\",\"capacity\":%v,\"rate_per_sec\":%v}\n","123.4.5.6", 10, 1),
			expectedStatus: http.StatusOK,
		},
		{
			name:  "invalid input data",
			input: "{ \"invalid\": \"1234\" }",
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				return mock_service.NewMockLimitsManagerService(c)
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrInvalidInput),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "user not exists",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().GetClientLimits(
					gomock.Any(),
					"123.4.5.6").
					Return(model.ClientLimits{}, model.ErrClientNotExists)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrClientNotExists),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "db error",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().GetClientLimits(
					gomock.Any(),
					"123.4.5.6").
					Return(model.ClientLimits{}, model.ErrDb)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrDb),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/limits/get", strings.NewReader(testCase.input))

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := NewProxyHandler(nil, nil, nil, testCase.mockBehavior(controller))

			handler.GetLimits(w, req)

			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
			require.Equal(t, testCase.expectedData, w.Body.String())
		})
	}
}

func TestUpdateLimits(t *testing.T) {
	validInput := fmt.Sprintf("{ \"client_id\": \"%s\", \"capacity\": %v, \"rate_per_sec\": %v }", "123.4.5.6", 10, 1)

	testTable := []struct {
		name           string
		input          string
		mockBehavior   func(c *gomock.Controller) service.LimitsManagerService
		expectedData   string
		expectedStatus int
	}{
		{
			name:  "OK",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().UpdateClientLimits(
					gomock.Any(),
					model.ClientLimits{
						ClientId:   "123.4.5.6",
						Capacity:   10,
						RatePerSec: 1}).
					Return(nil)
				return mock
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:  "invalid input data",
			input: "{ \"invalid\": \"1234\" }",
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				return mock_service.NewMockLimitsManagerService(c)
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrInvalidInput),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "user not exists",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().UpdateClientLimits(
					gomock.Any(),
					model.ClientLimits{
						ClientId:   "123.4.5.6",
						Capacity:   10,
						RatePerSec: 1}).
					Return(model.ErrClientNotExists)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrClientNotExists),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "db error",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().UpdateClientLimits(
					gomock.Any(),
					model.ClientLimits{
						ClientId:   "123.4.5.6",
						Capacity:   10,
						RatePerSec: 1}).
					Return(model.ErrDb)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrDb),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/limits/update", strings.NewReader(testCase.input))

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := NewProxyHandler(nil, nil, nil, testCase.mockBehavior(controller))

			handler.UpdateLimits(w, req)

			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
			require.Equal(t, testCase.expectedData, w.Body.String())
		})
	}
}

func TestDeleteLimits(t *testing.T) {
	validInput := fmt.Sprintf("{ \"client_id\": \"%s\"}", "123.4.5.6")

	testTable := []struct {
		name           string
		input          string
		mockBehavior   func(c *gomock.Controller) service.LimitsManagerService
		expectedData   string
		expectedStatus int
	}{
		{
			name:  "OK",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().DeleteClientLimits(
					gomock.Any(), "123.4.5.6").
					Return(nil)
				return mock
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:  "invalid input data",
			input: "{ \"invalid\": \"1234\" }",
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				return mock_service.NewMockLimitsManagerService(c)
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrInvalidInput),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "user not exists",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().DeleteClientLimits(
					gomock.Any(), "123.4.5.6").
					Return(model.ErrClientNotExists)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrClientNotExists),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "db error",
			input: validInput,
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().DeleteClientLimits(
					gomock.Any(), "123.4.5.6").
					Return(model.ErrDb)
				return mock
			},
			expectedData:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrDb),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/limits/delete", strings.NewReader(testCase.input))

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := NewProxyHandler(nil, nil, nil, testCase.mockBehavior(controller))

			handler.DeleteLimits(w, req)

			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
			require.Equal(t, testCase.expectedData, w.Body.String())
		})
	}
}
