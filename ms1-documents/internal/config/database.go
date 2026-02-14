package config

import (
	"context"
	"log"
	"ms1-documents/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewDatabase(uri, dbName, collectionName string) (*Database, error) {
	ctx, cancel := utils.CrearContextoConTimeoutDB(context.Background())
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	db := &Database{
		Client:     client,
		Collection: collection,
	}

	err = db.createIndexes(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Conectado a MongoDB")
	return db, nil
}

func (db *Database) createIndexes(contexto context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    map[string]interface{}{"idDocumento": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := db.Collection.Indexes().CreateOne(contexto, indexModel)
	if err != nil {
		return err
	}

	log.Println("Indices creados correctamente")
	return nil
}

func (db *Database) Disconnect() {
	if db.Client != nil {
		ctx, cancel := utils.CrearContextoConTimeoutDB(context.Background())
		defer cancel()
		db.Client.Disconnect(ctx)
	}
}
