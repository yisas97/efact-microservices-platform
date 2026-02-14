package repository

import (
	"context"
	"ms1-documents/internal/domain"
)

type DocumentRepository interface {
	Crear(contexto context.Context, documento *domain.Document) error
	BuscarTodos(contexto context.Context) ([]domain.Document, error)
	BuscarPorID(contexto context.Context, id string) (*domain.Document, error)
	Actualizar(contexto context.Context, id string, documento *domain.Document) error
	Eliminar(contexto context.Context, id string) error
}
