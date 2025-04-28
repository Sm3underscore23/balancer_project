package repository

import (
	"balancer/internal/model"
	"context"
)

type Repository interface {
	CreateUserLimits(ctx context.Context, info *model.UserLimits) (string, error)
	GetUserLimits(ctx context.Context, clientId string) (model.UserLimits, error)
	IsClientExists(ctx context.Context, userId string) (bool, error)
	UpdateUserLimits(ctx context.Context, updateData *model.UserLimits) error
}
