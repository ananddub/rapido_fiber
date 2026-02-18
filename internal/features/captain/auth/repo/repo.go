package repo

import (
	"errors"

	"encore.app/gen/pgdb"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"go.temporal.io/sdk/client"
)

var ErrNotFound = errors.New("captain not found")

type CaptainRepo struct {
	queries  *pgdb.Queries
	redis    *redis.Client
	kafka    sarama.SyncProducer
	temporal client.Client
}

func NewCaptainRepo(
	queries *pgdb.Queries,
	redis *redis.Client,
	kafka sarama.SyncProducer,
	temporal client.Client,
) *CaptainRepo {
	return &CaptainRepo{
		queries:  queries,
		redis:    redis,
		kafka:    kafka,
		temporal: temporal,
	}
}
