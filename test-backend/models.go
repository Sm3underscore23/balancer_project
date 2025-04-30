package main

import "encoding/json"

type ErrorResponse struct {
	Message string `json:"errors"`
}

type TestRespose struct {
	Method string          `json:"http_method"`
	Uri    string          `json:"request_uri"`
	Data   json.RawMessage `json:"your_message"`
}
