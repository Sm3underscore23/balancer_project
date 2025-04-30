package tockenmanager

import (
	"balancer/internal/model"
	"balancer/internal/repository"
	inmemorycache "balancer/internal/service/in-memory-cache"
	"balancer/internal/service/interfaces"
)

type tokenService struct {
	cache         *inmemorycache.InMemoryTokenBucketCache
	defoultLimits *model.DefaultClientLimits
	db            repository.LimitsRepository
}

func NewTockenService(
	cache *inmemorycache.InMemoryTokenBucketCache,
	db repository.LimitsRepository,
	defoultLimits *model.DefaultClientLimits) interfaces.TokenService {
	return &tokenService{cache: cache, db: db, defoultLimits: defoultLimits}
}
