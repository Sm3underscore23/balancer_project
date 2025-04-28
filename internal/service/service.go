package service

import (
	"balancer/internal/model"
	"balancer/internal/repository"
	"balancer/internal/service/interfaces"
	limitsmanagergo "balancer/internal/service/limits-manager.go"
	poolservice "balancer/internal/service/pool-service"
	tockenmanager "balancer/internal/service/tocken-manager"
)

type Service struct {
	interfaces.PoolService
	interfaces.TockenService
	interfaces.LimitsManagerService
}

func NewService(pool *model.BackendPool, db repository.Repository, defoultLimits *model.DefoultUserLimits) *Service {
	return &Service{
		PoolService:          poolservice.NewPoolService(pool),
		TockenService:        tockenmanager.NewTockenService(db, defoultLimits),
		LimitsManagerService: limitsmanagergo.NewPoolService(db),
	}
}
