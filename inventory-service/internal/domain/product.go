package domain

type Product struct {
	ID       string  `bson:"_id,omitempty"`
	Name     string  `bson:"name"`
	Category string  `bson:"category"`
	Price    float64 `bson:"price"`
	Stock    int32   `bson:"stock"`
}
