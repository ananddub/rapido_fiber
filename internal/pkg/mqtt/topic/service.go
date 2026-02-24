package topic

import (
	"encoding/json"
	"fmt"

	booking_common "encore.app/internal/common/booking/redis"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Topic struct {
	mqtt mqtt.Client
}

func New(mqtt mqtt.Client) *Topic {
	return &Topic{
		mqtt: mqtt,
	}
}

func (t *Topic) PublishBookingConfirm(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_CONFIRM, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishBookingError(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_ERROR, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishBookingCancel(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_CANCEL, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishBookingComplete(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_COMPLETE, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishBookingStart(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_START, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishBookingCancelByCaptain(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_CANCEL_BY_CAPTAIN, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishBookingCancelByUser(payload booking_common.BookingData) error {
	booking_id := payload.BookingId.String()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(BOOKING_CANCEL_BY_USER, booking_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishNotificatoinToUser(user_id int32, payload Notification) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(USER_NOTIFICATION, user_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishNotificatoinToCaptain(captain_id int32, payload Notification) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(fmt.Sprintf(CAPTAIN_NOTIFICATION, captain_id), 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (t *Topic) PublishNotificatoin(payload Notification) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := t.mqtt.Publish(NOTIFICATOIN, 1, false, data)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
