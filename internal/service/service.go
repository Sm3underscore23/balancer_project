package service

import (
	"context"
	"net/http/httputil"
	"time"

	"balancer/internal/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type TokenService interface {
	RequestFromUser(ctx context.Context, ip string) error
	StartRefillWorker(ctx context.Context)
}

type BalanceStrategyService interface {
	CheckerWithTicker(ctx context.Context, t *time.Ticker) error
	Balance(ctx context.Context) (httputil.ReverseProxy, error)
}

type LimitsManagerService interface {
	CreateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error
	GetClientLimits(ctx context.Context, clientId string) (model.ClientLimits, error)
	UpdateClientLimits(ctx context.Context, clientLimits model.ClientLimits) error
	DeleteClientLimits(ctx context.Context, clientId string) error
}
