package repo

import (
	"context"
	"encoding/json"
	"strconv"

	booking_common "encore.app/internal/common/booking/redis"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

func (r *BRepo) PublishBookingEvent(ctx context.Context, topic string, event booking_common.BookingData) error {
	str, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = r.producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(r.GenereateKeyUserAndCaptainKey(strconv.Itoa(int(event.UserId)), strconv.Itoa(int(event.CaptainId)))),
			Value: sarama.StringEncoder(str),
		},
	)
	return err
}

func (r *BRepo) PaymentSucess(ctx context.Context, bookingId string) error {
	id, err := uuid.Parse(bookingId)
	if err != nil {
		return err
	}
	str, err := json.Marshal(booking_common.BookingData{BookingId: id})
	if err != nil {
		return err
	}
	_, _, err = r.producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: "payment_success",
			Key:   sarama.StringEncoder(bookingId),
			Value: sarama.StringEncoder(str),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *BRepo) CancelBooking(ctx context.Context, bookingId string, canceledBy string) error {
	id, err := uuid.Parse(bookingId)
	if err != nil {
		return err
	}
	str, err := json.Marshal(booking_common.BookingData{BookingId: id, CanceledBy: canceledBy})
	if err != nil {
		return err
	}
	_, _, err = r.producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: "booking_cancel",
			Key:   sarama.StringEncoder(bookingId),
			Value: sarama.StringEncoder(str),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
