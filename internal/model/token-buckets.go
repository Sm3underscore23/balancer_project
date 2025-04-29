package model

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Mu             sync.Mutex
	Tokens         float64
	MaxTokens      float64
	RefillRate     float64
	LastRefillTime time.Time
}

func ConverUserLimitstoTB(userLimits UserLimits) *TokenBucket {
	return &TokenBucket{
		Tokens:         userLimits.Capacity,
		MaxTokens:      userLimits.Capacity,
		RefillRate:     userLimits.RatePerSec,
		LastRefillTime: time.Now(),
	}
}
