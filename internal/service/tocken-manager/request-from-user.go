package tockenmanager

import (
	"context"

	"balancer/internal/model"
)

func (s *tokenService) RequestFromUser(ctx context.Context, clientId string) error {
	userTb, ok := s.cache.Get(clientId)
	if ok {
		switch {
		case userTb.Token() < 1:
			return model.ErrRateLimit
		default:
			userTb.UseToken()
			return nil
		}

	}
	isExists, err := s.db.IsClientExists(ctx, clientId)
	if err != nil {
		return err
	}

	if isExists {
		clientLimit, err := s.db.GetUserLimits(ctx, clientId)
		if err != nil {
			return err
		}

		tb := model.NewTokenBucket(clientLimit.Capacity, clientLimit.RatePerSec)

		s.cache.Set(clientId, tb)
		tb.UseToken()
		return nil
	}

	tb := model.NewTokenBucket(s.defoultLimits.Capacity, s.defoultLimits.RatePerSec)
	err = s.db.CreateUserLimits(ctx, &model.ClientLimits{
		ClientId:   clientId,
		Capacity:   tb.MaxTokens,
		RatePerSec: tb.RefillRate})

	if err != nil {
		return err
	}

	s.cache.Set(clientId, tb)

	tb.UseToken()

	return nil
}
