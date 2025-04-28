package limitsmanagergo

import (
	"balancer/internal/repository"
	"balancer/internal/service/interfaces"
)

type limitsManagerService struct {
	repo repository.Repository
}

func NewPoolService(repo repository.Repository) interfaces.LimitsManagerService {
	return &limitsManagerService{repo: repo}
}
