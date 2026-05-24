package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorsResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, ErrorsResponse{
		Error: message,
	})
}
