package limitsmanager

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) UpdateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error {
	isExists, err := s.repo.IsClientExists(ctx, clientLimits.ClientId)
	if err != nil {
		return err
	}

	if !isExists {
		return model.ErrClientNotExists
	}

	err = s.repo.UpdateClientLimits(ctx, clientLimits)
	if err != nil {
		return err
	}

	if _, ok := s.cache.Get(clientLimits.ClientId); ok {
		s.cache.Set(clientLimits.ClientId, model.ConverClientLimitstoTB(clientLimits))
	}

	return nil
}
