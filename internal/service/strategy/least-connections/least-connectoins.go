package leastconnections

import (
	"balancer/internal/model"
	"context"
	"math"
	"net/http/httputil"
)

func (l *leastConnectionsService) Balance(ctx context.Context) (httputil.ReverseProxy, error) {

	l.mu.Lock()
	defer l.mu.Unlock()

	var (
		minConections uint64 = math.MaxUint64
		index         int
	)

	for i, b := range l.Pool {
		if !b.IsOnline() {
			continue
		}

		numOfCon := l.conList[i]

		if minConections > numOfCon {
			minConections = numOfCon
			index = i
		}
	}

	if minConections != math.MaxUint64 {
		l.conList[index] += 1
		return *l.Pool[index].Prx, nil
	}

	return httputil.ReverseProxy{}, model.ErrNoAvilibleServers
}
