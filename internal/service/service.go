package service

import (
	"context"

	"balancer/internal/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type TokenService interface {
	RequestFromUser(ctx context.Context, ip string) error
	StartRefillWorker(ctx context.Context)
}

type Checker interface {
	CheckerWithTicker(ctx context.Context, rate uint64) error
}

type BalanceStrategyService interface {
	Balance(ctx context.Context) (model.Proxy, error)
}

type LimitsManagerService interface {
	CreateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error
	GetClientLimits(ctx context.Context, clientId string) (model.ClientLimits, error)
	UpdateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error
	DeleteClientLimits(ctx context.Context, clientId string) error
}
