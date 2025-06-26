package tokenmanager

import (
	"balancer/internal/model"
	"balancer/pkg/logger"
	"context"
	"time"
)

func refill(ctx context.Context, tb *model.TokenBucket) {
	if tb.TokenAmount() == tb.MaxTokens() {
		return
	}
	now := time.Now()
	elapsed := now.Sub(tb.LastRefillTime()).Seconds()

	tb.AddToken(elapsed * tb.RefillRate())

	if tb.TokenAmount() > tb.MaxTokens() {
		tb.SetToken(tb.MaxTokens())
		ctx = logger.AddValuesToContext(ctx, logger.ClientTokenAmount, tb.TokenAmount())
		logger.FromContext(ctx).Info("full refill")
	}

	tb.SetLastRefillTime(now)
}

func (s *tokenService) StartRefillWorker(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cache.Range(func(clientID string, bucket *model.TokenBucket) {
					ctx = logger.AddValuesToContext(ctx, logger.ClientID, clientID)
					refill(ctx, bucket)
				})
			}
		}
	}()
}
