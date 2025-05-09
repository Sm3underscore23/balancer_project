package api

import (
	"encoding/json"
	"net/http"

	"balancer/internal/model"
)

func (h *Handler) CreateLimits(w http.ResponseWriter, r *http.Request) {
	var clientId model.ClientLimits
	err := json.NewDecoder(r.Body).Decode(&clientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	r.Body.Close()

	err = h.limitsManager.CreateClientLimits(r.Context(), clientId) // убрать указатель
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
}

func (h *Handler) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	var clientId model.ClientLimits
	err := json.NewDecoder(r.Body).Decode(&clientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}
	r.Body.Close()

	err = h.limitsManager.UpdateClientLimits(r.Context(), clientId)
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

	err = h.limitsManager.DeleteClientLimits(r.Context(), clientId.ClientId)
	if err != nil {
		writeJSONError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
