package api

import (
	"balancer/internal/model"
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
	err := h.srv.RequestFromUser(r.Context(), getClientID(r))
	if err != nil {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	s, err := h.srv.BalanceStrategyRoundRobin(r.Context())
	if err != nil {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	s.Prx.ServeHTTP(w, r)
}
