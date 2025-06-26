package roundrobbin

import (
	"balancer/internal/model"
	"balancer/pkg/logger"
	"context"
)

func (r *roundRobbinService) Balance(ctx context.Context) (model.Proxy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	l := uint64(len(r.pool))

	for i := range l {
		idx := (r.current + i) % l
		if !r.pool[idx].IsOnline() {
			continue
		}

		r.current = (uint64(idx) + 1)
		prx := r.pool[idx]
		ctx = logger.AddValuesToContext(ctx, logger.BackendUrl, r.pool[idx].BckndUrl)
		logger.FromContext(ctx).Info("get backend address")

		return prx, nil
	}

	return nil, model.ErrNoAvilibleServers
}
