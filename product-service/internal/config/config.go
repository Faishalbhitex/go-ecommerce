package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getenvStrict("DB_HOST"),
		DBPort:     getenvStrict("DB_PORT"),
		DBUser:     getenvStrict("DB_USER"),
		DBPassword: getenvStrict("DB_PASSWORD"),
		DBName:     getenvStrict("DB_NAME"),
		Port:       getenv("PORT", "8080"),
	}

	return cfg
}

func getenvStrict(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is required but not set", key)
	}
	return value
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
