package errors

import (
	"fmt"
	"net/http"
)

// AppError representa un error de aplicación
type AppError struct {
	Code      int    `json:"status" example:"400"`
	ErrorType string `json:"error" example:"Bad Request"`
	Message   string `json:"message" example:"Error en la solicitud"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) AJson() map[string]interface{} {
	return map[string]interface{}{
		"status":  e.Code,
		"error":   e.ErrorType,
		"message": e.Message,
	}
}

func ErrorValidacion(mensaje string) *AppError {
	return &AppError{
		Code:      http.StatusBadRequest,
		ErrorType: "Bad Request",
		Message:   mensaje,
	}
}

func ErrorNoEncontrado(mensaje string) *AppError {
	return &AppError{
		Code:      http.StatusNotFound,
		ErrorType: "Not Found",
		Message:   mensaje,
	}
}

func ErrorConflicto(mensaje string) *AppError {
	return &AppError{
		Code:      http.StatusConflict,
		ErrorType: "Conflict",
		Message:   mensaje,
	}
}

func ErrorInterno(mensaje string) *AppError {
	return &AppError{
		Code:      http.StatusInternalServerError,
		ErrorType: "Internal Server Error",
		Message:   mensaje,
	}
}

func (e *AppError) String() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.ErrorType, e.Message)
}
