package repo

import (
	"context"
	"fmt"
	"time"

	"encore.app/internal/pkg/errs"

	"github.com/redis/go-redis/v9"
)

func (r *UserRepo) StoreOTP(ctx context.Context, phone, otp string, ttl time.Duration) error {
	key := fmt.Sprintf("user:otp:%s", phone)
	if err := r.redis.Set(ctx, key, otp, ttl).Err(); err != nil {
		return errs.Internal(err, "failed to store OTP in redis")
	}
	return nil
}

func (r *UserRepo) GetOTP(ctx context.Context, phone string) (string, error) {
	key := fmt.Sprintf("user:otp:%s", phone)
	otp, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errs.BadRequest("OTP expired or not found")
	}
	if err != nil {
		return "", errs.Internal(err, "failed to get OTP from redis")
	}
	return otp, nil
}

func (r *UserRepo) DeleteOTP(ctx context.Context, phone string) error {
	key := fmt.Sprintf("user:otp:%s", phone)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return errs.Internal(err, "failed to delete OTP from redis")
	}
	return nil
}
