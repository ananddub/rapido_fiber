package repo

import (
	"context"

	"encore.app/gen/pgdb"
)

func (b *BRepo) GetBookingByCaptainId(ctx context.Context, captain_id int32) ([]pgdb.Booking, error) {
	data, err := b.pg.GetCaptainBookingsByCaptainId(ctx, captain_id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *BRepo) GetCurrentBookingByCaptainId(ctx context.Context, captain_id int32) (*pgdb.Booking, error) {
	data, err := b.pg.GetCurrentBookingByCaptainId(ctx, captain_id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
