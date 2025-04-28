package poolservice

import (
	"balancer/internal/model"
	"balancer/internal/service/interfaces"
)

type poolService struct {
	pool *model.BackendPool
}

func NewPoolService(pool *model.BackendPool) interfaces.PoolService {
	return &poolService{pool: pool}
}
