package config

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
    Ctx         context.Context
    Port        string
    MongoURI    string
    MongoDBName string
    Client      *mongo.Client
    RedisAddr    string // Redis адресі
    RedisPassword string // Redis паролі
}

func Load() *Config {
    port := os.Getenv("PORT")
    if port == "" {
        port = "50053"
    }
    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }
    dbName := os.Getenv("MONGO_DB")
    if dbName == "" {
        dbName = "inventory_db"
    }
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379" // Redis әдепкі адресі
    }

    redisPassword := os.Getenv("REDIS_PASSWORD")
    if redisPassword == "" {
        redisPassword = "" // Redis әдепкі пароль жоқ
    }

    // Контекст для подключения
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("❌ MongoDB connect error: %v", err)
    }

    return &Config{
        Ctx:         context.Background(),
        Port:        port,
        MongoURI:    uri,
        MongoDBName: dbName,
        Client:      client,
        RedisAddr:     redisAddr,
        RedisPassword: redisPassword,
    }
}
