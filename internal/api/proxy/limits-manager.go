package api

import (
	"encoding/json"
	"net/http"

	"balancer/internal/model"
)

func (h *Handler) UpdateLimits(w http.ResponseWriter, r *http.Request) { // переименовать // вынести в два хендлера отделных
	if r.Method != http.MethodPost {
		err := model.ErrMethodNotAllowed
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
	r.Body.Close()

	err = h.srv.UpdateUserLimits(r.Context(), &userLimits)
	if err != nil {
		if err := writeJSONError(w, model.ErrWithStatus[err], err.Error()); err != nil {
			http.Error(w, err.Error(), model.ErrWithStatus[err])
			return
		}
		return
	}

	w.WriteHeader(http.StatusAccepted) // status
}
