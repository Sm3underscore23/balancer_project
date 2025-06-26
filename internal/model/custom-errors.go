package model

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"errors"`
}

var (
	ErrParseConfig = fmt.Errorf("failed to parse config")

	ErrCreateProxy = fmt.Errorf("create nil proxy")

	ErrInvalidInput = fmt.Errorf("invalid input data")

	ErrClientNotExists     = fmt.Errorf("client not exists,  please use servise for create")
	ErrClientAlreadyExists = fmt.Errorf("clinet already exists, please use servise for update")
	ErrMethodNotAllowed    = fmt.Errorf("method not allowed")
	ErrRateLimit           = fmt.Errorf("rate limit exceeded")

	ErrWriteMessage      = fmt.Errorf("failed to write json message")
	ErrNoAvilibleServers = fmt.Errorf("no healthy backends available")

	ErrDb = fmt.Errorf("db error")

	errWithStatus = map[error]int{
		ErrInvalidInput:        http.StatusBadRequest,
		ErrClientNotExists:     http.StatusBadRequest,
		ErrClientAlreadyExists: http.StatusBadRequest,
		ErrMethodNotAllowed:    http.StatusMethodNotAllowed,
		ErrRateLimit:           http.StatusTooManyRequests,
		ErrNoAvilibleServers:   http.StatusServiceUnavailable,
	}
)

func GetStatusCode(err error) int {
	for mapError, code := range errWithStatus {
		if errors.Is(err, mapError) {
			return code
		}
	}
	return http.StatusInternalServerError
}
