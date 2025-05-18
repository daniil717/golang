package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(uri string) (*mongo.Client, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	log.Println("✅ Connected to MongoDB (v1 driver)")
	return client, ctx, nil
}
