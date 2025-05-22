package repository

import (
	"context"

	"balancer/internal/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type LimitsRepository interface {
	CreateClientLimits(ctx context.Context, info model.ClientLimits) error
	GetClientLimits(ctx context.Context, clientId string) (model.ClientLimits, error)
	IsClientExists(ctx context.Context, userId string) (bool, error)
	UpdateClientLimits(ctx context.Context, updateData model.ClientLimits) error
	DeleteClientLimits(ctx context.Context, clientId string) error
}
