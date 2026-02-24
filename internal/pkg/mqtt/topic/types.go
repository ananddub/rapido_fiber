package topic

const BOOKING_CONFIRM string = "/booking/%s/booking_confirm"
const BOOKING_START string = "/booking/%s/booking_start"
const BOOKING_ERROR string = "/booking/%s/booking_error"
const BOOKING_CANCEL string = "/booking/%s/booking_cancel"
const BOOKING_CONFIRM_PAYMENT string = "/booking/%s/booking_confirm_payment"
const BOOKING_COMPLETE string = "/booking/%s/booking_complete"
const BOOKING_CANCEL_BY_CAPTAIN string = "/booking/%s/booking_cancel_by_captain"
const BOOKING_CANCEL_BY_USER string = "/booking/%s/booking_cancel_by_user"

const NOTIFICATOIN string = "/notification"
const CAPTAIN_NOTIFICATION string = "/captain/%s/notification"
const USER_NOTIFICATION string = "/user/%s/notification"

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
