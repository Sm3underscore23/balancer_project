package api

import (
	"balancer/internal/model"
	"balancer/internal/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	pool *model.BackendPool
	srv  *service.Service
}

func NewProxyHandler(pool *model.BackendPool, srv *service.Service) *Handler {
	return &Handler{
		pool: pool,
		srv:  srv,
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	response := model.ErrorResponse{
		Message: message,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return model.ErrWriteMessage
	}

	return nil
}
