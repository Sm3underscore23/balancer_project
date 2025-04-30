package repository

import (
	"context"

	"balancer/internal/model"
)

type LimitsRepository interface {
	CreateUserLimits(ctx context.Context, info *model.ClientLimits) error
	GetUserLimits(ctx context.Context, clientId string) (model.ClientLimits, error)
	IsClientExists(ctx context.Context, userId string) (bool, error)
	UpdateUserLimits(ctx context.Context, updateData *model.ClientLimits) error
	DeleteUserLimits(ctx context.Context, clientId string) error
}
