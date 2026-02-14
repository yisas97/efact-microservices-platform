package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI        string
	MongoDB         string
	MongoCollection string
	RabbitMQURI     string
	Port            string
	LogDir          string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		MongoURI:        getEnv("MONGO_URI", "mongodb://admin:admin123@localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "efact_db"),
		MongoCollection: getEnv("MONGO_COLLECTION", "documents"),
		RabbitMQURI:     getEnv("RABBITMQ_URI", "amqp://admin:admin123@localhost:5672/"),
		Port:            getEnv("PORT", "3500"),
		LogDir:          getEnv("LOG_DIR", "./logs"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
