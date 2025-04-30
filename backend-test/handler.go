package main

import (
	"balancer/internal/model"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	resp := &TestRespose{
		ServerAdd: hostPort,
		Method:    r.Method,
		Uri:       r.RequestURI,
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil && err != io.EOF {
		writeJSONError(w, http.StatusBadRequest, "cannot read body")
		return
	}

	if buf.Len() == 0 {
		resp.Data = json.RawMessage(`null`)
	} else {
		if err := json.Unmarshal(buf.Bytes(), &resp.Data); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "cannot encode json", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	response := model.ErrorResponse{
		Message: message,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to write JSONE: %s", err)
	}
}
