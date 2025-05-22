package tokenmanager

import (
	"balancer/internal/model"
	"context"
	"time"
)

func refill(tb *model.TokenBucket) {
	now := time.Now()
	elapsed := now.Sub(tb.LastRefillTime()).Seconds()

	tb.AddToken(elapsed * tb.RefillRate())

	if tb.TokenAmount() > tb.MaxTokens() {
		tb.SetToken(tb.MaxTokens())
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
				s.cache.Range(func(clientId string, bucket *model.TokenBucket) {
					refill(bucket)
				})
			}
		}
	}()
}
