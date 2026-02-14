package service

import (
	"context"
	"ms1-documents/internal/domain"
)

type DocumentService interface {
	CrearDocumento(contexto context.Context, documento *domain.Document) error
	ObtenerTodosDocumentos(contexto context.Context) ([]domain.Document, error)
	ObtenerDocumentoPorID(contexto context.Context, id string) (*domain.Document, error)
	ActualizarDocumento(contexto context.Context, id string, documento *domain.Document) error
	EliminarDocumento(contexto context.Context, id string) error
	VerificarDocumento(contexto context.Context, documento *domain.Document, firma string) (bool, error)
}
