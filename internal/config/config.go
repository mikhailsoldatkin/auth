package config

import (
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "./.env"

// DB represents the configuration for the database.
type DB struct {
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host             string `env:"DB_HOST" env-required:"true"`
	Port             int    `env:"DB_PORT" env-required:"true"`
	PostgresDSN      string `env:"-"`
}

// GRPC represents the configuration for the gRPC server.
type GRPC struct {
	Port    int    `env:"GRPC_PORT" env-required:"true"`
	Host    string `env:"GRPC_HOST" env-required:"true"`
	Address string `env:"-"`
}

// HTTP represents the configuration for the http server.
type HTTP struct {
	Port    int    `env:"HTTP_PORT" env-required:"true"`
	Host    string `env:"HTTP_HOST" env-required:"true"`
	Address string `env:"-"`
}

// Swagger represents the configuration for the swagger server.
type Swagger struct {
	Port    int    `env:"SWAGGER_PORT" env-required:"true"`
	Host    string `env:"SWAGGER_HOST" env-required:"true"`
	Address string `env:"-"`
}

// Redis represents configuration for Redis.
type Redis struct {
	Host        string `env:"REDIS_HOST" env-required:"true"`
	Port        int    `env:"REDIS_PORT" env-required:"true"`
	ConnTimeout int    `env:"REDIS_CONNECTION_TIMEOUT_SEC" env-required:"true"`
	MaxIdle     int    `env:"REDIS_MAX_IDLE" env-required:"true"`
	MaxActive   int    `env:"REDIS_MAX_ACTIVE" env-required:"true"`
	IdleTimeout int    `env:"REDIS_IDLE_TIMEOUT_SEC" env-required:"true"`
	Address     string `env:"-"`
}

// KafkaConsumer represents configuration for KafkaConsumer.
type KafkaConsumer struct {
	Brokers []string `env:"KAFKA_BROKERS" env-required:"true"`
	GroupID string   `env:"KAFKA_GROUP_ID" env-required:"true"`
	Topic   string   `env:"KAFKA_TOPIC" env-required:"true"`
	Config  *sarama.Config
}

// Auth represents configuration for authentication.
type Auth struct {
	TokenSecretKey            string `env:"TOKEN_SECRET_KEY" env-required:"true"`
	RefreshTokenExpirationMin int    `env:"REFRESH_TOKEN_EXPIRATION_MIN" env-required:"true"`
	AccessTokenExpirationMin  int    `env:"ACCESS_TOKEN_EXPIRATION_MIN" env-required:"true"`
}

// Logger represents configuration for logger.
type Logger struct {
	Level      string `env:"LOG_LEVEL" env-required:"true"`
	Filename   string `env:"LOG_FILENAME" env-required:"true"`
	MaxSizeMB  int    `env:"LOG_MAX_SIZE_MB" env-required:"true"`
	MaxBackups int    `env:"LOG_MAX_BACKUPS" env-required:"true"`
	MaxAgeDays int    `env:"LOG_MAX_AGE_DAYS" env-required:"true"`
}

// Config represents the overall application configuration.
type Config struct {
	DB            DB
	GRPC          GRPC
	Redis         Redis
	HTTP          HTTP
	Swagger       Swagger
	KafkaConsumer KafkaConsumer
	Auth          Auth
	Logger        Logger
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

	cfg.DB.PostgresDSN = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.PostgresDB,
		cfg.DB.PostgresUser,
		cfg.DB.PostgresPassword,
	)

	cfg.Redis.Address = fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	cfg.GRPC.Address = fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	cfg.HTTP.Address = fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	cfg.Swagger.Address = fmt.Sprintf("%s:%d", cfg.Swagger.Host, cfg.Swagger.Port)

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V3_6_0_0
	kafkaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.KafkaConsumer.Config = kafkaConfig

	return &cfg, nil
}
