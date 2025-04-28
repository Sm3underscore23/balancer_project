package model

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"errors"`
}

var (
	ErrCreateProxy = fmt.Errorf("create nil proxy")

	ErrObjectNotExists  = fmt.Errorf("no user object")
	ErrMethodNotAllowed = fmt.Errorf("method not allowed")
	ErrRateLimit        = fmt.Errorf("rate limit exceeded")

	ErrDbBuilder         = fmt.Errorf("failed to build query error")
	ErrDbScan            = fmt.Errorf("failed to scan query")
	ErrDbQuery           = fmt.Errorf("failed to exec query")
	ErrWriteMessage      = fmt.Errorf("failed to write json message")
	ErrNoAvilibleServers = fmt.Errorf("no healthy backends available")

	ErrWithStatus = map[error]int{
		ErrObjectNotExists:   http.StatusBadRequest,
		ErrMethodNotAllowed:  http.StatusMethodNotAllowed,
		ErrRateLimit:         http.StatusTooManyRequests,
		ErrDbBuilder:         http.StatusInternalServerError,
		ErrDbScan:            http.StatusInternalServerError,
		ErrDbQuery:           http.StatusInternalServerError,
		ErrWriteMessage:      http.StatusInternalServerError,
		ErrNoAvilibleServers: http.StatusServiceUnavailable,
	}
)
