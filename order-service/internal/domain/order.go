package domain

import (
	"time"
)

type OrderItem struct {
	ProductID string  `bson:"product_id"`
	Quantity  int32   `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type Order struct {
	ID        string      `bson:"_id,omitempty"`
	UserID    string      `bson:"user_id"`
	Items     []OrderItem `bson:"items"`
	Total     float64     `bson:"total"`
	Status    string      `bson:"status"`
	CreatedAt time.Time   `bson:"created_at"`
}