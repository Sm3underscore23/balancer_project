package limitsmanager

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) UpdateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error {

	if _, ok := s.cache.Get(clientLimits.ClientId); !ok {

		isExists, err := s.repo.IsClientExists(ctx, clientLimits.ClientId)
		if err != nil {
			return err
		}

		if !isExists {
			return model.ErrClientNotExists
		}

	}

	err := s.repo.UpdateClientLimits(ctx, clientLimits)
	if err != nil {
		return err
	}

	if currClientBucket, ok := s.cache.Get(clientLimits.ClientId); ok {
		updatedClientBucket := model.ConvertClientLimitstoTB(clientLimits)
		updatedClientBucket.SetToken(currClientBucket.TokenAmount())
		s.cache.Set(clientLimits.ClientId, updatedClientBucket)
	}

	return nil
}
