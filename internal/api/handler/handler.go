package api

import (
	"encoding/json"
	"log"
	"net/http"

	"balancer/internal/model"
	"balancer/internal/service"
)

type Handler struct {
	pool []*model.BackendServer
	srv  *service.Service
}

func NewProxyHandler(pool []*model.BackendServer, srv *service.Service) *Handler {
	return &Handler{
		pool: pool,
		srv:  srv,
	}
}

func writeJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := model.ErrorResponse{
		Message: err.Error(),
	}

	w.WriteHeader(model.GetStatusCode(err))
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to write JSONE: %s", err)
	}
}
