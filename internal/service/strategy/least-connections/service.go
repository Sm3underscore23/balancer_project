package leastconnections

import (
	"balancer/internal/model"
	"balancer/internal/service"
	"sync"
)

type leastConnectionsService struct {
	conList []uint64
	mu      sync.Mutex
	Pool    []*model.BackendServer
}

func NewLeastConnectionsService(pool []*model.BackendServer) service.BalanceStrategyService {
	return &leastConnectionsService{conList: make([]uint64, len(pool)), Pool: pool}
}
