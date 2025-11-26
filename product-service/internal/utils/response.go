package utils

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func Err(w http.ResponseWriter, appErr *AppError) {
	JSON(w, appErr.Status, map[string]any{
		"error": map[string]string{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}
