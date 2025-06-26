package tokenmanager

import (
	"balancer/internal/model"
	"balancer/internal/repository"
	"balancer/internal/service"
	inmemorycache "balancer/internal/service/in-memory-cache"
)

type tokenService struct {
	cache         *inmemorycache.InMemoryTokenBucketCache
	defoultLimits model.DefaultClientLimits
	db            repository.LimitsRepository
}

func NewTockenService(
	cache *inmemorycache.InMemoryTokenBucketCache,
	db repository.LimitsRepository,
	defoultLimits model.DefaultClientLimits) service.TokenService {
	return &tokenService{cache: cache, db: db, defoultLimits: defoultLimits}
}

