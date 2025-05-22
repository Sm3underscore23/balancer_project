package checker

import (
	"balancer/internal/model"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckerWithProgressiveHealth(t *testing.T) {

	var globalCheckCount int

	type serverWithCounter struct {
		server    *httptest.Server
		callCount int
		id        int
	}

	newTestServer := func(id int) *serverWithCounter {
		callCount := 0

		handler := func(w http.ResponseWriter, r *http.Request) {
			callCount++
			if globalCheckCount >= id {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		}

		srv := httptest.NewServer(http.HandlerFunc(handler))

		return &serverWithCounter{
			server:    srv,
			callCount: callCount,
			id:        id,
		}
	}

	countOnline := func(pool []*model.BackendServer) int {
		count := 0
		for _, b := range pool {
			if b.IsOnline() {
				count++
			}
		}
		return count
	}

	servers := []*serverWithCounter{
		newTestServer(1),
		newTestServer(2),
		newTestServer(3),
	}

	pool := make([]*model.BackendServer, 0, len(servers))
	for _, s := range servers {
		pool = append(pool, &model.BackendServer{
			BackendServerSettings: model.BackendServerSettings{
				BckndUrl: s.server.URL,
				Method:   "GET",
				HelthUrl: s.server.URL,
			},
		})
	}

	ctx := context.Background()
	ch := checker{pool: pool}

	ch.check(ctx)
	assert.Equal(t, 0, countOnline(pool))

	globalCheckCount += 1
	ch.check(ctx)
	assert.Equal(t, 1, countOnline(pool))

	globalCheckCount += 1
	ch.check(ctx)
	assert.Equal(t, 2, countOnline(pool))

	globalCheckCount += 1
	ch.check(ctx)
	assert.Equal(t, 3, countOnline(pool))
}
