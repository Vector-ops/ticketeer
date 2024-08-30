package config

import (
	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBUser     string `env:"DB_USER,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
}

func NewEnvConfig() *EnvConfig {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load env file %e", err)
	}

	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Failed to load env vars %e", err)
	}

	return config
}
