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

func TestCreateClientLimits(t *testing.T) {
	testTable := []struct {
		name          string
		mockRepo      func(controller *gomock.Controller) repository.LimitsRepository
		expectedError error
	}{
		{
			name: "OK",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, nil)
				mock.EXPECT().CreateClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "test-client-id",
					Capacity:   10,
					RatePerSec: 1,
				}).Return(nil)
				return mock
			},
			expectedError: nil,
		},
		{
			name: "client already exists",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(true, nil)
				return mock
			},
			expectedError: model.ErrClientAlreadyExists,
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
			name: "db error: failed to create new client",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, nil)
				mock.EXPECT().CreateClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "test-client-id",
					Capacity:   10,
					RatePerSec: 1,
				}).Return(model.ErrDb)
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

			err := limitsManager.CreateClientLimits(ctx, model.ClientLimits{
				ClientId:   "test-client-id",
				Capacity:   10,
				RatePerSec: 1,
			})

			assert.Equal(t, testCase.expectedError, err)

			// сделать проверку кеша
		})
	}
}
