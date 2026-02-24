package booking_background

import (
	"database/sql"
	"encoding/json"
	"log"

	"encore.app/gen/pgdb"
	booking_common "encore.app/internal/common/booking/redis"
	"encore.app/internal/pkg/mqtt/topic"
	"github.com/IBM/sarama"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PaymentHandler struct {
	redis   *redis.Client
	mqtt    mqtt.Client
	topic   *topic.Topic
	booking *booking_common.Repo
	pg      *sql.DB
}

func PaymentConfirmmNew(redis *redis.Client, pg *sql.DB, mqtt mqtt.Client) *PaymentHandler {

	return &PaymentHandler{
		redis:   redis,
		pg:      pg,
		mqtt:    mqtt,
		topic:   topic.New(mqtt),
		booking: booking_common.NewRepo(redis),
	}
}
func (h *PaymentHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PaymentHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PaymentHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		var payload booking_common.BookingData

		err := json.Unmarshal(message.Value, &payload)

		if err != nil {
			h.topic.PublishBookingError(payload)
			log.Printf("Error unmarshalling message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		err = h.Worker(session, payload, false)
		if err != nil {
			h.topic.PublishBookingError(payload)
			log.Printf("Error processing booking: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (h *PaymentHandler) Worker(session sarama.ConsumerGroupSession, payload booking_common.BookingData, isError bool) error {

	tx, err := h.pg.BeginTx(session.Context(), nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}

	q := pgdb.New(h.pg).WithTx(tx)
	err = q.LockUser(session.Context(), payload.UserId)

	if err != nil {
		tx.Rollback()
		log.Printf("Error locking user: %v", err)
		return err
	}
	err = q.LockCaptain(session.Context(), payload.CaptainId)
	if err != nil {
		tx.Rollback()
		log.Printf("Error locking captain: %v", err)
		return err
	}
	booking, err := q.UpdateBookingSucess(session.Context(), payload.BookingId)
	if err != nil {
		tx.Rollback()
		log.Printf("Error processing booking: %v", err)
		return err
	}

	err = q.UpdateCaptainStatus(
		session.Context(),
		pgdb.UpdateCaptainStatusParams{
			ID: payload.CaptainId,
			Status: pgdb.NullUserStatus{
				Valid:      true,
				UserStatus: pgdb.UserStatusAVAILABLE,
			},
			CurrentBookingID: uuid.NullUUID{
				Valid: true,
				UUID:  uuid.Nil,
			},
		},
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error updating captain status: %v", err)
		return err
	}

	err = q.UpdateUserStatus(session.Context(),
		pgdb.UpdateUserStatusParams{
			ID: payload.UserId,
			Status: pgdb.NullUserStatus{
				Valid:      true,
				UserStatus: pgdb.UserStatusAVAILABLE,
			},
			CurrentBookingID: uuid.NullUUID{
				Valid: true,
				UUID:  uuid.Nil,
			},
		},
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error updating user status: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	h.booking.UpdateCaptainStatus(
		session.Context(),
		string(booking.CaptainID),
		"available",
	)

	h.booking.UpdateBookingStatus(
		session.Context(),
		string(booking.UserID),
		"available",
	)
	return nil
}
