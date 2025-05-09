package api

import (
	"balancer/internal/model"
	"balancer/internal/service"
	mock_service "balancer/internal/service/mocks"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateLimits(t *testing.T) {

	testTamble := []struct {
		name           string
		input          string
		mockBehavior   func(c *gomock.Controller) service.LimitsManagerService
		expectedData   string
		expectedStatus int
	}{
		{
			name:  "OK",
			input: "{ \"client_id\": \"123.5.6.7\", \"capacity\": 10, \"rate_per_sec\": 1 }",
			mockBehavior: func(c *gomock.Controller) service.LimitsManagerService {
				mock := mock_service.NewMockLimitsManagerService(c)
				mock.EXPECT().CreateClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "123.5.6.7",
					Capacity:   10,
					RatePerSec: 1,
				}).Return(nil)
				return mock
			},
			expectedData:   "",
			expectedStatus: 201,
		},
	}

	for _, testCase := range testTamble {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/limits/create", strings.NewReader(testCase.input))

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := NewProxyHandler(nil, nil, nil, testCase.mockBehavior(controller))

			handler.CreateLimits(w, req)

			require.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
			require.Equal(t, testCase.expectedData, w.Body.String())
		})
	}
}
