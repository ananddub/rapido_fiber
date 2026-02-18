package auth_service

import (
	"context"
	"time"

	"encore.app/gen/pgdb"
	"encore.app/internal/config"
)

type Repository interface {
	GetByPhone(ctx context.Context, phone string) (*pgdb.Captain, error)
	Create(ctx context.Context, name, phone string) (*pgdb.Captain, error)

	StoreOTP(ctx context.Context, phone, otp string, ttl time.Duration) error
	GetOTP(ctx context.Context, phone string) (string, error)
	DeleteOTP(ctx context.Context, phone string) error

	SendSMS(ctx context.Context, phone, message string) error

	StartVerificationWorkflow(ctx context.Context, captainID int32) (string, error)
}

type AuthService struct {
	repo Repository
	cfg  *config.Config
}

func NewAuthService(repo Repository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}
