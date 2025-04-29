package api

import (
	"net/http"

	"balancer/internal/model"
)

func getClientID(r *http.Request) string {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		return apiKey
	}
	return r.RemoteAddr
}

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.srv.RequestFromUser(ctx, getClientID(r))
	if err != nil {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	s, err := h.srv.BalanceStrategyRoundRobin(ctx)
	if err != nil {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	s.Prx.ServeHTTP(w, r)
}
