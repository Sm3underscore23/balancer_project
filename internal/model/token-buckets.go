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

func (tb *TokenBucket) Refill() {
	now := time.Now()
	elapsed := now.Sub(tb.LastRefillTime).Seconds()

	tb.Tokens += elapsed * tb.RefillRate
	if tb.Tokens > tb.MaxTokens {
		tb.Tokens = tb.MaxTokens
	}
	tb.LastRefillTime = now
}

func ConverUserLimitstoTB(userLimits UserLimits) *TokenBucket {
	return &TokenBucket{
		Tokens:         userLimits.Capacity,
		MaxTokens:      userLimits.Capacity,
		RefillRate:     userLimits.RatePerSec,
		LastRefillTime: time.Now(),
	}
}
