package service

import (
	"context"
	"database/sql"

	"encore.app/internal/features/user/booking/repo"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rp *repo.BRepo
}

func NewService(rediss *redis.Client, pg *sql.DB, pr sarama.SyncProducer) *Service {
	return &Service{
		rp: repo.NewRepo(
			rediss,
			pg,
			pr,
		),
	}
}

func (s *Service) FindNearbyCaptains(ctx context.Context, lat, long float64, km int) ([]map[string]string, error) {
	return s.rp.FindNearbyCaptains(ctx, lat, long, km)
}
