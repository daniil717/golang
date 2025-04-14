package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig содержит настройки подключения к MongoDB.
type MongoConfig struct {
	URI      string
	Database string
}

// NewMongoClient создает и возвращает клиент MongoDB.
func NewMongoClient(ctx context.Context, cfg MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.New("failed to ping MongoDB: " + err.Error())
	}

	return client, nil
}

// NewMongoDatabase возвращает экземпляр базы данных MongoDB.
func NewMongoDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}