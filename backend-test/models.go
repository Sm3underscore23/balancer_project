package main

import "encoding/json"

type ErrorResponse struct {
	Message string `json:"errors"`
}

type TestResponse struct {
	ServerAdd string          `json:"server_address"`
	Method    string          `json:"http_method"`
	Uri       string          `json:"request_uri"`
	Data      json.RawMessage `json:"your_message"`
}
