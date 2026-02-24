package booking_background

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"encore.app/gen/pgdb"
	booking_common "encore.app/internal/common/booking/redis"
	"encore.app/internal/pkg/mqtt/topic"
	"github.com/IBM/sarama"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type BookingConfirmHandler struct {
	redis   *redis.Client
	mqtt    mqtt.Client
	booking *booking_common.Repo
	topic   *topic.Topic
	pg      *sql.DB
}

func NewBookingConfirmNew(redis *redis.Client, pg *sql.DB, mqtt mqtt.Client) *BookingConfirmHandler {
	return &BookingConfirmHandler{
		redis:   redis,
		pg:      pg,
		mqtt:    mqtt,
		topic:   topic.New(mqtt),
		booking: booking_common.NewRepo(redis),
	}
}

func (h *BookingConfirmHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BookingConfirmHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BookingConfirmHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {

		var payload booking_common.BookingData

		err := json.Unmarshal(message.Value, &payload)

		if err != nil {
			h.topic.PublishBookingError(payload)
			log.Printf("Error unmarshalling message: %v", err)
			session.MarkMessage(message, "")
			continue
		}
		err = h.Worker(session, payload)
		session.MarkMessage(message, "")
		if err != nil {
			h.topic.PublishBookingError(payload)
			log.Printf("Error processing booking: %v", err)
			continue
		}
		err = h.topic.PublishBookingConfirm(payload)
		if err != nil {
			log.Printf("Error publishing booking confirm: %v", err)
			continue
		}

	}
	return nil
}

func (h *BookingConfirmHandler) Worker(session sarama.ConsumerGroupSession, payload booking_common.BookingData) error {

	tx, err := h.pg.BeginTx(session.Context(), nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}

	q := pgdb.New(h.pg).WithTx(tx)
	err = q.LockUser(session.Context(), payload.UserId)
	_, err = q.GetCurrentBookingByUserId(session.Context(), payload.UserId)
	if err == nil {
		tx.Rollback()
		log.Printf("User already has an active booking: %v")
		return fmt.Errorf("user already has an active booking")
	}
	_, err = q.GetCurrentBookingByCaptainId(session.Context(), payload.CaptainId)
	if err == nil {
		tx.Rollback()
		log.Printf("Captain already has an active booking: %v")
		return fmt.Errorf("captain already has an active booking")
	}
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
	err = q.CreateBooking(session.Context(), pgdb.CreateBookingParams{
		ID:             payload.BookingId,
		UserID:         payload.UserId,
		CaptainID:      payload.CaptainId,
		PickupLocation: payload.PickupLocation,
		DropLocation:   payload.DropLocation,
		Status: pgdb.NullBookingStatus{
			Valid:         true,
			BookingStatus: pgdb.BookingStatusPENDING,
		},
	})

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
				UserStatus: pgdb.UserStatusBUSY,
			},
			CurrentBookingID: uuid.NullUUID{
				Valid: true,
				UUID:  payload.BookingId,
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
				UserStatus: pgdb.UserStatusBUSY,
			},
			CurrentBookingID: uuid.NullUUID{
				Valid: true,
				UUID:  payload.BookingId,
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

	otp := OptGenerator()
	h.redis.Set(session.Context(), fmt.Sprintf("booking:%s", payload.BookingId), otp, 0)
	return nil
}

func OptGenerator() string {
	otp := rand.Int31n(900000) + 100000
	return string(otp)
}
