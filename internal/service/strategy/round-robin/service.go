package roundrobbin

import (
	"sync"

	"balancer/internal/model"
	"balancer/internal/service"
)

type roundRobbinService struct {
	current uint64
	mu      sync.Mutex
	Pool    []*model.BackendServer
}

func RoundRobbinService(pool []*model.BackendServer) service.BalanceStrategyService {
	return &roundRobbinService{Pool: pool}
}
