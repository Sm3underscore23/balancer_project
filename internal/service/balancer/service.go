package checker

import (
	"balancer/internal/model"
	"balancer/internal/service"
)

type poolChecker struct {
	pool *model.BackendPool
}

func NewCheckerService(pool *model.BackendPool) service.CheckerService {
	return &poolChecker{pool: pool}
}
