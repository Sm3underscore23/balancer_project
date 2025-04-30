package limitsmanagergo

import (
	"balancer/internal/repository"
	inmemorycache "balancer/internal/service/in-memory-cache"
	"balancer/internal/service/interfaces"
)

type limitsManagerService struct {
	cache *inmemorycache.InMemoryTokenBucketCache
	repo  repository.LimitsRepository
}

func NewPoolService(
	cache *inmemorycache.InMemoryTokenBucketCache,
	repo repository.LimitsRepository) interfaces.LimitsManagerService {
	return &limitsManagerService{cache: cache, repo: repo}
}
