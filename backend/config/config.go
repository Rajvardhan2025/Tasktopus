package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server    ServerConfig
	MongoDB   MongoDBConfig
	JWT       JWTConfig
	CORS      CORSConfig
	WebSocket WebSocketConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type CORSConfig struct {
	Origins []string
}

type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
}

func Load() *Config {
	expiry := getEnvDuration("JWT_EXPIRY", 24*time.Hour)

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "project_management"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Expiry: expiry,
		},
		CORS: CORSConfig{
			Origins: getEnvStringSlice("CORS_ORIGINS"),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  getEnvInt("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize: getEnvInt("WS_WRITE_BUFFER_SIZE", 1024),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func getEnvStringSlice(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	if len(origins) == 0 {
		return nil
	}

	return origins
}
