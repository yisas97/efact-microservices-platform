package service

import (
	"context"
	"ms1-documents/internal/domain"
	"ms1-documents/internal/validator"
	"ms1-documents/pkg/errors"
	"testing"
)

type mockRepository struct {
	createFunc   func(ctx context.Context, doc *domain.Document) error
	findAllFunc  func(ctx context.Context) ([]domain.Document, error)
	findByIDFunc func(ctx context.Context, id string) (*domain.Document, error)
	updateFunc   func(ctx context.Context, id string, doc *domain.Document) error
	deleteFunc   func(ctx context.Context, id string) error
}

func (m *mockRepository) Crear(ctx context.Context, doc *domain.Document) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, doc)
	}
	return nil
}

func (m *mockRepository) BuscarTodos(ctx context.Context) ([]domain.Document, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return []domain.Document{}, nil
}

func (m *mockRepository) BuscarPorID(ctx context.Context, id string) (*domain.Document, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.NuevoErrorNoEncontrado("not found")
}

func (m *mockRepository) Actualizar(ctx context.Context, id string, doc *domain.Document) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, doc)
	}
	return nil
}

func (m *mockRepository) Eliminar(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

type mockPublisher struct {
	publishFunc func(documentID, uuid string) error
	called      bool
	lastDocID   string
	lastUUID    string
}

func (m *mockPublisher) PublishDocumentCreated(documentID, uuid string) error {
	m.called = true
	m.lastDocID = documentID
	m.lastUUID = uuid

	if m.publishFunc != nil {
		return m.publishFunc(documentID, uuid)
	}
	return nil
}

func TestCreateDocument_Success(t *testing.T) {
	repo := &mockRepository{}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
	}

	ctx := context.Background()
	err := svc.CrearDocumento(ctx, doc)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if doc.UUID == "" {
		t.Error("Expected UUID to be generated")
	}

	if !publisher.called {
		t.Error("Expected publisher to be called")
	}

	if publisher.lastDocID != doc.IDDocumento {
		t.Errorf("Expected docID %s, got %s", doc.IDDocumento, publisher.lastDocID)
	}

	if publisher.lastUUID != doc.UUID {
		t.Errorf("Expected UUID %s, got %s", doc.UUID, publisher.lastUUID)
	}
}

func TestCreateDocument_RepositoryError(t *testing.T) {
	repo := &mockRepository{
		createFunc: func(ctx context.Context, doc *domain.Document) error {
			return errors.NuevoErrorServidorInterno("database error")
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
	}

	ctx := context.Background()
	err := svc.CrearDocumento(ctx, doc)

	if err == nil {
		t.Error("Expected repository error")
	}

	if publisher.called {
		t.Error("Expected publisher NOT to be called on repository error")
	}
}

func TestCreateDocument_PublisherError(t *testing.T) {
	repo := &mockRepository{}
	publisher := &mockPublisher{
		publishFunc: func(documentID, uuid string) error {
			return errors.NuevoErrorServidorInterno("rabbitmq error")
		},
	}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
	}

	ctx := context.Background()
	err := svc.CrearDocumento(ctx, doc)

	if err == nil {
		t.Error("Expected publisher error")
	}

	appErr, ok := err.(*errors.AppError)
	if !ok {
		t.Error("Expected AppError")
	}

	if appErr.Message != "Error al publicar mensaje a RabbitMQ" {
		t.Errorf("Expected RabbitMQ error message, got: %s", appErr.Message)
	}
}

