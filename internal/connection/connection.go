package connection

import (
	"database/sql"
	"fmt"

	"encore.app/gen/pgdb"
	"encore.app/internal/config"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/client"
)

type Connections struct {
	DB            *sql.DB
	Redis         *redis.Client
	KafkaProducer sarama.SyncProducer
	S3Client      *s3.Client
	Temporal      client.Client
	Query         *pgdb.Queries
}

var connection *Connections

func InitConnection() (*Connections, error) {
	if connection != nil {
		return connection, nil
	}

	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	db, err := NewDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	redisClient, err := NewRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	kafkaProducer, err := NewKafkaProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %w", err)
	}

	s3Client, err := NewS3Client(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	temporalClient, err := NewTemporalClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to temporal: %w", err)
	}

	queries := pgdb.New(db)

	connection = &Connections{
		DB:            db,
		Redis:         redisClient,
		KafkaProducer: kafkaProducer,
		S3Client:      s3Client,
		Temporal:      temporalClient,
		Query:         queries,
	}

	return connection, nil
}
