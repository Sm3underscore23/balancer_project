package api

import (
	"encoding/json"
	"log"
	"net/http"

	"balancer/internal/model"
	"balancer/internal/service"
)

type Handler struct {
	pool            []*model.BackendServer
	tokenService    service.TokenService
	balanceStrategy service.BalanceStrategyService
	limitsManager   service.LimitsManagerService
}

func NewProxyHandler(
	pool []*model.BackendServer,
	tokenService service.TokenService,
	balanceStrategy service.BalanceStrategyService,
	limitsManager service.LimitsManagerService,
) *Handler {
	return &Handler{
		pool:            pool,
		tokenService:    tokenService,
		balanceStrategy: balanceStrategy,
		limitsManager:   limitsManager,
	}
}

func writeJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := model.ErrorResponse{
		Message: err.Error(),
	}

	w.WriteHeader(model.GetStatusCode(err))
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write JSONE: %s", err)
	}
}
