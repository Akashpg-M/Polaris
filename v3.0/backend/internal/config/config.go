package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig
	Server ServerConfig
	Redis  RedisConfig
	DB     DBConfig
}

type AppConfig struct {
	Env      string
	LogLevel string
}

type ServerConfig struct {
	GatewayPort string
	EnginePort  string
}

type RedisConfig struct {
	URL string
}

type DBConfig struct {
	URL string
}

func Load() *Config {
	// Don't crash if .env is missing (e.g., in Docker). Just log it.
	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found. Relying on OS environment variables.")
	}

	return &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Server: ServerConfig{
			GatewayPort: getEnv("GATEWAY_PORT", "6080"),
			EnginePort:  getEnv("ENGINE_PORT", "6081"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
		},
		DB: DBConfig{
			URL: getEnv("POSTGRES_URL", "postgres://polaris_user:polaris_password@localhost:5432/polaris_core?sslmode=disable"),
		},
	}
}

// getEnv returns the environment variable or a safe fallback
func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return fallback
}