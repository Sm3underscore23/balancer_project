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

// func refill(tb *model.TokenBucket) {
// 	now := time.Now()
// 	elapsed := now.Sub(tb.LastRefillTime).Seconds()
// 	if tb.Token() > tb.MaxTokens {
// 		tb.SetToken(tb.MaxTokens)
// 	}
// loop:
// 	for range int(min(math.Floor(tb.AddAndReturnToken(elapsed*tb.RefillRate)), float64(tb.MaxTokens))) {
// 		select {
// 		case tb.TokensChan <- struct{}{}:
// 		default:
// 			break loop
// 		}
// 	}
// 	tb.LastRefillTime = now
// }

// func (s *tokenService) StartRefillWorker(ctx context.Context) {
// 	ticker := time.NewTicker(time.Millisecond * 100)

// 	go func() {
// 		defer ticker.Stop()
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case <-ticker.C:
// 				for _, tb := range s.cache {
// 					refill(tb)
// 				}
// 			}
// 		}
// 	}()
// }
