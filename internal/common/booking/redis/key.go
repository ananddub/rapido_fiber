package booking_common

import (
	"fmt"
)

func (r *Repo) GetBookingKey(bookingId string) string {
	return fmt.Sprintf("booking:%s", bookingId)
}

func (r *Repo) GetCaptainKey(captainId string) string {
	return fmt.Sprintf("captain:%s", captainId)
}

func (r *Repo) GetUserKey(userId string) string {
	return fmt.Sprintf("user:%s", userId)
}
