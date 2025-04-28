package tockenmanager

import (
	"balancer/internal/model"
	"context"
	"time"
)

func newTokenBucket(maxTokens float64, refillRate float64) *model.TokenBucket {
	return &model.TokenBucket{
		Tokens:         maxTokens,
		MaxTokens:      maxTokens,
		RefillRate:     refillRate,
		LastRefillTime: time.Now(),
	}
}

func (s *tockenService) RequestFromUser(ctx context.Context, clientId string) error {
	userTb, ok := s.cache[clientId]
	if !ok {
		isExists, err := s.db.IsClientExists(ctx, clientId)
		if err != nil {
			return err
		}

		if isExists {
			userLimit, err := s.db.GetUserLimits(ctx, clientId)
			if err != nil {
				return err
			}

			tb := model.ConverUserLimitstoTB(userLimit)

			s.cache[clientId] = tb
			userTb = tb
		}

		if !isExists {
			tb := newTokenBucket(s.defoultLimits.Capacity, s.defoultLimits.RatePerSec)
			_, err := s.db.CreateUserLimits(ctx, model.ConverTBtoUserLimits(clientId, tb))

			if err != nil {
				return err
			}

			s.cache[clientId] = tb

			userTb = tb
		}
	}

	userTb.Mu.Lock()
	defer userTb.Mu.Unlock()

	userTb.Refill()
	if userTb.Tokens > 0 {
		userTb.Tokens -= 1
		return nil
	}
	return model.ErrRateLimit
}
