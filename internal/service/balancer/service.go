package balancer

import (
	"balancer/internal/model"
	"balancer/internal/service"
)

type poolBalancer struct {
	pool *model.BackendPool
}

func NewBalancerService(pool *model.BackendPool) service.BalancerService {
	return &poolBalancer{pool: pool}
}
