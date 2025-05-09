package limitsmanagergo

import (
	"balancer/internal/repository"
	"balancer/internal/service"
	inmemorycache "balancer/internal/service/in-memory-cache"
)

type limitsManagerService struct {
	cache *inmemorycache.InMemoryTokenBucketCache
	repo  repository.LimitsRepository
}

func NewPoolService(
	cache *inmemorycache.InMemoryTokenBucketCache,
	repo repository.LimitsRepository) service.LimitsManagerService {
	return &limitsManagerService{cache: cache, repo: repo}
}
