package testing

import (
    "context"
    "log"
    "os"
    "testing"
    "time"

    "inventory-service/internal/model"
    "inventory-service/internal/redis"
    "inventory-service/internal/repository"
    "inventory-service/internal/usecase"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    testRepo   repository.ProductRepository
    productUc  *usecase.ProductUsecase
    mongoClient *mongo.Client
)

func TestMain(m *testing.M) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // MongoDB қосу
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("MongoDB қосыла алмады: %v", err)
    }
    mongoClient = client

    // Redis инициализациясы
    err = redis.InitRedisWithParams("localhost:6379", "", 0)
    if err != nil {
        log.Fatalf("Redis инициализациясы сәтсіз: %v", err)
    }

    db := client.Database("test_inventory_db")
    coll := db.Collection("products")

    // Репозиторий жасау
    testRepo = repository.NewMongoProductRepository(coll)
    productUc = usecase.NewProductUsecase(testRepo)

    // Тесттерді іске қосу
    code := m.Run()

    // Тазалау
    _ = coll.Drop(context.Background())
    _ = client.Disconnect(context.Background())

    os.Exit(code)
}

func TestIntegrationProductUsecase(t *testing.T) {
    ctx := context.Background()

    // 1. Жаңа өнім құру
    prod := &model.Product{
        Name:        "Test Product",
        Description: "Test Description",
        Category:    "TestCategory",
        Stock:       100,
        Price:       50.5,
    }

    id, err := productUc.CreateProduct(ctx, prod)
    if err != nil {
        t.Fatalf("Өнімді құру сәтсіз: %v", err)
    }

    if id == "" {
        t.Fatal("Құрылған өнімнің ID бос болмауы тиіс")
    }

    prod.ID = id

    // 2. Өнімді алу (Redis кэшсіз — тікелей репозиториден)
    productFromRepo, err := productUc.GetProduct(ctx, id)
    if err != nil {
        t.Fatalf("Өнімді алу сәтсіз: %v", err)
    }
    if productFromRepo == nil || productFromRepo.Name != prod.Name {
        t.Fatalf("Алынған өнім дұрыс емес: %+v", productFromRepo)
    }

    // 3. Өнімді Redis кэштен алу (екінші рет шақырғанда кэштен алынуы тиіс)
    productFromCache, err := productUc.GetProduct(ctx, id)
    if err != nil {
        t.Fatalf("Redis-тен өнімді алу сәтсіз: %v", err)
    }
    if productFromCache == nil || productFromCache.Name != prod.Name {
        t.Fatalf("Redis-тен алынған өнім дұрыс емес: %+v", productFromCache)
    }

    // 4. Өнімді жаңарту
    prod.Description = "Updated Description"
    prod.Stock = 90
    prod.Price = 45.0
    err = productUc.UpdateProduct(ctx, prod)
    if err != nil {
        t.Fatalf("Өнімді жаңарту сәтсіз: %v", err)
    }

    // 5. Жаңартылған өнімді қайтадан алу (кэш тазаланған соң, қайта репозиториден алынады)
    updatedProduct, err := productUc.GetProduct(ctx, id)
    if err != nil {
        t.Fatalf("Жаңартылған өнімді алу сәтсіз: %v", err)
    }
    if updatedProduct.Description != "Updated Description" || updatedProduct.Stock != 90 || updatedProduct.Price != 45.0 {
        t.Fatalf("Жаңартылған өнім дұрыс емес: %+v", updatedProduct)
    }

    // 6. Тізімін алу
    products, err := productUc.ListProducts(ctx, "TestCategory", 1, 10)
    if err != nil {
        t.Fatalf("Өнімдер тізімін алу сәтсіз: %v", err)
    }
    if len(products) == 0 {
        t.Fatal("Өнімдер тізімі бос")
    }

    // 7. Қойманы азайту
    err = productUc.DecreaseStock(ctx, id, 10)
    if err != nil {
        t.Fatalf("Қойманы азайту сәтсіз: %v", err)
    }

    decreasedProduct, err := productUc.GetProduct(ctx, id)
    if err != nil {
        t.Fatalf("Қойманы азайтудан кейін өнімді алу сәтсіз: %v", err)
    }
    if decreasedProduct.Stock != 80 {
        t.Fatalf("Қойма дұрыс азайтылмады, күтілген: 80, шыққан: %d", decreasedProduct.Stock)
    }

    // 8. Өнімді өшіру
    err = productUc.DeleteProduct(ctx, id)
    if err != nil {
        t.Fatalf("Өнімді өшіру сәтсіз: %v", err)
    }

    deletedProduct, err := productUc.GetProduct(ctx, id)
    if err == nil && deletedProduct != nil {
        t.Fatalf("Өшіру сәтті болған соң өнім әлі де табылады: %+v", deletedProduct)
    }
}
