package repo

import (
	"context"

	booking_common "encore.app/internal/common/booking/redis"
)

func (r *BRepo) CreateCaptain(ctx context.Context, user booking_common.CaptainData) error {
	return r.booking.CreateCaptain(ctx, user)
}

func (r *BRepo) FindNearbyUser(ctx context.Context, lat, long float64, km int) ([]map[string]string, error) {
	return r.booking.FindNearbyUsers(ctx, lat, long, km)
}

func (r *BRepo) GetCaptain(ctx context.Context, userId string) (*booking_common.CaptainData, error) {
	return r.booking.GetCaptain(ctx, userId)
}

func (r *BRepo) UpdateCaptainLocation(ctx context.Context, captainId string, lat, long float64) error {
	return r.booking.UpdateCaptainLocation(ctx, captainId, lat, long)
}

func (r *BRepo) CreateBooking(ctx context.Context, booking booking_common.BookingData) error {
	err := r.booking.CreateBooking(ctx, booking.BookingId.String(), booking)
	if err != nil {
		return err
	}

	return nil
}

func (r *BRepo) GenereateKeyUserAndCaptainKey(userId, captainid string) string {
	return userId + ":" + captainid
}
