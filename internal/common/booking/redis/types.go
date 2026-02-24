package booking_common

import "github.com/google/uuid"

type CaptainData struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsBooked  bool    `json:"is_booked"`
}

type UserData struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Status    string  `json:"status"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsBooked  bool    `json:"is_booked"`
}

type BookingData struct {
	BookingId      uuid.UUID `json:"booking_id"`
	CaptainId      int32     `json:"captain_id"`
	UserId         int32     `json:"user_id"`
	ActualPrice    string    `json:"actual_price"`
	PaidPrice      string    `json:"paid_price"`
	Accepted       bool      `json:"accepted"`
	PickupLocation string    `json:"pickup_location"`
	DropLocation   string    `json:"drop_location"`
	IsCanceled     bool      `json:"is_canceled"`
	IsVerified     bool      `json:"is_verified"`
	CanceledBy     string    `json:"canceled_by,omitempty"` // "user", "captain", "system"
	IsSuccess      bool      `json:"is_success"`
	PaymentType    string    `json:"payment_type"` // "cash", "upi", "card"
	Status         string    `json:"status"`       // "pending", "accepted", "completed", "canceled"
	CreatedAt      int64     `json:"created_at"`
	StartedAt      int64     `json:"started_at,omitempty"`
	CompletedAt    int64     `json:"completed_at,omitempty"`
}

type NearbyCaptain struct {
	CaptainID string  `json:"captain_id"`
	Name      string  `json:"name"`
	Distance  float64 `json:"distance_km"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}
