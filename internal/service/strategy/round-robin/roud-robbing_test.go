package roundrobbin

import (
	"balancer/internal/model"
	"context"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	testTable := []struct {
		name          string
		newPool       func() []*model.BackendServer
		inputCurrent  uint64
		expectedPrx   httputil.ReverseProxy
		expectedError error
	}{
		{
			name: "OK",
			newPool: func() []*model.BackendServer {
				pool := []*model.BackendServer{
					{
						BackendServerSettings: model.BackendServerSettings{
							BckndUrl: "backend 1",
						},
					},
					{},
					{
						BackendServerSettings: model.BackendServerSettings{
							BckndUrl: "backend 3",
						},
					},
				}
				pool[0].ChangeHealthStatus(true)
				pool[2].ChangeHealthStatus(true)

				return pool
			},
		},
		{
			name: "no healthy backends available",
			newPool: func() []*model.BackendServer {
				pool := []*model.BackendServer{
					{},
					{},
					{},
				}
				return pool
			},
			expectedError: model.ErrNoAvilibleServers,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			balancer := RoundRobbinService(testCase.newPool())

			ctx := context.Background()

			prx, err := balancer.Balance(ctx)

			assert.Equal(t, testCase.expectedError, err)

			if err != nil {
				assert.Nil(t, prx)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, prx.URL(), "backend 1")

			prx, err = balancer.Balance(ctx)

			assert.NoError(t, err)
			assert.Equal(t, prx.URL(), "backend 3")

			prx, err = balancer.Balance(ctx)

			assert.NoError(t, err)
			assert.Equal(t, prx.URL(), "backend 1")

		})
	}
}
