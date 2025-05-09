package roundrobbin

import (
	"context"
	"net/http/httputil"

	"balancer/internal/model"
)

func (r *roundRobbinService) Balance(ctx context.Context) (httputil.ReverseProxy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	l := uint64(len(r.Pool))

	for i := range l {
		idx := (r.current + i) % l
		if r.Pool[idx].IsOnline() {
			r.current = (uint64(idx) + 1)
			return *r.Pool[idx].Prx, nil
		}
	}

	return httputil.ReverseProxy{}, model.ErrNoAvilibleServers
}
