package leastconnections

import (
	"balancer/internal/model"
	"context"
	"math"
)

func (l *leastConnectionsService) Balance(ctx context.Context) (model.Proxy, error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	var (
		minConections uint64 = math.MaxUint64
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

	if minConections != math.MaxUint64 {
		l.pool[index].numOfCon += 1
		prx := l.pool[index]
		return prx, nil
	}

	return nil, model.ErrNoAvilibleServers
}
