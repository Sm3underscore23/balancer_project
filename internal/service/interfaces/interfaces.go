package interfaces

import (
	"context"
	"time"

	"balancer/internal/model"
)

type TockenService interface { // Token
	RequestFromUser(ctx context.Context, ip string) error
}

type PoolService interface {
	CheckerWithTicker(ctx context.Context, t *time.Ticker) error

	BalanceStrategyRoundRobin(ctx context.Context) (*model.BackendServer, error)
}

type LimitsManagerService interface {
	UpdateUserLimits(ctx context.Context, clientLimits *model.UserLimits) error
}
