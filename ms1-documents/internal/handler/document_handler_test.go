package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"ms1-documents/internal/domain"
	"ms1-documents/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockService struct {
	createDocumentFunc  func(ctx context.Context, doc *domain.Document) error
	getAllDocumentsFunc func(ctx context.Context) ([]domain.Document, error)
	getDocumentByIDFunc func(ctx context.Context, id string) (*domain.Document, error)
	updateDocumentFunc  func(ctx context.Context, id string, doc *domain.Document) error
	deleteDocumentFunc  func(ctx context.Context, id string) error
	verifyDocumentFunc  func(ctx context.Context, documento *domain.Document, firma string) (bool, error)
}

func (m *mockService) CrearDocumento(ctx context.Context, doc *domain.Document) error {
	if m.createDocumentFunc != nil {
		return m.createDocumentFunc(ctx, doc)
	}
	return nil
}

func (m *mockService) ObtenerTodosDocumentos(ctx context.Context) ([]domain.Document, error) {
	if m.getAllDocumentsFunc != nil {
		return m.getAllDocumentsFunc(ctx)
	}
	return []domain.Document{}, nil
}

func (m *mockService) ObtenerDocumentoPorID(ctx context.Context, id string) (*domain.Document, error) {
	if m.getDocumentByIDFunc != nil {
		return m.getDocumentByIDFunc(ctx, id)
	}
	return nil, errors.ErrorNoEncontrado("not found")
}

func (m *mockService) ActualizarDocumento(ctx context.Context, id string, doc *domain.Document) error {
	if m.updateDocumentFunc != nil {
		return m.updateDocumentFunc(ctx, id, doc)
	}
	return nil
}

func (m *mockService) EliminarDocumento(ctx context.Context, id string) error {
	if m.deleteDocumentFunc != nil {
		return m.deleteDocumentFunc(ctx, id)
	}
	return nil
}

func (m *mockService) VerificarDocumento(ctx context.Context, documento *domain.Document, firma string) (bool, error) {
	if m.verifyDocumentFunc != nil {
		return m.verifyDocumentFunc(ctx, documento, firma)
	}
	return false, errors.ErrorNoEncontrado("not found")
}

func setupRouter(handler *DocumentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/documents", handler.CrearDocumento)
	router.GET("/documents", handler.ObtenerDocumentos)
	router.GET("/documents/:id", handler.ObtenerDocumento)
	router.PUT("/documents/:id", handler.ActualizarDocumento)
	router.DELETE("/documents/:id", handler.EliminarDocumento)
	router.POST("/documents/verify", handler.VerificarDocumento)

	return router
}

