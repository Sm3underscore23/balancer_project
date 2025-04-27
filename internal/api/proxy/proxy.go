package api

import (
	"balancer/internal/model"
	"fmt"
	"net/http"
)

func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("some text")
	s := h.srv.BalanceStrategyRoundRobin()
	if s == (model.BackendServer{}) {
		http.Error(w, "all servers are dead", http.StatusInternalServerError)
		return
	}
	fmt.Println(s)
	s.Prx.ServeHTTP(w, r)
}
