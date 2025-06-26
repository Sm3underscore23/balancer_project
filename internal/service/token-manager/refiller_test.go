package tokenmanager

import (
	"balancer/internal/model"
	mock_repository "balancer/internal/repository/mocks"
	inmemorycache "balancer/internal/service/in-memory-cache"
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStartRefillWorker(t *testing.T) {
	testTable := []struct {
		name                string
		clientId            string
		ratePerSec          float64
		workingTime         int
		expectedNumOfTokens float64
	}{
		{
			name:                "1 seconds - ratePerSec 1 - expected 1",
			clientId:            "test-client-id",
			ratePerSec:          1,
			workingTime:         1,
			expectedNumOfTokens: 1,
		},
		{
			name:                "10 seconds - ratePerSec 0.1 - expected 1",
			clientId:            "test-client-id",
			ratePerSec:          0.2,
			workingTime:         5,
			expectedNumOfTokens: 1,
		},
		{
			name:                "11 seconds - ratePerSec 0.1 - expected 1.1",
			clientId:            "test-client-id",
			ratePerSec:          0.2,
			workingTime:         6,
			expectedNumOfTokens: 1.2,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			emptyBucket := model.NewTokenBucket(math.MaxFloat64, testCase.ratePerSec)

			cache := inmemorycache.NewInMemoryTokenBucketCache()
			cache.Set(testCase.clientId, emptyBucket)

			tokensInCacheBucket, ok := cache.Get(testCase.clientId)

			assert.True(t, ok)
			assert.Zero(t, tokensInCacheBucket.TokenAmount())

			controller := gomock.NewController(t)
			repo := mock_repository.NewMockLimitsRepository(controller)

			tockenService := NewTockenService(cache, repo, model.DefaultClientLimits{})

			ctx, cancel := context.WithCancel(context.Background())

			tockenService.StartRefillWorker(ctx)
			time.Sleep(time.Millisecond * 50)
			time.Sleep(time.Second * time.Duration(testCase.workingTime))
			cancel()

			tokensInCacheBucket, ok = cache.Get(testCase.clientId)

			assert.True(t, ok)
			assert.InDelta(t, testCase.expectedNumOfTokens, tokensInCacheBucket.TokenAmount(), 0.01)
		})
	}
}
