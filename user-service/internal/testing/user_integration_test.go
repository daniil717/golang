package testing

import (
    "context"
    "testing"
    "time"
    "user-service/internal/model"
    "user-service/internal/repository"
    "user-service/internal/usecase"

    "github.com/stretchr/testify/assert"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

var client *mongo.Client
var userRepo repository.UserRepository

func TestMain(m *testing.M) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("MongoDB-ге қосылу сәтсіз аяқталды: %v", err)
    }

    db := client.Database("user_service_test")

    userRepo = repository.NewMongoUserRepository(db.Collection("users"))

    m.Run()

    if err := client.Disconnect(ctx); err != nil {
        log.Fatalf("MongoDB-ден ажырату сәтсіз аяқталды: %v", err)
    }
}
func cleanupCollection(t *testing.T) {
    ctx := context.Background()
    err := userRepo.Cleanup(ctx)
    if err != nil {
        t.Fatalf("Failed to clean up collection: %v", err)
    }
}


func TestIntegration(t *testing.T) {
    cleanupCollection(t) // Тазалау — тестті таза бастау үшін

    uc := usecase.NewUserUsecase(userRepo)

    user := &model.User{
        Username: "integration_user",
        Email:    "integration@example.com",
        Password: "password123",
    }

    id, err := uc.CreateUser(context.Background(), user)
    assert.NoError(t, err)
    assert.NotEmpty(t, id)

    savedUser, err := uc.GetUserByID(context.Background(), id)
    assert.NoError(t, err)
    assert.Equal(t, user.Username, savedUser.Username)
    assert.Equal(t, user.Email, savedUser.Email)

    cleanupCollection(t) 
}
