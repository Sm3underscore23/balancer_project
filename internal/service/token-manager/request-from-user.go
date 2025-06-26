package tokenmanager

import (
	"context"
	"log/slog"

	"balancer/internal/model"
	"balancer/pkg/logger"
)

func logTokenGet(ctx context.Context, clientTb *model.TokenBucket) context.Context {
	ctx = logger.AddValuesToContext(ctx, logger.ClientTokenAmount, clientTb.TokenAmount())
	logger.FromContext(ctx).Info("get client token amount")
	return ctx
}

func logTokenRemain(ctx context.Context, clientTb *model.TokenBucket) {
	logger.FromContext(ctx).Info("remain client token amount", logger.ClientTokenAmount, clientTb.TokenAmount())
}

func (s *tokenService) RequestFromUser(ctx context.Context, clientID string) error {
	clientTb, ok := s.cache.Get(clientID)

	if ok {
		logTokenGet(ctx, clientTb)
		switch {
		case clientTb.TokenAmount() < 1:
			return model.ErrRateLimit

		default:
			clientTb.UseToken()
			logTokenRemain(ctx, clientTb)
			return nil
		}
	}

	isExists, err := s.db.IsClientExists(ctx, clientID)
	if err != nil {
		return err
	}

	if isExists {
		clientLimit, err := s.db.GetClientLimits(ctx, clientID)
		if err != nil {
			return err
		}

		clientTb := model.NewTokenBucket(clientLimit.Capacity, clientLimit.RatePerSec)

		s.cache.Set(clientID, clientTb)

		ctx = logTokenGet(ctx, clientTb)

		clientTb.UseToken()

		logTokenRemain(ctx, clientTb)

		return nil
	}

	clientTb = model.NewTokenBucket(s.defoultLimits.Capacity, s.defoultLimits.RatePerSec)
	err = s.db.CreateClientLimits(ctx, model.ClientLimits{
		ClientId:   clientID,
		Capacity:   clientTb.MaxTokens(),
		RatePerSec: clientTb.RefillRate()})

	if err != nil {
		return err
	}

	slog.Info("new client created", slog.Group(logger.GroupClientLimits,
		logger.ClientID, clientID,
		logger.ClinetTokenCapacity, clientTb.MaxTokens(),
		logger.ClinetRateRefill, clientTb.RefillRate(),
	))

	s.cache.Set(clientID, clientTb)

	ctx = logTokenGet(ctx, clientTb)

	clientTb.UseToken()

	logTokenRemain(ctx, clientTb)

	return nil
}
