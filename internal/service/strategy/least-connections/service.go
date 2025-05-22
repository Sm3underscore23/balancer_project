package leastconnections

import (
	"balancer/internal/model"
	"balancer/internal/service"
	"sync"
)

type backendServerWithNumOfCon struct {
	numOfCon uint64
	*model.BackendServer
}

type leastConnectionsService struct {
	mu   sync.Mutex
	pool []*backendServerWithNumOfCon
}

func LeastConnectionsService(pool []*model.BackendServer) service.BalanceStrategyService {
	wPool := make([]*backendServerWithNumOfCon, len(pool))
	for i, b := range pool {
		wPool[i] = &backendServerWithNumOfCon{
			BackendServer: b,
		}
	}
	return &leastConnectionsService{pool: wPool}
}

func (l *leastConnectionsService) Lock() {
	l.mu.Lock()
}

func (l *leastConnectionsService) UnLock() {
	l.mu.Unlock()
}
