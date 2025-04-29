package tockenmanager

import (
	"balancer/internal/model"
	"context"
	"time"
)

func refill(tb *model.TokenBucket) {
	now := time.Now()
	elapsed := now.Sub(tb.LastRefillTime).Seconds()

	tb.Tokens += elapsed * tb.RefillRate
	if tb.Tokens > tb.MaxTokens {
		tb.Tokens = tb.MaxTokens
	}
	tb.LastRefillTime = now
}

func (s *tockenService) StartRefillWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, tb := range s.cache {
					tb.Mu.Lock()
					refill(tb)
					tb.Mu.Unlock()
				}
			}
		}
	}()
}
