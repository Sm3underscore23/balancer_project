package limitsmanagergo

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
		return model.ClientLimits{}, model.ErrUserNotExists
	}

	clientLimits, err := s.repo.GetUserLimits(ctx, clientId)
	if err != nil {
		return model.ClientLimits{}, err
	}

	if _, ok := s.cache.Get(clientId); !ok {
		s.cache.Set(clientId, model.ConverClientLimitstoTB(clientLimits))
	}

	return clientLimits, nil
}
