package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
	RedisAddr   string
	RabbitMQURL string
}

func Load() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return Config{
		Port:        port,
		Env:         env,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   redisAddr,
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
	}
}
