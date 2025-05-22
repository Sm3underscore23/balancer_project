package roundrobbin

import (
	"balancer/internal/model"
	"context"
)

func (r *roundRobbinService) Balance(ctx context.Context) (model.Proxy, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	l := uint64(len(r.Pool))

	for i := range l {
		idx := (r.current + i) % l
		if r.Pool[idx].IsOnline() {
			r.current = (uint64(idx) + 1)
			prx := r.Pool[idx]
			return prx, nil
		}
	}

	return nil, model.ErrNoAvilibleServers
}

// не верю в опыт по рассказу - подача
// уверенное написание кода -
// подкрепление позитивного результата
