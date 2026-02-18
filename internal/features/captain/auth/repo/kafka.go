package repo

import (
	"context"
	"encoding/json"

	"encore.app/internal/pkg/errs"

	"github.com/IBM/sarama"
)

func (r *CaptainRepo) SendSMS(ctx context.Context, phone, message string) error {
	payload := map[string]string{
		"phone": phone,
		"otp":   message,
	}

	msgBytes, err := json.Marshal(payload)
	if err != nil {
		return errs.Internal(err, "failed to marshal SMS payload")
	}

	_, _, err = r.kafka.SendMessage(&sarama.ProducerMessage{
		Topic: "sms-notification",
		Key:   sarama.StringEncoder("sms"),
		Value: sarama.ByteEncoder(msgBytes),
	})
	if err != nil {
		return errs.Internal(err, "failed to send SMS to kafka")
	}

	return nil
}
