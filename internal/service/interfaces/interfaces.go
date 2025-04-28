package interfaces

import (
	"balancer/internal/model"
	"context"
	"time"
)

type TockenService interface {
	RequestFromUser(ctx context.Context, ip string) error
}

type PoolService interface {
	CheckerWithTicker(ctx context.Context, t *time.Ticker) error

	BalanceStrategyRoundRobin(ctx context.Context) (*model.BackendServer, error)
}

type LimitsManagerService interface {
	UpdateUserLimits(ctx context.Context, clientLimits *model.UserLimits) error
}
