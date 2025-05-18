package testing

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"order-service/internal/model"
	"order-service/internal/repository"
	"order-service/internal/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testRepo    repository.OrderRepository
	orderUc     *usecase.OrderUsecase
	mongoClient *mongo.Client
)

func TestMain(m *testing.M) {
	// MongoDB-ге қосылу
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("MongoDB-ге қосыла алмады: %v", err)
	}

	mongoClient = client
	db := client.Database("test_orders_db")
	coll := db.Collection("orders")

	// Репозиторийді бастау
	testRepo = repository.NewMongoOrderRepository(coll)

	// Паблишердің орнына nil береміз (publish тексермейміз)
	orderUc = usecase.NewOrderUsecase(testRepo, nil)

	// Тесттерді іске қосу
	code := m.Run()

	// Тест бітсе — мәліметтер базасын тазалау
	_ = coll.Drop(context.Background())
	_ = client.Disconnect(context.Background())

	os.Exit(code)
}

func TestIntegrationOrder(t *testing.T) {
	ctx := context.Background()

	// 📦 Жаңа тапсырыс
	order := &model.Order{
		UserID: "user123",
		Products: []model.Product{
			{ProductID: "507f1f77bcf86cd799439011", Quantity: 2},
		},
	}

	id, err := orderUc.CreateOrder(ctx, order)
	if err != nil {
		t.Fatalf("Тапсырысты құру сәтсіз: %v", err)
	}
	if id == "" {
		t.Fatal("Құрылған тапсырыстың ID бос болмауы тиіс")
	}

	// 🔍 Тапсырысты ID арқылы іздеу
	foundOrder, err := orderUc.GetOrder(ctx, id)
	if err != nil {
		t.Fatalf("Тапсырысты алу сәтсіз: %v", err)
	}
	if foundOrder == nil || foundOrder.UserID != "user123" {
		t.Fatalf("Күтілген тапсырыс табылмады, тапты: %+v", foundOrder)
	}
}
