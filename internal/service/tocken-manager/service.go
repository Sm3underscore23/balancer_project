package tockenmanager

import (
	"balancer/internal/model"
	"balancer/internal/repository"
	"balancer/internal/service/interfaces"
)

type tockenService struct {
	cache         map[string]*model.TokenBucket
	defoultLimits *model.DefoultUserLimits
	db            repository.Repository
}

func NewTockenService(db repository.Repository, defoultLimits *model.DefoultUserLimits) interfaces.TockenService {
	return &tockenService{cache: make(map[string]*model.TokenBucket), db: db, defoultLimits: defoultLimits}
}
