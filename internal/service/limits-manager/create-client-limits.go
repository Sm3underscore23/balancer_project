package limitsmanagergo

import (
	"balancer/internal/model"
	"context"
)

func (s *limitsManagerService) CreateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error {
	isExists, err := s.repo.IsClientExists(ctx, clientLimits.ClientId)
	if err != nil {
		return err
	}

	if isExists {
		return model.ErrUserAlreadyExists
	}

	err = s.repo.CreateUserLimits(ctx, clientLimits)

	if err != nil {
		return err
	}

	return nil
}
