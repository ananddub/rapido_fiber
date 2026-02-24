package repo

import (
	"context"
	"database/sql"

	"encore.app/gen/pgdb"
	booking_common "encore.app/internal/common/booking/redis"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type BRepo struct {
	redis_client *redis.Client
	booking      *booking_common.Repo
	db           *pgdb.Queries
	producer     sarama.SyncProducer
	pg           *sql.DB
}

func NewRepo(rediss *redis.Client, pg *sql.DB, pr sarama.SyncProducer) *BRepo {
	db := pgdb.New(pg)
	return &BRepo{
		redis_client: rediss,
		booking:      booking_common.NewRepo(rediss),
		producer:     pr,
		db:           db,
		pg:           pg,
	}
}

func (r *BRepo) FindNearbyCaptains(ctx context.Context, lat, long float64, km int) ([]map[string]string, error) {
	return r.booking.FindNearbyCaptains(ctx, lat, long, km)
}

func (r *BRepo) GetUser(ctx context.Context, userId string) (*booking_common.UserData, error) {
	return r.booking.GetUser(ctx, userId)
}

func (r *BRepo) UpdateUserLocation(ctx context.Context, userId string, lat, long float64) error {
	return r.booking.UpdateUserLocation(ctx, userId, lat, long)
}
