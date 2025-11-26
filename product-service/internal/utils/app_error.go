package utils

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string { return e.Message }

func BadRequest(msg string) *AppError {
	return &AppError{Code: "BAD_REQUEST", Message: msg, Status: http.StatusBadRequest}
}

func Validation(msg string) *AppError {
	return &AppError{Code: "VALIDATION_ERROR", Message: msg, Status: http.StatusUnprocessableEntity}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: msg, Status: http.StatusNotFound}
}

func Internal(msg string) *AppError {
	return &AppError{Code: "INTERNAL_ERROR", Message: msg, Status: http.StatusInternalServerError}
}
