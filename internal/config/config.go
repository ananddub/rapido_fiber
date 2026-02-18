package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// PostgreSQL
	PostgresHost     string `json:"postgres_host"`
	PostgresUser     string `json:"postgres_user"`
	PostgresPassword string `json:"postgres_password"`
	PostgresDB       string `json:"postgres_db"`
	PostgresPort     int    `json:"postgres_port"`

	// Redis
	RedisHost string `json:"redis_host"`
	RedisPort int    `json:"redis_port"`

	// Kafka (Redpanda)
	KafkaHost    string   `json:"kafka_host"`
	KafkaPort    int      `json:"kafka_port"`
	KafkaBrokers []string `json:"kafka_brokers"`

	// MinIO
	MinioEndpoint string `json:"minio_endpoint"`
	MinioUser     string `json:"minio_user"`
	MinioPassword string `json:"minio_password"`
	MinioUseSSL   bool   `json:"minio_use_ssl"`
	MinioRegion   string `json:"minio_region"`

	// Temporal
	TemporalHost string `json:"temporal_host"`
	TemporalPort int    `json:"temporal_port"`

	// JWT
	JWTSecret string `json:"jwt_secret"`
}

var cfg *Config

func InitConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	cfg = &Config{}

	cfg.PostgresHost = getEnv("POSTGRES_HOST", "localhost")
	cfg.PostgresUser = getEnv("POSTGRES_USER", "postgres")
	cfg.PostgresPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	cfg.PostgresDB = getEnv("POSTGRES_DB", "rapido")
	cfg.PostgresPort = getEnvAsInt("POSTGRES_PORT", 5432)

	cfg.RedisHost = getEnv("REDIS_HOST", "localhost")
	cfg.RedisPort = getEnvAsInt("REDIS_PORT", 6379)

	cfg.KafkaHost = getEnv("KAFKA_HOST", "localhost")
	cfg.KafkaPort = getEnvAsInt("KAFKA_PORT", 9092)
	kafkaBroker := getEnv("KAFKA_BROKER", fmt.Sprintf("%s:%d", cfg.KafkaHost, cfg.KafkaPort))
	cfg.KafkaBrokers = []string{kafkaBroker}

	minioHost := getEnv("MINIO_HOST", "localhost")
	minioPort := getEnvAsInt("MINIO_API_PORT", 9000)
	cfg.MinioEndpoint = fmt.Sprintf("http://%s:%d", minioHost, minioPort)
	cfg.MinioUser = getEnv("MINIO_ROOT_USER", "admin")
	cfg.MinioPassword = getEnv("MINIO_ROOT_PASSWORD", "password")
	cfg.MinioUseSSL = false
	cfg.MinioRegion = "us-east-1"

	cfg.TemporalHost = getEnv("TEMPORAL_HOST", "localhost")
	cfg.TemporalPort = getEnvAsInt("TEMPORAL_PORT", 7233)

	cfg.JWTSecret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}
