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
	ClientId string `json:"client_id" binding:"required"`
}

func (c *ClientLimits) ValidateClientLimits() error {
	if c.ClientId == "" {
		return ErrInvalidInput
	}
	if c.Capacity == 0 {
		return ErrInvalidInput
	}
	if c.RatePerSec == 0 {
		return ErrInvalidInput
	}
	return nil
}

func (c *ClientIdRequest) ValidateClientIdRequest() error {
	if c.ClientId == "" {
		return ErrInvalidInput
	}
	return nil
}
