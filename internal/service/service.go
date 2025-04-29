package service

import (
	"balancer/internal/model"
	"balancer/internal/repository"
	inmemorycache "balancer/internal/service/in-memory-cache"
	"balancer/internal/service/interfaces"
	limitsmanagergo "balancer/internal/service/limits-manager"
	poolservice "balancer/internal/service/pool-service"
	tockenmanager "balancer/internal/service/tocken-manager"
)

// я бы это удалил, сервис сервисов не нужен
type Service struct {
	interfaces.PoolService
	interfaces.TokenService
	interfaces.LimitsManagerService
}

func NewService(
	cache *inmemorycache.InMemoryTokenBucketCache,
	pool []*model.BackendServer,
	db repository.LimitsRepository,
	defaultLimits *model.DefaultClientLimits) *Service {
	return &Service{
		PoolService:          poolservice.NewPoolService(pool),
		TokenService:         tockenmanager.NewTockenService(cache, db, defaultLimits),
		LimitsManagerService: limitsmanagergo.NewPoolService(cache, db),
	}
}
