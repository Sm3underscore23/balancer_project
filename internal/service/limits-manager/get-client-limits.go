package limitsmanager

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) GetClientLimits(ctx context.Context, clientID string) (model.ClientLimits, error) {

	if tb, ok := s.cache.Get(clientID); ok {
		return model.ClientLimits{
			ClientId:   clientID,
			Capacity:   tb.MaxTokens(),
			RatePerSec: tb.RefillRate(),
		}, nil
	}

	isExists, err := s.repo.IsClientExists(ctx, clientID)
	if err != nil {
		return model.ClientLimits{}, err
	}

	if !isExists {
		return model.ClientLimits{}, model.ErrClientNotExists
	}

	clientLimits, err := s.repo.GetClientLimits(ctx, clientID)
	if err != nil {
		return model.ClientLimits{}, err
	}

	if _, ok := s.cache.Get(clientID); !ok {
		s.cache.Set(clientID, model.ConvertClientLimitstoTB(clientLimits))
	}

	return clientLimits, nil
}
