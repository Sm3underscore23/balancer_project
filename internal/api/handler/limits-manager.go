package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"balancer/internal/model"
	"balancer/pkg/logger"
)

// func logErrDecoder(ctx context.Context, err error) {
// 	ctx = logger.WithLogErr(ctx, err)
// 	slog.ErrorContext(ctx, "json decoder error")
// }

// func logErrValidate(ctx context.Context, err error) {
// 	ctx = logger.WithLogErr(ctx, err)
// 	slog.ErrorContext(ctx, "json validate error")
// }

func (h *Handler) CreateLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := logger.FromContext(ctx)
	l.Info("handler CreateLimits started")

	var clientLimits model.ClientLimits
	err := json.NewDecoder(r.Body).Decode(&clientLimits)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}
	r.Body.Close()

	if err = clientLimits.ValidateClientLimits(); err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	err = h.limitsManager.CreateClientLimits(r.Context(), clientLimits)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	ctx = logger.AddValuesToContext(ctx,
		logger.StatusCode, http.StatusCreated,
		slog.Group(logger.GroupClientLimits,
			logger.ClientID, clientLimits.ClientId,
			logger.ClinetTokenCapacity, clientLimits.Capacity,
			logger.ClinetRateRefill, clientLimits.RatePerSec,
		))
	logger.FromContext(ctx).Info("handler CreateLimits comlete")
}

func (h *Handler) GetLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := logger.FromContext(ctx)
	l.Info("handler GetLimits started")

	var clientId model.ClientIdRequest
	err := json.NewDecoder(r.Body).Decode(&clientId)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}
	r.Body.Close()

	if err := clientId.ValidateClientIdRequest(); err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	clientLimits, err := h.limitsManager.GetClientLimits(r.Context(), clientId.ClientId)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clientLimits)
	if err != nil {
		ctx := logger.AddValuesToContext(ctx, logger.Error, err)
		logger.FromContext(ctx).Info("json encode error")

		http.Error(w, "can not encode json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	ctx = logger.AddValuesToContext(ctx,
		logger.StatusCode, http.StatusOK,
		logger.ClientID, clientId,
	)
	logger.FromContext(ctx).Info("handler GetLimits comlete")
}

func (h *Handler) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := logger.FromContext(ctx)
	l.Info("handler UpdateLimits started")

	var clientLimits model.ClientLimits
	err := json.NewDecoder(r.Body).Decode(&clientLimits)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}
	r.Body.Close()

	err = clientLimits.ValidateClientLimits()
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	err = h.limitsManager.UpdateClientLimits(r.Context(), clientLimits)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx = logger.AddValuesToContext(ctx,
		logger.StatusCode, http.StatusCreated,
		slog.Group(logger.GroupClientLimits,
			logger.ClientID, clientLimits.ClientId,
			logger.ClinetTokenCapacity, clientLimits.Capacity,
			logger.ClinetRateRefill, clientLimits.RatePerSec,
		))
	logger.FromContext(ctx).Info("handler UpdateLimits comlete")
}

func (h *Handler) DeleteLimits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := logger.FromContext(ctx)
	l.Info("handler UpdateLimits started")

	var clientId model.ClientIdRequest
	err := json.NewDecoder(r.Body).Decode(&clientId)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}
	r.Body.Close()

	if err := clientId.ValidateClientIdRequest(); err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	err = h.limitsManager.DeleteClientLimits(r.Context(), clientId.ClientId)
	if err != nil {
		writeJSONError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	ctx = logger.AddValuesToContext(ctx,
		logger.StatusCode, http.StatusOK,
		logger.ClientID, clientId,
	)
	logger.FromContext(ctx).Info("handler UpdateLimits comlete")
}
