package api

import (
	"balancer/internal/model"
	"encoding/json"
	"net/http"
)

func (h *Handler) LimitsManager(w http.ResponseWriter, r *http.Request) {
	if err := model.ErrMethodNotAllowed; r.Method != http.MethodPost {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), model.ErrWithStatus[err])
			return
		}
		return
	}

	var userLimits model.UserLimits
	err := json.NewDecoder(r.Body).Decode(&userLimits)
	if err != nil {
		if err := writeJSONError(w, http.StatusInternalServerError, err.Error()); err != nil {
			http.Error(w, err.Error(), model.ErrWithStatus[err])
			return
		}
		return
	}

	err = h.srv.UpdateUserLimits(r.Context(), &userLimits)

	if err != nil {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), model.ErrWithStatus[err])
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
