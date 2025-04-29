package model

import (
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	tokens         float64
	MaxTokens      float64
	RefillRate     float64
	mu             sync.Mutex
	TokensChan     chan struct{}
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

func (tb *TokenBucket) Token() float64 {
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

func ConverClientLimitstoTB(clientLimits ClientLimits) *TokenBucket {
	return &TokenBucket{
		MaxTokens:      clientLimits.Capacity,
		RefillRate:     clientLimits.RatePerSec,
		LastRefillTime: time.Now(),
	}
}

func NewTokenBucket(maxTokens float64, refillRate float64) *TokenBucket {
	tb := &TokenBucket{
		tokens:         maxTokens,
		MaxTokens:      maxTokens,
		RefillRate:     refillRate,
		LastRefillTime: time.Now(),
	}
	return tb

}
