package config

import "os"

type Config struct {
	Port           string
	InventoryAddr  string
	OrderAddr      string
	JWTSecret      string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		InventoryAddr: getEnv("INVENTORY_SERVICE_ADDR", "localhost:50051"),
		OrderAddr:     getEnv("ORDER_SERVICE_ADDR", "localhost:50052"),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}