package api

import (
	"balancer/internal/model"
	"balancer/internal/service"
	mock_service "balancer/internal/service/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProxy(t *testing.T) {
	testBackendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer testBackendServer.Close()

	testBackendURL, _ := url.Parse(testBackendServer.URL)
	testProxy := *httputil.NewSingleHostReverseProxy(testBackendURL)

	testTable := []struct {
		name             string
		clientId         string
		isClientIdAPIKey bool
		mockToken        func(controller *gomock.Controller) service.TokenService
		mockBalance      func(controller *gomock.Controller) service.BalanceStrategyService
		expectedStatus   int
		isError          bool
		expectedBody     string
	}{
		{
			name:             "success with API key",
			clientId:         "test-api-key",
			isClientIdAPIKey: true,
			mockToken: func(controller *gomock.Controller) service.TokenService {
				mock := mock_service.NewMockTokenService(controller)
				mock.EXPECT().RequestFromUser(gomock.Any(), "test-api-key").Return(nil)
				return mock
			},
			mockBalance: func(controller *gomock.Controller) service.BalanceStrategyService {
				mock := mock_service.NewMockBalanceStrategyService(controller)
				mock.EXPECT().Balance(gomock.Any()).Return(testProxy, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			isError:        false,
		},
		{
			name:             "success with remote address",
			clientId:         "123.4.5.6",
			isClientIdAPIKey: false,
			mockToken: func(controller *gomock.Controller) service.TokenService {
				mock := mock_service.NewMockTokenService(controller)
				mock.EXPECT().RequestFromUser(gomock.Any(), "123.4.5.6").Return(nil)
				return mock
			},
			mockBalance: func(controller *gomock.Controller) service.BalanceStrategyService {
				mock := mock_service.NewMockBalanceStrategyService(controller)
				mock.EXPECT().Balance(gomock.Any()).Return(testProxy, nil)
				return mock
			},
			expectedStatus: http.StatusOK,
			isError:        false,
		},
		{
			name:             "rate limit error",
			clientId:         "123.4.5.6",
			isClientIdAPIKey: false,
			mockToken: func(controller *gomock.Controller) service.TokenService {
				mock := mock_service.NewMockTokenService(controller)
				mock.EXPECT().RequestFromUser(gomock.Any(), "123.4.5.6").Return(model.ErrRateLimit)
				return mock
			},
			mockBalance: func(controller *gomock.Controller) service.BalanceStrategyService {
				return mock_service.NewMockBalanceStrategyService(controller)
			},
			expectedStatus: http.StatusTooManyRequests,
			isError:        true,
			expectedBody:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrRateLimit),
		},
		{
			name:             "no av",
			clientId:         "123.4.5.6",
			isClientIdAPIKey: false,
			mockToken: func(controller *gomock.Controller) service.TokenService {
				mock := mock_service.NewMockTokenService(controller)
				mock.EXPECT().RequestFromUser(gomock.Any(), "123.4.5.6").Return(nil)
				return mock
			},
			mockBalance: func(controller *gomock.Controller) service.BalanceStrategyService {
				mock := mock_service.NewMockBalanceStrategyService(controller)
				mock.EXPECT().Balance(gomock.Any()).Return(httputil.ReverseProxy{}, model.ErrNoAvilibleServers)
				return mock
			},
			expectedStatus: http.StatusServiceUnavailable,
			isError:        true,
			expectedBody:   fmt.Sprintf("{\"errors\":\"%s\"}\n", model.ErrNoAvilibleServers),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)

			if testCase.isClientIdAPIKey {
				req.Header.Set("X-API-Key", testCase.clientId)
			} else {
				req.RemoteAddr = testCase.clientId
			}

			w := httptest.NewRecorder()

			controller := gomock.NewController(t)

			handler := NewProxyHandler(nil, testCase.mockToken(controller), testCase.mockBalance(controller), nil)

			handler.Proxy(w, req)

			if testCase.isError {
				assert.Equal(t, testCase.expectedBody, w.Body.String())
			}

			assert.Equal(t, testCase.expectedStatus, w.Result().StatusCode)
		})
	}
}
