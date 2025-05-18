package migration

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func MigrateUp(db *mongo.Database) error {
	_, err := db.Collection("products").UpdateMany(
		context.Background(),
		bson.M{"quantity": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"quantity": 0}},
	)
	return err
}

func MigrateDown(db *mongo.Database) error {
	_, err := db.Collection("products").UpdateMany(
		context.Background(),
		bson.M{},
		bson.M{"$unset": bson.M{"quantity": ""}},
	)
	return err
}
