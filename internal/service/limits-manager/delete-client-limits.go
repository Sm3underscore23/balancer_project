package limitsmanagergo

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) DeleteClientLimits(ctx context.Context, clientId string) error {
	isExists, err := s.repo.IsClientExists(ctx, clientId)
	if err != nil {
		return err
	}

	if !isExists {
		return model.ErrUserNotExists
	}

	err = s.repo.DeleteUserLimits(ctx, clientId)
	if err != nil {
		return err
	}

	if _, ok := s.cache.Get(clientId); !ok {
		s.cache.Delete(clientId)
	}

	return nil
}
