package tokenmanager

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

func TestRequestFromUser(t *testing.T) {
	defoultLimits := model.DefaultClientLimits{
		Capacity:   10,
		RatePerSec: 1,
	}

	testTable := []struct {
		name          string
		clientId      string
		isRateLimit   bool
		isCacheEmpty  bool
		mockRepo      func(controller *gomock.Controller) repository.LimitsRepository
		expectedError error
	}{
		{
			name:     "OK - cache not empty",
			clientId: "test-client-id",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				return mock_repository.NewMockLimitsRepository(controller)
			},
			expectedError: nil,
		},
		{
			name:         "OK - empty cashe",
			clientId:     "test-client-id",
			isCacheEmpty: true,
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(true, nil)
				mock.EXPECT().GetClientLimits(gomock.Any(), "test-client-id").Return(model.ClientLimits{
					ClientId:   "123.4.5.6",
					Capacity:   10,
					RatePerSec: 1,
				},
					nil,
				)
				return mock
			},
			expectedError: nil,
		},
		{
			name:         "OK - new client",
			clientId:     "test-client-id",
			isCacheEmpty: true,
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, nil)
				mock.EXPECT().GetClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "test-client-id",
					Capacity:   10,
					RatePerSec: 1}).
					Return(nil)
				return mock
			},
			expectedError: nil,
		},
		{
			name:        "rate limit",
			clientId:    "test-client-id",
			isRateLimit: true,
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				return mock_repository.NewMockLimitsRepository(controller)
			},
			expectedError: model.ErrRateLimit,
		},
		{
			name:         "db error: failed to get IsClientExists",
			isCacheEmpty: true,
			clientId:     "test-client-id",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, model.ErrDb)
				return mock
			},
			expectedError: model.ErrDb,
		},
		{
			name:         "db error: failed to get client limits",
			isCacheEmpty: true,
			clientId:     "test-client-id",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(true, nil)
				mock.EXPECT().GetClientLimits(gomock.Any(), "test-client-id").Return(model.ClientLimits{},
					model.ErrDb,
				)
				return mock
			},
			expectedError: model.ErrDb,
		},
		{
			name:         "db error: failed to create new client",
			isCacheEmpty: true,
			clientId:     "test-client-id",
			mockRepo: func(controller *gomock.Controller) repository.LimitsRepository {
				mock := mock_repository.NewMockLimitsRepository(controller)
				mock.EXPECT().IsClientExists(gomock.Any(), "test-client-id").Return(false, nil)
				mock.EXPECT().GetClientLimits(gomock.Any(), model.ClientLimits{
					ClientId:   "test-client-id",
					Capacity:   10,
					RatePerSec: 1}).
					Return(model.ErrDb)
				return mock
			},
			expectedError: model.ErrDb,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			repo := testCase.mockRepo(controller)

			cache := inmemorycache.NewInMemoryTokenBucketCache()

			if !testCase.isCacheEmpty && !testCase.isRateLimit { // сделать функцию
				cache.Set(testCase.clientId, model.NewTokenBucket(defoultLimits.Capacity, defoultLimits.RatePerSec))
			}

			if !testCase.isCacheEmpty && testCase.isRateLimit {
				cache.Set(testCase.clientId, &model.TokenBucket{})
			}

			tokenService := NewTockenService(cache, repo, defoultLimits)

			err := tokenService.RequestFromUser(context.Background(), testCase.clientId)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
