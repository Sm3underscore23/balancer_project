package tockenmanager

import (
	"context"
	"time"

	"balancer/internal/model"
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
	// mutex
	userTb, ok := s.cache[clientId]
	if !ok {
		userTb.Mu.Lock()
		defer userTb.Mu.Unlock()

		select {
		case <-userTb.TokensChan:
			return nil
		default:
			return model.ErrRateLimit
		}

	}
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
		return nil
	}

	tb := newTokenBucket(s.defoultLimits.Capacity, s.defoultLimits.RatePerSec)
	_, err = s.db.CreateUserLimits(ctx, model.ConverTBtoUserLimits(clientId, tb))

	if err != nil {
		return err
	}

	s.cache[clientId] = tb

	userTb = tb
	return nil
	// я бы убрал рефилл здесь, пусть будет только по тикеру
	refill(userTb)
	if userTb.Tokens > 0 {
		userTb.Tokens -= 1
		return nil
	}
	return model.ErrRateLimit
}
