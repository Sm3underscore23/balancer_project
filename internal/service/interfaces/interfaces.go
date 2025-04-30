package interfaces

import (
	"context"
	"time"

	"balancer/internal/model"
)

type TokenService interface {
	RequestFromUser(ctx context.Context, ip string) error
	StartRefillWorker(ctx context.Context)
}

type PoolService interface {
	CheckerWithTicker(ctx context.Context, t *time.Ticker) error
	BalanceStrategyRoundRobin(ctx context.Context) (*model.BackendServer, error)
}

type LimitsManagerService interface {
	CreateClientLimits(ctx context.Context, clientLimits *model.ClientLimits) error
	GetClientLimits(ctx context.Context, clientId string) (model.ClientLimits, error)
	UpdateClientLimits(ctx context.Context, clientLimits *model.ClientLimits) error
	DeleteClientLimits(ctx context.Context, clientId string) error
}
