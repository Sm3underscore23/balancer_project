package model

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens         float64
	maxTokens      float64
	refillRate     float64
	mu             sync.Mutex
	lastRefillTime time.Time
}

func (tb *TokenBucket) TokenAmount() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

func (tb *TokenBucket) UseToken() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens -= 1
}

func (tb *TokenBucket) AddToken(a float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens += a
}

func (tb *TokenBucket) AddAndReturnToken(a float64) float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens += a
	return tb.tokens
}

func (tb *TokenBucket) SetToken(a float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens = a
}

func (tb *TokenBucket) MaxTokens() float64 {
	return tb.maxTokens
}

func (tb *TokenBucket) RefillRate() float64 {
	return tb.refillRate
}

func (tb *TokenBucket) LastRefillTime() time.Time {
	return tb.lastRefillTime
}

func (tb *TokenBucket) SetLastRefillTime(t time.Time) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.lastRefillTime = t
}

func ConvertClientLimitstoTB(clientLimits ClientLimits) *TokenBucket {
	return &TokenBucket{
		maxTokens:      clientLimits.Capacity,
		refillRate:     clientLimits.RatePerSec,
		lastRefillTime: time.Now(),
	}
}

func NewTokenBucket(maxTokens float64, refillRate float64) *TokenBucket {
	tb := &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
	return tb

}
