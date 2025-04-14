package repository

import (
	"context"
	"order-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	ListByUser(ctx context.Context, userID string, limit, offset int32) ([]*domain.Order, error)
}

type OrderMongoRepository struct {
	collection *mongo.Collection
}

func NewOrderMongoRepository(db *mongo.Database) *OrderMongoRepository {
	return &OrderMongoRepository{
		collection: db.Collection("orders"),
	}
}

func (r *OrderMongoRepository) Create(ctx context.Context, order *domain.Order) error {
	order.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *OrderMongoRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	return &order, err
}

func (r *OrderMongoRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *OrderMongoRepository) ListByUser(ctx context.Context, userID string, limit, offset int32) ([]*domain.Order, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*domain.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
