package api

import (
	"encoding/json"
	"net/http"

	"balancer/internal/model"
)

func (h *Handler) CreateLimits(w http.ResponseWriter, r *http.Request) {
	var clientLimits model.ClientLimits
	err := json.NewDecoder(r.Body).Decode(&clientLimits)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	r.Body.Close()

	if err = clientLimits.ValidateClientLimits(); err != nil {
		writeJSONError(w, err)
		return
	}

	err = h.limitsManager.CreateClientLimits(r.Context(), clientLimits)
	if err != nil {
		writeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetLimits(w http.ResponseWriter, r *http.Request) {
	var clientId model.ClientIdRequest
	err := json.NewDecoder(r.Body).Decode(&clientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	r.Body.Close()

	if err := clientId.ValidateClientIdRequest(); err != nil {
		writeJSONError(w, err)
		return
	}

	clientLimits, err := h.limitsManager.GetClientLimits(r.Context(), clientId.ClientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clientLimits)
	if err != nil {
		http.Error(w, "can not encode json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	var clientLimits model.ClientLimits
	err := json.NewDecoder(r.Body).Decode(&clientLimits)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	r.Body.Close()

	err = clientLimits.ValidateClientLimits()
	if err != nil {
		writeJSONError(w, err)
		return
	}

	err = h.limitsManager.UpdateClientLimits(r.Context(), clientLimits)
	if err != nil {
		writeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteLimits(w http.ResponseWriter, r *http.Request) {
	var clientId model.ClientIdRequest
	err := json.NewDecoder(r.Body).Decode(&clientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	r.Body.Close()

	if err := clientId.ValidateClientIdRequest(); err != nil {
		writeJSONError(w, err)
		return
	}

	err = h.limitsManager.DeleteClientLimits(r.Context(), clientId.ClientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
