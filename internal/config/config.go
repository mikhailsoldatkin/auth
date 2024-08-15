package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "./.env"

// DatabaseConfig represents the configuration for the database.
type DatabaseConfig struct {
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host             string `env:"DB_HOST" env-required:"true"`
	Port             int    `env:"DB_PORT" env-required:"true"`
	PostgresDSN      string `env:"-"`
}

// GRPCConfig represents the configuration for the gRPC server.
type GRPCConfig struct {
	Port int `env:"GRPC_PORT" env-required:"true"`
}

// RedisConfig represents configuration for Redis.
type RedisConfig struct {
	Host        string `env:"REDIS_HOST" env-required:"true"`
	Port        int    `env:"REDIS_PORT" env-required:"true"`
	ConnTimeout int    `env:"REDIS_CONNECTION_TIMEOUT_SEC" env-required:"true"`
	MaxIdle     int    `env:"REDIS_MAX_IDLE" env-required:"true"`
	MaxActive   int    `env:"REDIS_MAX_ACTIVE" env-required:"true"`
	IdleTimeout int    `env:"REDIS_IDLE_TIMEOUT_SEC" env-required:"true"`
	Address     string `env:"-"`
}

// Config represents the overall application configuration.
type Config struct {
	AppName  string `env:"APP_NAME" env-required:"true"`
	Database DatabaseConfig
	GRPC     GRPCConfig
	Redis    RedisConfig
}

// Load reads configuration from .env file.
func Load() (*Config, error) {
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file does not exist in project's root")
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config from environment variables: %w", err)
	}

	cfg.Database.PostgresDSN = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.PostgresDB,
		cfg.Database.PostgresUser,
		cfg.Database.PostgresPassword,
	)

	cfg.Redis.Address = fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	return &cfg, nil
}
