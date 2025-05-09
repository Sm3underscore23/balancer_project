package tockenmanager

import (
	"balancer/internal/model"
	"context"
	"time"
)

func refill(tb *model.TokenBucket) {
	now := time.Now()
	elapsed := now.Sub(tb.LastRefillTime).Seconds()
	if tb.Token() > tb.MaxTokens {
		tb.SetToken(tb.MaxTokens)
	} else {
		tb.AddToken(elapsed * tb.RefillRate)
	}
	tb.LastRefillTime = now
}

func (s *tokenService) StartRefillWorker(ctx context.Context) {
	ticker := time.NewTicker(time.Microsecond * 100)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cache.Range(func(clientID string, bucket *model.TokenBucket) {
					refill(bucket)
				})
			}
		}
	}()
}