package poolservice

import (
	"sync"

	"balancer/internal/model"
	"balancer/internal/service/interfaces"
)

type poolService struct {
	current uint64
	mu      sync.RWMutex
	Pool    []*model.BackendServer
}

func NewPoolService(pool *model.BackendPool) interfaces.PoolService {
	return &poolService{pool: pool}
}
