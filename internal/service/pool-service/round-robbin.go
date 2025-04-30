package poolservice

import (
	"context"

	"balancer/internal/model"
)

func (p *poolService) BalanceStrategyRoundRobin(ctx context.Context) (*model.BackendServer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	l := uint64(len(p.Pool))

	for i := range l {
		idx := (p.current + i) % l
		if p.Pool[idx].Load() {
			p.current = (uint64(idx) + 1)
			return p.Pool[idx], nil
		}
	}

	return nil, model.ErrNoAvilibleServers
}
