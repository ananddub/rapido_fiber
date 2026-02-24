package booking_common

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (r *Repo) CreateBooking(ctx context.Context, bookingId string, data BookingData) error {
	err := r.LockCaptain(ctx, string(data.CaptainId))
	if err != nil {
		return err
	}
	err = r.LockUser(ctx, string(data.UserId))
	if err != nil {
		r.UnlockCaptain(ctx, string(data.CaptainId))
		return err
	}
	_, err = r.redis.HSet(ctx, r.GetBookingKey(bookingId),
		"user_id", data.UserId,
		"captain_id", data.CaptainId,
		"user_id", data.UserId,
		"status", data.Status,
		"pickup_location", data.PickupLocation,
		"drop_location", data.DropLocation,
		"paid_price", data.PaidPrice,
		"actual_price", data.ActualPrice,
		"payment_type", data.PaymentType,
		"accepted", data.Accepted,
		"is_canceled", data.IsCanceled,
		"is_verified", data.IsVerified,
		"is_success", data.IsSuccess,
		"started_at", data.StartedAt,
		"completed_at", data.CompletedAt,
		"created_at", time.Now().Unix(),
	).Result()
	return err
}

func (r *Repo) GetBooking(ctx context.Context, bookingId string) (*BookingData, error) {
	result, err := r.redis.HGetAll(ctx, r.GetBookingKey(bookingId)).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	captainId, err := strconv.ParseInt(result["captain_id"], 10, 32)
	if err != nil {
		return nil, err
	}
	userId, err := strconv.ParseInt(result["user_id"], 10, 32)
	if err != nil {
		return nil, err
	}

	bookingData := &BookingData{
		BookingId:      uuid.MustParse(bookingId),
		CaptainId:      int32(captainId),
		UserId:         int32(userId),
		Status:         result["status"],
		PickupLocation: result["pickup_location"],
		DropLocation:   result["drop_location"],
		PaidPrice:      result["paid_price"],
		ActualPrice:    result["actual_price"],
		PaymentType:    result["payment_type"],
		IsCanceled:     result["is_canceled"] == "1",
		IsVerified:     result["is_verified"] == "1",
		IsSuccess:      result["is_success"] == "1",
		Accepted:       result["accepted"] == "1",
	}
	return bookingData, nil
}

func (r *Repo) UpdateBookingStatus(ctx context.Context, bookingId string, status string) error {
	_, err := r.redis.HSet(ctx, r.GetBookingKey(bookingId),
		"status", status,
		"updated_at", time.Now().Unix(),
	).Result()
	return err
}

func (r *Repo) CancelBooking(ctx context.Context, bookingId string) error {
	booking, err := r.GetBooking(ctx, bookingId)
	if err != nil {
		return err
	}
	if booking == nil {
		return nil
	}
	err = r.UnlockCaptain(ctx, string(booking.CaptainId))
	if err != nil {
		return err
	}
	err = r.UnlockUser(ctx, string(booking.UserId))
	if err != nil {
		return err
	}
	return r.UpdateBookingStatus(ctx, bookingId, "canceled")
}
