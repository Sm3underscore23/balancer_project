package tockenmanager

import (
	"context"
	"time"

	"balancer/internal/model"
)

func refill(tb *model.TokenBucket) {
	now := time.Now()
	elapsed := now.Sub(tb.LastRefillTime).Seconds()

	for range int(min(elapsed*tb.RefillRate, tb.MaxTokens)) {
		tb.TokensChan <- struct{}{}
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
				// мьютекс
				for _, tb := range s.cache {
					tb.Mu.Lock()
					refill(tb)
					tb.Mu.Unlock()
				}
			}
		}
	}()
}
