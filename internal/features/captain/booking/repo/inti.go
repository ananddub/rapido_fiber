package repo

import (
	"encore.app/gen/pgdb"
	booking_common "encore.app/internal/common/booking/redis"
	"github.com/IBM/sarama"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
)

type BRepo struct {
	redis_client *redis.Client
	booking      *booking_common.Repo
	producer     sarama.SyncProducer
	mqtt         mqtt.Client
	pg           *pgdb.Queries
}

func NewRepo(rediss *redis.Client, mqtt_client mqtt.Client, pg *pgdb.Queries, pr sarama.SyncProducer) *BRepo {
	return &BRepo{
		redis_client: rediss,
		booking:      booking_common.NewRepo(rediss),
		producer:     pr,
		mqtt:         mqtt_client,
		pg:           pg,
	}
}
