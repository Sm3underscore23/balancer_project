package service

import "balancer/internal/model"

type CheckerService interface {
	CheckerOnce()
	CheckerWithTicker()
}

type BalancerService interface {
	BalanceStrategyRoundRobin() model.BackendServer
}
