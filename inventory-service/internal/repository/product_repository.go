package repository

import (
	"context"
	"inventory-servicee/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, id string, product *domain.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int32) ([]*domain.Product, error)
}

type ProductMongoRepository struct {
	collection *mongo.Collection
}

func NewProductMongoRepository(db *mongo.Database) *ProductMongoRepository {
	return &ProductMongoRepository{
		collection: db.Collection("products"),
	}
}

func (r *ProductMongoRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *ProductMongoRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductMongoRepository) Update(ctx context.Context, id string, product *domain.Product) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": product},
	)
	return err
}

func (r *ProductMongoRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *ProductMongoRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Product, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{},
		options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
