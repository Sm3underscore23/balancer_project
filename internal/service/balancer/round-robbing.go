package balancer

import (
	"balancer/internal/model"
)

func (p *poolBalancer) BalanceStrategyRoundRobin() model.BackendServer {

	l := uint64(len(p.pool.Pool))

	for i := range l {
		idx := (p.pool.Current + i) % l
		if p.pool.Pool[idx].IsOnline {
			p.pool.Current = uint64(idx) + 1
			return *p.pool.Pool[idx]
		}
	}

	return model.BackendServer{}
}
