package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"balancer/internal/model"
	"balancer/internal/service"
	"balancer/pkg/logger"
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

func writeJSONError(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := model.ErrorResponse{
		Message: err.Error(),
	}
	statusCode := model.GetStatusCode(err)
	w.WriteHeader(statusCode)
	
	ctx = logger.AddValuesToContext(ctx,
		logger.StatusCode, statusCode,
		logger.Error, err,
	)
	logger.FromContext(ctx).Error("handler error")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write JSONE: %s", err)
	}
}
