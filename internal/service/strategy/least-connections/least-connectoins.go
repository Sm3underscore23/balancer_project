package leastconnections

import (
	"balancer/internal/model"
	"balancer/pkg/logger"
	"context"
)

func (l *leastConnectionsService) Balance(ctx context.Context) (model.Proxy, error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	var (
		minConections = l.pool[0].numOfCon
		idx           int
	)

	for i, b := range l.pool {
		if !b.IsOnline() {
			continue
		}

		if minConections > b.numOfCon {
			minConections = b.numOfCon
			idx = i
		}
	}

	if l.pool[idx].IsOnline() {
		l.pool[idx].numOfCon += 1
		prx := l.pool[idx]
		ctx = logger.AddValuesToContext(ctx, logger.BackendUrl, l.pool[idx].BckndUrl)
		logger.FromContext(ctx).Info("get backend address")

		return prx, nil
	}

	return nil, model.ErrNoAvilibleServers
}
