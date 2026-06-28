package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getRequiredEnv("DATABASE_URL"),
		JWTSecret:   getRequiredEnv("JWT_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Fatal: Environment variable %s is required but not set", key)
	}

	return value
}
