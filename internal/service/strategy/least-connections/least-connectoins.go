package leastconnections

import (
	"balancer/internal/model"
	"context"
)

func (l *leastConnectionsService) Balance(ctx context.Context) (model.Proxy, error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	var (
		minConections = l.pool[0].numOfCon
		index         int
	)

	for i, b := range l.pool {
		if !b.IsOnline() {
			continue
		}

		if minConections > b.numOfCon {
			minConections = b.numOfCon
			index = i
		}
	}

	if l.pool[index].IsOnline() {
		l.pool[index].numOfCon += 1
		prx := l.pool[index]
		return prx, nil
	}

	return nil, model.ErrNoAvilibleServers
}
