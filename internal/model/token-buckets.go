package model

import (
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	Mu             sync.Mutex
	Tokens         float64
	MaxTokens      float64
	RefillRate     float64
	LastRefillTime time.Time

	TokensChan chan struct{}
}

func (t *TokenBucket) UseToken() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.Tokens -= 1
}

func (t *TokenBucket) AddToken() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.Tokens += 1
}

func ConverUserLimitstoTB(userLimits UserLimits) *TokenBucket {
	return &TokenBucket{
		Tokens:         userLimits.Capacity,
		MaxTokens:      userLimits.Capacity,
		RefillRate:     userLimits.RatePerSec,
		LastRefillTime: time.Now(),
	}
}