func TestGetAllDocuments_Success(t *testing.T) {
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

	repo := &mockRepository{
		findAllFunc: func(ctx context.Context) ([]domain.Document, error) {
			return expectedDocs, nil
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	ctx := context.Background()
	docs, err := svc.ObtenerTodosDocumentos(ctx)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(docs) != len(expectedDocs) {
		t.Errorf("Expected %d documents, got %d", len(expectedDocs), len(docs))
	}
}

func TestGetAllDocuments_Error(t *testing.T) {
	repo := &mockRepository{
		findAllFunc: func(ctx context.Context) ([]domain.Document, error) {
			return nil, errors.NuevoErrorServidorInterno("database error")
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	ctx := context.Background()
	_, err := svc.ObtenerTodosDocumentos(ctx)

	if err == nil {
		t.Error("Expected error")
	}
}

func TestGetDocumentByID_Success(t *testing.T) {
	expectedDoc := &domain.Document{
		IDDocumento: "FACT-123456789",
		UUID:        "uuid-1",
	}

	repo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*domain.Document, error) {
			if id == "FACT-123456789" {
				return expectedDoc, nil
			}
			return nil, errors.NuevoErrorNoEncontrado("not found")
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	ctx := context.Background()
	doc, err := svc.ObtenerDocumentoPorID(ctx, "FACT-123456789")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if doc.IDDocumento != expectedDoc.IDDocumento {
		t.Errorf("Expected ID %s, got %s", expectedDoc.IDDocumento, doc.IDDocumento)
	}
}

func TestGetDocumentByID_NotFound(t *testing.T) {
	repo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*domain.Document, error) {
			return nil, errors.NuevoErrorNoEncontrado("not found")
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	ctx := context.Background()
	_, err := svc.ObtenerDocumentoPorID(ctx, "NONEXISTENT")

	if err == nil {
		t.Error("Expected not found error")
	}
}

func TestUpdateDocument_Success(t *testing.T) {
	existingDoc := &domain.Document{
		IDDocumento: "FACT-123456789",
		UUID:        "existing-uuid",
	}

	repo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*domain.Document, error) {
			return existingDoc, nil
		},
		updateFunc: func(ctx context.Context, id string, doc *domain.Document) error {
			return nil
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	updatedDoc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 200.0,
		IgvTotal:               36.0,
		MontoTotal:             236.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 100.0,
				Cantidad:       2,
				PrecioTotal:    200.0,
				IgvTotal:       36.0,
			},
		},
	}

	ctx := context.Background()
	err := svc.ActualizarDocumento(ctx, "FACT-123456789", updatedDoc)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if updatedDoc.UUID != "existing-uuid" {
		t.Errorf("Expected UUID to be preserved as 'existing-uuid', got '%s'", updatedDoc.UUID)
	}
}

func TestUpdateDocument_NotFound(t *testing.T) {
	repo := &mockRepository{
		findByIDFunc: func(ctx context.Context, id string) (*domain.Document, error) {
			return nil, errors.NuevoErrorNoEncontrado("not found")
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	doc := &domain.Document{
		IDDocumento:            "FACT-123456789",
		RucEmisor:              "20123456789",
		RucReceptor:            "20987654321",
		MontoTotalSinImpuestos: 100.0,
		IgvTotal:               18.0,
		MontoTotal:             118.0,
		Items: []domain.Item{
			{
				Descripcion:    "Item 1",
				PrecioUnitario: 50.0,
				Cantidad:       2,
				PrecioTotal:    100.0,
				IgvTotal:       18.0,
			},
		},
	}

	ctx := context.Background()
	err := svc.ActualizarDocumento(ctx, "NONEXISTENT", doc)

	if err == nil {
		t.Error("Expected not found error")
	}
}

func TestDeleteDocument_Success(t *testing.T) {
	repo := &mockRepository{
		deleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	ctx := context.Background()
	err := svc.EliminarDocumento(ctx, "FACT-123456789")

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDeleteDocument_NotFound(t *testing.T) {
	repo := &mockRepository{
		deleteFunc: func(ctx context.Context, id string) error {
			return errors.NuevoErrorNoEncontrado("not found")
		},
	}
	publisher := &mockPublisher{}
	svc := NewDocumentService(repo, publisher, validator.NewDocumentValidator(), "amqp://guest:guest@localhost:5672/")

	ctx := context.Background()
	err := svc.EliminarDocumento(ctx, "NONEXISTENT")

	if err == nil {
		t.Error("Expected not found error")
	}
}

