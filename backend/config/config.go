package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
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
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGO_DATABASE", "project_management")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("WS_READ_BUFFER_SIZE", 1024)
	viper.SetDefault("WS_WRITE_BUFFER_SIZE", 1024)

	expiry, err := time.ParseDuration(viper.GetString("JWT_EXPIRY"))
	if err != nil {
		expiry = 24 * time.Hour
	}

	return &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Env:  viper.GetString("ENV"),
		},
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGO_URI"),
			Database: viper.GetString("MONGO_DATABASE"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			Expiry: expiry,
		},
		CORS: CORSConfig{
			Origins: viper.GetStringSlice("CORS_ORIGINS"),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  viper.GetInt("WS_READ_BUFFER_SIZE"),
			WriteBufferSize: viper.GetInt("WS_WRITE_BUFFER_SIZE"),
		},
	}
}
