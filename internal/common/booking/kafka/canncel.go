package booking_background

import (
	"context"
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

type CancelHandler struct {
	redis   *redis.Client
	mqtt    mqtt.Client
	topic   *topic.Topic
	booking *booking_common.Repo
	pg      *sql.DB
}

func CancelNew(redis *redis.Client, pg *sql.DB, mqtt mqtt.Client) *CancelHandler {

	return &CancelHandler{
		redis:   redis,
		pg:      pg,
		mqtt:    mqtt,
		topic:   topic.New(mqtt),
		booking: booking_common.NewRepo(redis),
	}
}
func (h *CancelHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *CancelHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *CancelHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		var payload booking_common.BookingData

		err := json.Unmarshal(message.Value, &payload)
		if err != nil {
			h.topic.PublishBookingError(payload)
			log.Printf("Error unmarshalling message: %v", err)
			session.MarkMessage(message, "")
			continue
		}
		err = h.Worker(session.Context(), payload)

		session.MarkMessage(message, "")

		if err != nil {
			h.topic.PublishBookingError(payload)
			log.Printf("Error processing message: %v", err)
		}
		err = h.topic.PublishBookingCancel(payload)
		if err != nil {
			log.Printf("Error publishing booking cancel: %v", err)
		}

	}
	return nil
}

func (h *CancelHandler) Worker(ctx context.Context, payload booking_common.BookingData) error {
	tx, err := h.pg.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}

	q := pgdb.New(h.pg).WithTx(tx)
	err = q.LockUser(ctx, payload.UserId)

	if err != nil {
		tx.Rollback()
		log.Printf("Error locking user: %v", err)
		return err
	}

	err = q.LockCaptain(ctx, payload.CaptainId)

	if err != nil {
		h.topic.PublishBookingError(payload)
		tx.Rollback()
		log.Printf("Error locking captain: %v", err)
		return err
	}

	_, err = q.UpdateBookingCancel(ctx, pgdb.UpdateBookingCancelParams{
		CancelledBy: sql.NullString{
			String: payload.CanceledBy,
			Valid:  payload.CanceledBy != "",
		},
		ID: payload.BookingId,
	})

	if err != nil {
		tx.Rollback()
		log.Printf("Error processing booking: %v", err)
		return err
	}

	err = q.UpdateCaptainStatus(
		ctx,
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
		h.topic.PublishBookingError(payload)
		tx.Rollback()
		log.Printf("Error updating captain status: %v", err)
		return err
	}

	err = q.UpdateUserStatus(ctx,
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

	booking, err := q.GetBooking(ctx, payload.BookingId)
	if err != nil {
		log.Printf("Error fetching booking: %v", err)
		return err
	}

	h.booking.UpdateCaptainStatus(
		ctx,
		string(booking.CaptainID),
		string(pgdb.UserStatusAVAILABLE),
	)

	h.booking.UpdateBookingStatus(
		ctx,
		string(booking.UserID),
		string(pgdb.UserStatusAVAILABLE),
	)

	return nil

}
