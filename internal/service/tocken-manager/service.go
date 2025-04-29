package tockenmanager

import (
	"sync"

	"balancer/internal/model"
	"balancer/internal/repository"
	"balancer/internal/service/interfaces"

	"golang.org/x/time/rate"
)

type tockenService struct {
	cache         map[string]*model.TokenBucket
	mu            sync.Mutex
	defoultLimits *model.DefoultUserLimits // default - typo
	db            repository.LimitsRepository
}

func NewTockenService(db repository.LimitsRepository, defoultLimits *model.DefoultUserLimits) interfaces.TockenService {
	return &tockenService{cache: make(map[string]*model.TokenBucket), db: db, defoultLimits: defoultLimits}
}
