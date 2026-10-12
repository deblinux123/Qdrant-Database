package config

import "os"

type Config struct {
	Port       string
	QdrantHost string
	QdrantPort string
	OllamaURL  string
	RedisAddr  string
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		QdrantHost: getEnv("QDRANT_HOST", "localhost"),
		QdrantPort: getEnv("QDRANT_PORT", "6334"),
		OllamaURL:  getEnv("OLLAMA_URL", "http://localhost:11434"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6378"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
