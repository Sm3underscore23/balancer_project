package api

import (
	"balancer/internal/model"
	"balancer/internal/service"
)

type Handler struct {
	srv  service.BalancerService
	pool *model.BackendPool
}

func NewUserImplementation(srv service.BalancerService, pool *model.BackendPool) *Handler {
	return &Handler{
		srv:  srv,
		pool: pool,
	}
}
