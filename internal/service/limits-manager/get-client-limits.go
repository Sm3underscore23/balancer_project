package limitsmanager

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) GetClientLimits(ctx context.Context, clientId string) (model.ClientLimits, error) {
	isExists, err := s.repo.IsClientExists(ctx, clientId)
	if err != nil {
		return model.ClientLimits{}, err
	}

	if !isExists {
		return model.ClientLimits{}, model.ErrClientNotExists
	}

	clientLimits, err := s.repo.GetClientLimits(ctx, clientId)
	if err != nil {
		return model.ClientLimits{}, err
	}

	if _, ok := s.cache.Get(clientId); !ok {
		s.cache.Set(clientId, model.ConvertClientLimitstoTB(clientLimits))
	}

	return clientLimits, nil
}
