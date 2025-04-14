package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	MongoDatabase string
	GRPCPort      string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	return &Config{
		MongoURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGODB_DATABASE", "inventory"),
		GRPCPort:      getEnv("GRPC_PORT", ":50051"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
