package leastconnections

import (
	"balancer/internal/model"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	testTable := []struct {
		name          string
		newPool       func() []*model.BackendServer
		inputCurrent  uint64
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
					{
						BackendServerSettings: model.BackendServerSettings{
							BckndUrl: "backend 2",
						},
					},
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
			expectedError: nil,
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
			pool := testCase.newPool()
			balancer := LeastConnectionsService(pool)

			ctx := context.Background()

			prx, err := balancer.Balance(ctx)

			if err != nil {
				assert.Equal(t, testCase.expectedError, err)
				assert.Nil(t, prx)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, prx.URL(), "backend 1")

			prx, err = balancer.Balance(ctx)

			assert.NoError(t, err)
			assert.Equal(t, prx.URL(), "backend 3")

			pool[1].ChangeHealthStatus(true)

			prx, err = balancer.Balance(ctx)

			assert.NoError(t, err)
			assert.Equal(t, prx.URL(), "backend 2")
		})
	}
}
