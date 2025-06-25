package api

import (
	"balancer/pkg/logger"
	"net/http"
)

func getClientID(r *http.Request) string {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		return apiKey
	}
	return r.RemoteAddr
}

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	clientID := getClientID(r)

	ctx := logger.AddValuesToContext(r.Context(),
		logger.ClientID, clientID,
	)

	logger.FromContext(ctx).Info("handler Proxy started")

	err := h.tokenService.RequestFromUser(ctx, clientID)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}
	prx, err := h.balanceStrategy.Balance(ctx)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}
	prx.Proxy().ServeHTTP(w, r)

	ctx = logger.AddValuesToContext(r.Context(), logger.StatusCode, http.StatusOK)
	logger.FromContext(ctx).Info("handler Proxy complete")
}
