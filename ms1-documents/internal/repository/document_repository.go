package repository

import (
	"context"
	"ms1-documents/internal/config"
	"ms1-documents/internal/domain"
	"ms1-documents/internal/utils"
	"ms1-documents/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type documentRepository struct {
	db *config.Database
}

func NewDocumentRepository(db *config.Database) DocumentRepository {
	return &documentRepository{
		db: db,
	}
}

func (r *documentRepository) Crear(contexto context.Context, documento *domain.Document) error {
	_, err := r.db.Collection.InsertOne(contexto, documento)
	if err != nil {
		if dupErr := utils.HandleMongoDuplicateError(err, documento.IDDocumento); dupErr != nil {
			return dupErr
		}
		return errors.ErrorInterno("Error al crear documento en la base de datos")
	}
	return nil
}

func (r *documentRepository) BuscarTodos(contexto context.Context) ([]domain.Document, error) {
	cursor, err := r.db.Collection.Find(contexto, bson.M{})
	if err != nil {
		return nil, errors.ErrorInterno("Error al obtener documentos de la base de datos")
	}
	defer cursor.Close(contexto)

	var documentos []domain.Document
	if err = cursor.All(contexto, &documentos); err != nil {
		return nil, errors.ErrorInterno("Error al decodificar documentos")
	}

	if documentos == nil {
		documentos = []domain.Document{}
	}

	return documentos, nil
}

func (r *documentRepository) BuscarPorID(contexto context.Context, id string) (*domain.Document, error) {
	var documento domain.Document
	err := r.db.Collection.FindOne(contexto, utils.DocumentIDFilter(id)).Decode(&documento)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.DocumentNotFoundError(id)
		}
		return nil, errors.ErrorInterno("Error al buscar documento en la base de datos")
	}
	return &documento, nil
}

func (r *documentRepository) Actualizar(contexto context.Context, id string, documento *domain.Document) error {
	result, err := r.db.Collection.UpdateOne(
		contexto,
		utils.DocumentIDFilter(id),
		bson.M{"$set": documento},
	)

	if err != nil {
		if dupErr := utils.HandleMongoDuplicateError(err, documento.IDDocumento); dupErr != nil {
			return dupErr
		}
		return errors.ErrorInterno("Error al actualizar documento en la base de datos")
	}

	if result.MatchedCount == 0 {
		return utils.DocumentNotFoundError(id)
	}

	return nil
}

func (r *documentRepository) Eliminar(contexto context.Context, id string) error {
	result, err := r.db.Collection.DeleteOne(contexto, utils.DocumentIDFilter(id))
	if err != nil {
		return errors.ErrorInterno("Error al eliminar documento de la base de datos")
	}

	if result.DeletedCount == 0 {
		return utils.DocumentNotFoundError(id)
	}

	return nil
}
