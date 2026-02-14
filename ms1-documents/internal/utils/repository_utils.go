package utils

import (
	"fmt"
	"ms1-documents/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DocumentIDFilter(id string) bson.M {
	return bson.M{"idDocumento": id}
}

func HandleMongoDuplicateError(err error, idDocumento string) error {
	if mongo.IsDuplicateKeyError(err) {
		return errors.ErrorConflicto(fmt.Sprintf("Ya existe un documento con ID %s", idDocumento))
	}
	return nil
}

func DocumentNotFoundError(id string) error {
	return errors.ErrorNoEncontrado(fmt.Sprintf("Documento con ID %s no encontrado", id))
}
