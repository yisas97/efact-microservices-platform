package errors

import (
	"net/http"
	"testing"
)

func TestNuevoErrorValidacion(t *testing.T) {
	err := NuevoErrorValidacion("invalid input")

	if err.Code != http.StatusBadRequest {
		t.Errorf("Expected code %d, got %d", http.StatusBadRequest, err.Code)
	}

	if err.ErrorType != "Bad Request" {
		t.Errorf("Expected error type 'Bad Request', got '%s'", err.ErrorType)
	}

	if err.Message != "invalid input" {
		t.Errorf("Expected message 'invalid input', got '%s'", err.Message)
	}
}

func TestNuevoErrorNoEncontrado(t *testing.T) {
	err := NuevoErrorNoEncontrado("resource not found")

	if err.Code != http.StatusNotFound {
		t.Errorf("Expected code %d, got %d", http.StatusNotFound, err.Code)
	}

	if err.ErrorType != "Not Found" {
		t.Errorf("Expected error type 'Not Found', got '%s'", err.ErrorType)
	}

	if err.Message != "resource not found" {
		t.Errorf("Expected message 'resource not found', got '%s'", err.Message)
	}
}

func TestNuevoErrorConflicto(t *testing.T) {
	err := NuevoErrorConflicto("resource already exists")

	if err.Code != http.StatusConflict {
		t.Errorf("Expected code %d, got %d", http.StatusConflict, err.Code)
	}

	if err.ErrorType != "Conflict" {
		t.Errorf("Expected error type 'Conflict', got '%s'", err.ErrorType)
	}

	if err.Message != "resource already exists" {
		t.Errorf("Expected message 'resource already exists', got '%s'", err.Message)
	}
}

func TestNuevoErrorServidorInterno(t *testing.T) {
	err := NuevoErrorServidorInterno("something went wrong")

	if err.Code != http.StatusInternalServerError {
		t.Errorf("Expected code %d, got %d", http.StatusInternalServerError, err.Code)
	}

	if err.ErrorType != "Internal Server Error" {
		t.Errorf("Expected error type 'Internal Server Error', got '%s'", err.ErrorType)
	}

	if err.Message != "something went wrong" {
		t.Errorf("Expected message 'something went wrong', got '%s'", err.Message)
	}
}

func TestAppError_Error(t *testing.T) {
	err := &AppError{
		Code:      http.StatusBadRequest,
		ErrorType: "Bad Request",
		Message:   "test error",
	}

	expected := "test error"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestAppError_ToJSON(t *testing.T) {
	err := &AppError{
		Code:      http.StatusNotFound,
		ErrorType: "Not Found",
		Message:   "not found",
	}

	json := err.AJson()

	if json["status"] != http.StatusNotFound {
		t.Errorf("Expected status %d, got %v", http.StatusNotFound, json["status"])
	}

	if json["error"] != "Not Found" {
		t.Errorf("Expected error 'Not Found', got '%v'", json["error"])
	}

	if json["message"] != "not found" {
		t.Errorf("Expected message 'not found', got '%v'", json["message"])
	}
}

func TestAppError_String(t *testing.T) {
	err := &AppError{
		Code:      http.StatusConflict,
		ErrorType: "Conflict",
		Message:   "duplicate entry",
	}

	expected := "[409] Conflict: duplicate entry"
	if err.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.String())
	}
}

func TestAppError_AllStatusCodes(t *testing.T) {
	testCases := []struct {
		name         string
		constructor  func(string) *AppError
		expectedCode int
		expectedType string
	}{
		{
			name:         "ValidationError",
			constructor:  NuevoErrorValidacion,
			expectedCode: http.StatusBadRequest,
			expectedType: "Bad Request",
		},
		{
			name:         "NotFoundError",
			constructor:  NuevoErrorNoEncontrado,
			expectedCode: http.StatusNotFound,
			expectedType: "Not Found",
		},
		{
			name:         "ConflictError",
			constructor:  NuevoErrorConflicto,
			expectedCode: http.StatusConflict,
			expectedType: "Conflict",
		},
		{
			name:         "InternalServerError",
			constructor:  NuevoErrorServidorInterno,
			expectedCode: http.StatusInternalServerError,
			expectedType: "Internal Server Error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.constructor("test message")

			if err.Code != tc.expectedCode {
				t.Errorf("Expected code %d, got %d", tc.expectedCode, err.Code)
			}

			if err.ErrorType != tc.expectedType {
				t.Errorf("Expected type '%s', got '%s'", tc.expectedType, err.ErrorType)
			}

			if err.Message != "test message" {
				t.Errorf("Expected message 'test message', got '%s'", err.Message)
			}
		})
	}
}

func TestAppError_ErrorInterface(t *testing.T) {
	var err error = NuevoErrorValidacion("test")

	if err.Error() != "test" {
		t.Errorf("AppError does not implement error interface correctly")
	}
}

func TestAppError_ToJSON_HasAllFields(t *testing.T) {
	err := NuevoErrorServidorInterno("critical error")
	json := err.AJson()

	requiredFields := []string{"status", "error", "message"}
	for _, field := range requiredFields {
		if _, exists := json[field]; !exists {
			t.Errorf("ToJSON() missing required field: %s", field)
		}
	}

	if len(json) != len(requiredFields) {
		t.Errorf("ToJSON() should have exactly %d fields, got %d", len(requiredFields), len(json))
	}
}
