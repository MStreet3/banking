package errs

import "net/http"

type AppError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func NewInternalServerError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}

func UnexpectedDatabaseError() *AppError {
	return NewInternalServerError("Unexpected database error")
}

func NewValidationError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
	}
}

func InsufficientFundsError() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "insufficient funds",
	}
}

func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}
