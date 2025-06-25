package limitsmanager

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) CreateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error {

	if _, ok := s.cache.Get(clientLimits.ClientId); !ok {
		return model.ErrClientAlreadyExists
	}

	isExists, err := s.repo.IsClientExists(ctx, clientLimits.ClientId)
	if err != nil {
		return err
	}

	if isExists {
		return err
	}

	err = s.repo.CreateClientLimits(ctx, clientLimits)

	if err != nil {
		return err
	}

	s.cache.Set(clientLimits.ClientId, model.ConvertClientLimitstoTB(clientLimits))

	return nil
}
