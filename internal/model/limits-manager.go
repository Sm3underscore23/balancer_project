package model

type DefoultUserLimits struct {
	Capacity   float64
	RatePerSec float64
}

type UserLimits struct {
	ClientId   string  `json:"client_id" binding:"required"`
	Capacity   float64 `json:"capacity" binding:"required"`
	RatePerSec float64 `json:"rate_per_sec" binding:"required"`
}

func ConverTBtoUserLimits(clientId string, tb *TokenBucket) *UserLimits {
	return &UserLimits{
		ClientId:   clientId,
		Capacity:   tb.MaxTokens,
		RatePerSec: tb.RefillRate,
	}
}
