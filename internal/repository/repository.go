package repository

import (
	"context"

	"balancer/internal/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type LimitsRepository interface {
	CreateUserLimits(ctx context.Context, info model.ClientLimits) error
	GetUserLimits(ctx context.Context, clientId string) (model.ClientLimits, error)
	IsClientExists(ctx context.Context, userId string) (bool, error)
	UpdateUserLimits(ctx context.Context, updateData model.ClientLimits) error
	DeleteUserLimits(ctx context.Context, clientId string) error
}
