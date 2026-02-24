package booking_background

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Config struct {
	subscriber sarama.ConsumerGroup
}

func NewConfig(subscriber sarama.ConsumerGroup) *Config {
	return &Config{
		subscriber: subscriber,
	}
}

func (c *Config) Subscribe(ctx context.Context) error {
	topics := []string{"booking_confirm"}
	handler := &BookingConfirmHandler{}
	for {
		err := c.subscriber.Consume(ctx, topics, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
