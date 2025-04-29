package poolservice

import (
	"context"

	"balancer/internal/model"
)

func (p *poolService) BalanceStrategyRoundRobin(ctx context.Context) (*model.BackendServer, error) {

	l := uint64(len(p.pool.Pool))

	for i := range l {
		idx := (p.pool.Current + i) % l
		if p.pool.Pool[idx].IsOnline.Load() {
			p.pool.Current = uint64(idx) + 1
			return p.pool.Pool[idx], nil
		}
	}

	return nil, model.ErrNoAvilibleServers
}
