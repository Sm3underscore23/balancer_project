package limitsmanager

import (
	"balancer/internal/model"
	"balancer/internal/repository"
	mock_repository "balancer/internal/repository/mocks"
	inmemorycache "balancer/internal/service/in-memory-cache"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetClientLimits(t *testing.T) {
	testTable := []struct {
		name          string
		mockRepo      func(controller *gomock.Controller) repository.LimitsRepository
		expectedData  model.ClientLimits
		expectedError error
	}{
		{
			name: "OK",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(true, nil)
				mock.EXPECT().GetClientLimits(gomock.Any(), "test-client-id").Return(
					model.ClientLimits{
						ClientId:   "test-client-limits",
						Capacity:   10,
						RatePerSec: 1,
					},
					nil)
				return mock
			},
			expectedData: model.ClientLimits{
				ClientId:   "test-client-limits",
				Capacity:   10,
				RatePerSec: 1,
			},
			expectedError: nil,
		},
		{
			name: "client not exists",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, nil)
				return mock
			},
			expectedError: model.ErrClientNotExists,
		},
		{
			name: "db error: failed to get if client is exists",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, model.ErrDb)
				return mock
			},
			expectedError: model.ErrDb,
		},
		{
			name: "db error: failed to delete client",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(true, nil)
				mock.EXPECT().GetClientLimits(gomock.Any(), "test-client-id").Return(
					model.ClientLimits{},
					model.ErrDb)
				return mock
			},
			expectedError: model.ErrDb,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			controller := gomock.NewController(t)

			limitsManager := NewLimitsManagerService(inmemorycache.NewInMemoryTokenBucketCache(), testCase.mockRepo(controller))

			clientLimits, err := limitsManager.GetClientLimits(ctx, "test-client-id")

			assert.Equal(t, testCase.expectedData, clientLimits)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
