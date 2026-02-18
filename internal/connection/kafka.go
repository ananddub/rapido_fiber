package connection

import (
	"fmt"

	"encore.app/internal/config"

	"github.com/IBM/sarama"
)

func NewKafkaProducer(cfg *config.Config) (sarama.SyncProducer, error) {
	kConfig := sarama.NewConfig()
	kConfig.Producer.Return.Successes = true
	kConfig.Producer.RequiredAcks = sarama.WaitForAll
	kConfig.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return producer, nil
}
