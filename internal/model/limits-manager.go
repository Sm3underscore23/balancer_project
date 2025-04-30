package model

type DefaultClientLimits struct {
	Capacity   float64
	RatePerSec float64
}

type ClientLimits struct {
	ClientId   string  `json:"client_id" binding:"required"`
	Capacity   float64 `json:"capacity" binding:"required"`
	RatePerSec float64 `json:"rate_per_sec" binding:"required"`
}

type ClientIdRequest struct {
	ClientId string `json:"client_id"`
}
