package rest

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal server error: failed to encode response", http.StatusInternalServerError)
	}
}

func JSONSuccess(w http.ResponseWriter, status int, resp APIResponse) {
	writeJSON(w, status, resp)
}

func JSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIResponse{
		Message: message,
		Data:    nil,
	})
}

func JSONValidationError(w http.ResponseWriter, errors map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, APIResponse{
		Message: "The given data was invalid.",
		Errors:  errors,
	})
}
