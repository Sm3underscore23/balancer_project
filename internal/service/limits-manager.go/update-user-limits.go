package limitsmanagergo

import (
	"balancer/internal/model"
	"context"
	"log"
)

func (s *limitsManagerService) UpdateUserLimits(ctx context.Context, clientLimits *model.UserLimits) error {
	log.Println("1")
	isExists, err := s.repo.IsClientExists(ctx, clientLimits.ClientId)
	if err != nil {
		return err
	}

	if !isExists {
		return model.ErrObjectNotExists
	}

	err = s.repo.UpdateUserLimits(ctx, clientLimits)
	if err != nil {
		return err
	}

	return nil
}
