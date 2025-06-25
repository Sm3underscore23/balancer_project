package limitsmanager

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) DeleteClientLimits(ctx context.Context, clientID string) error {

	if _, ok := s.cache.Get(clientID); !ok {
		s.cache.Delete(clientID)

		err := s.repo.DeleteClientLimits(ctx, clientID)
		if err != nil {
			return err
		}

		return nil
	}

	isExists, err := s.repo.IsClientExists(ctx, clientID)
	if err != nil {
		return err
	}

	if !isExists {
		return model.ErrClientNotExists
	}

	err = s.repo.DeleteClientLimits(ctx, clientID)
	if err != nil {
		return err
	}

	if _, ok := s.cache.Get(clientID); !ok {
		s.cache.Delete(clientID)
	}

	return nil
}
