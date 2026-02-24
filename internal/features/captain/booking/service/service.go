package service

import (
	"context"

	booking_common "encore.app/internal/common/booking/redis"
	"encore.app/internal/features/captain/booking/repo"
)

type Service struct {
	repo *repo.BRepo
}

func NewService(rp *repo.BRepo) *Service {
	return &Service{
		repo: rp,
	}
}

func (s *Service) CreateCaptain(ctx context.Context, user booking_common.CaptainData) error {
	return s.repo.CreateCaptain(ctx, user)
}

func (s *Service) FindNearbyUser(ctx context.Context, lat, long float64, km int) ([]map[string]string, error) {
	return s.repo.FindNearbyUser(ctx, lat, long, km)
}

func (s *Service) GetCaptain(ctx context.Context, userId string) (*booking_common.CaptainData, error) {
	return s.repo.GetCaptain(ctx, userId)
}

func (s *Service) CreateBooking(ctx context.Context, booking booking_common.BookingData) error {

	err := s.repo.CreateBooking(ctx, booking)
	if err != nil {
		return err
	}

	err = s.repo.PublishBookingEvent(ctx, "booking_confirm", booking)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetBooking(ctx context.Context, bookingId int32) ([]booking_common.BookingData, error) {
	data, err := s.repo.GetBookingByCaptainId(ctx, bookingId)
	if err != nil {
		return nil, err
	}
	var bookingData []booking_common.BookingData = make([]booking_common.BookingData, len(data))
	for _, v := range data {
		bookingData = append(bookingData, booking_common.BookingData{
			BookingId:      v.ID,
			CaptainId:      v.CaptainID,
			UserId:         v.UserID,
			ActualPrice:    v.ActualPrice,
			PaidPrice:      v.PaidPrice,
			Accepted:       v.IsVerified.Bool,
			PickupLocation: v.PickupLocation,
			DropLocation:   v.DropLocation,
			IsCanceled:     v.IsCancelled.Bool,
			IsVerified:     v.IsVerified.Bool,
			CanceledBy:     v.CancelledBy.String,
			IsSuccess:      v.IsSuccessful.Bool,
			PaymentType:    v.PaymentMethod.String,
			Status:         string(v.Status.BookingStatus),
			CreatedAt:      v.CreatedAt.Time.Unix(),
		})
	}
	return bookingData, nil
}

func (s *Service) GetCurrentBooking(ctx context.Context, captain_id int32) (*booking_common.BookingData, error) {
	v, err := s.repo.GetCurrentBookingByCaptainId(ctx, captain_id)
	if err != nil {
		return nil, err
	}
	return &booking_common.BookingData{
		BookingId:      v.ID,
		CaptainId:      v.CaptainID,
		UserId:         v.UserID,
		ActualPrice:    v.ActualPrice,
		PaidPrice:      v.PaidPrice,
		Accepted:       v.IsVerified.Bool,
		PickupLocation: v.PickupLocation,
		DropLocation:   v.DropLocation,
		IsCanceled:     v.IsCancelled.Bool,
		IsVerified:     v.IsVerified.Bool,
		CanceledBy:     v.CancelledBy.String,
		IsSuccess:      v.IsSuccessful.Bool,
		PaymentType:    v.PaymentMethod.String,
		Status:         string(v.Status.BookingStatus),
		CreatedAt:      v.CreatedAt.Time.Unix(),
	}, nil
}

func (s *Service) CancelBooking(ctx context.Context, bookingId string, canceledBy string) error {
	err := s.repo.CancelBooking(ctx, bookingId, canceledBy)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PaymentSuccess(ctx context.Context, bookingId string) error {
	err := s.repo.PaymentSucess(ctx, bookingId)
	if err != nil {
		return err
	}
	return nil
}