func TestCreateDocument_Success(t *testing.T) {
	svc := &mockService{
		createDocumentFunc: func(ctx context.Context, doc *domain.Document) error {
			doc.UUID = "test-uuid"
			return nil
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"idDocumento":            "FACT-123456789",
		"rucEmisor":              "20123456789",
		"rucReceptor":            "20987654321",
		"montoTotalSinImpuestos": 100.0,
		"igvTotal":               18.0,
		"montoTotal":             118.0,
		"items": []map[string]interface{}{
			{
				"descripcion":    "Item 1",
				"precioUnitario": 50.0,
				"cantidad":       2,
				"precioTotal":    100.0,
				"igvTotal":       18.0,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response domain.Document
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.UUID != "test-uuid" {
		t.Errorf("Expected UUID 'test-uuid', got '%s'", response.UUID)
	}
}

func TestCreateDocument_InvalidJSON(t *testing.T) {
	svc := &mockService{}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateDocument_ValidationError(t *testing.T) {
	svc := &mockService{
		createDocumentFunc: func(ctx context.Context, doc *domain.Document) error {
			return errors.ErrorValidacion("invalid document")
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"idDocumento": "INVALID",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/documents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetDocuments_Success(t *testing.T) {
	expectedDocs := []domain.Document{
		{
			IDDocumento: "FACT-123456789",
			UUID:        "uuid-1",
		},
		{
			IDDocumento: "FACT-987654321",
			UUID:        "uuid-2",
		},
	}

	svc := &mockService{
		getAllDocumentsFunc: func(ctx context.Context) ([]domain.Document, error) {
			return expectedDocs, nil
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/documents", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []domain.Document
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != len(expectedDocs) {
		t.Errorf("Expected %d documents, got %d", len(expectedDocs), len(response))
	}
}

func TestGetDocuments_Error(t *testing.T) {
	svc := &mockService{
		getAllDocumentsFunc: func(ctx context.Context) ([]domain.Document, error) {
			return nil, errors.ErrorInterno("database error")
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/documents", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetDocument_Success(t *testing.T) {
	expectedDoc := &domain.Document{
		IDDocumento: "FACT-123456789",
		UUID:        "uuid-1",
	}

	svc := &mockService{
		getDocumentByIDFunc: func(ctx context.Context, id string) (*domain.Document, error) {
			if id == "FACT-123456789" {
				return expectedDoc, nil
			}
			return nil, errors.ErrorNoEncontrado("not found")
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/documents/FACT-123456789", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response domain.Document
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.IDDocumento != expectedDoc.IDDocumento {
		t.Errorf("Expected ID %s, got %s", expectedDoc.IDDocumento, response.IDDocumento)
	}
}

func TestGetDocument_NotFound(t *testing.T) {
	svc := &mockService{
		getDocumentByIDFunc: func(ctx context.Context, id string) (*domain.Document, error) {
			return nil, errors.ErrorNoEncontrado("not found")
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/documents/NONEXISTENT", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateDocument_Success(t *testing.T) {
	svc := &mockService{
		updateDocumentFunc: func(ctx context.Context, id string, doc *domain.Document) error {
			return nil
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"idDocumento":            "FACT-123456789",
		"rucEmisor":              "20123456789",
		"rucReceptor":            "20987654321",
		"montoTotalSinImpuestos": 200.0,
		"igvTotal":               36.0,
		"montoTotal":             236.0,
		"items": []map[string]interface{}{
			{
				"descripcion":    "Item 1",
				"precioUnitario": 100.0,
				"cantidad":       2,
				"precioTotal":    200.0,
				"igvTotal":       36.0,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/documents/FACT-123456789", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateDocument_InvalidJSON(t *testing.T) {
	svc := &mockService{}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("PUT", "/documents/FACT-123456789", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateDocument_NotFound(t *testing.T) {
	svc := &mockService{
		updateDocumentFunc: func(ctx context.Context, id string, doc *domain.Document) error {
			return errors.ErrorNoEncontrado("not found")
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"idDocumento":            "FACT-123456789",
		"rucEmisor":              "20123456789",
		"rucReceptor":            "20987654321",
		"montoTotalSinImpuestos": 200.0,
		"igvTotal":               36.0,
		"montoTotal":             236.0,
		"items": []map[string]interface{}{
			{
				"descripcion":    "Item 1",
				"precioUnitario": 100.0,
				"cantidad":       2,
				"precioTotal":    200.0,
				"igvTotal":       36.0,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/documents/NONEXISTENT", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteDocument_Success(t *testing.T) {
	svc := &mockService{
		deleteDocumentFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("DELETE", "/documents/FACT-123456789", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Documento eliminado correctamente" {
		t.Errorf("Expected success message, got: %s", response["message"])
	}
}

func TestVerifyDocument_Success_ValidUUID(t *testing.T) {
	svc := &mockService{
		verifyDocumentFunc: func(ctx context.Context, documento *domain.Document, firma string) (bool, error) {
			return true, nil
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"documento": map[string]interface{}{
			"idDocumento":            "FACT-123456789",
			"rucEmisor":              "20123456789",
			"rucReceptor":            "20987654321",
			"montoTotalSinImpuestos": 100.0,
			"igvTotal":               18.0,
			"montoTotal":             118.0,
			"items": []map[string]interface{}{
				{
					"descripcion":    "Item 1",
					"precioUnitario": 50.0,
					"cantidad":       2,
					"precioTotal":    100.0,
					"igvTotal":       18.0,
				},
			},
		},
		"firma": "firma-encriptada-valida",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/documents/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if !response.Valid {
		t.Error("Expected valid to be true")
	}

	if response.Message != "La firma es valida y el documento no ha sido modificado" {
		t.Errorf("Expected success message, got: %s", response.Message)
	}
}

func TestVerifyDocument_Success_InvalidUUID(t *testing.T) {
	svc := &mockService{
		verifyDocumentFunc: func(ctx context.Context, documento *domain.Document, firma string) (bool, error) {
			return false, nil
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"documento": map[string]interface{}{
			"idDocumento":            "FACT-123456789",
			"rucEmisor":              "20123456789",
			"rucReceptor":            "20987654321",
			"montoTotalSinImpuestos": 100.0,
			"igvTotal":               18.0,
			"montoTotal":             118.0,
			"items": []map[string]interface{}{
				{
					"descripcion":    "Item 1",
					"precioUnitario": 50.0,
					"cantidad":       2,
					"precioTotal":    100.0,
					"igvTotal":       18.0,
				},
			},
		},
		"firma": "firma-encriptada-invalida",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/documents/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Valid {
		t.Error("Expected valid to be false")
	}

	if response.Message != "La firma es invalida o el documento ha sido modificado" {
		t.Errorf("Expected error message, got: %s", response.Message)
	}
}

func TestVerifyDocument_InvalidJSON(t *testing.T) {
	svc := &mockService{}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	req, _ := http.NewRequest("POST", "/documents/verify", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestVerifyDocument_NotFound(t *testing.T) {
	svc := &mockService{
		verifyDocumentFunc: func(ctx context.Context, documento *domain.Document, firma string) (bool, error) {
			return false, errors.ErrorNoEncontrado("not found")
		},
	}
	handler := NewDocumentHandler(svc)
	router := setupRouter(handler)

	reqBody := map[string]interface{}{
		"documento": map[string]interface{}{
			"idDocumento":            "NONEXISTENT",
			"rucEmisor":              "20123456789",
			"rucReceptor":            "20987654321",
			"montoTotalSinImpuestos": 100.0,
			"igvTotal":               18.0,
			"montoTotal":             118.0,
			"items": []map[string]interface{}{
				{
					"descripcion":    "Item 1",
					"precioUnitario": 50.0,
					"cantidad":       2,
					"precioTotal":    100.0,
					"igvTotal":       18.0,
				},
			},
		},
		"firma": "any-firma",
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/documents/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
