package config

import "os"

type Config struct {
	MongoURI      string
	MongoDatabase string
	GRPCPort      string
}

func Load() *Config {
	return &Config{
		MongoURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGODB_DATABASE", "order_db"),
		GRPCPort:      getEnv("GRPC_PORT", ":50052"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
