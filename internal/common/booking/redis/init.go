package booking_common

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	redis *redis.Client
}

func NewRepo(redis *redis.Client) *Repo {
	return &Repo{
		redis: redis,
	}
}

func (r *Repo) Init(ctx context.Context) error {
	if err := r.CreateSchemaFtQuery(ctx); err != nil {
		return fmt.Errorf("failed to create schema and indexes: %w", err)
	}
	return nil
}

func (r *Repo) createCaptainIndex(ctx context.Context) error {
	client := r.redis

	client.Do(ctx, "FT.DROPINDEX", "captainidx")

	_, err := client.Do(ctx,
		"FT.CREATE", "captainidx", "ON", "HASH", "PREFIX", "1", "captain:", "SCHEMA",
		"location", "GEO",
		"name", "TEXT",
		"phone", "TEXT",
		"status", "TAG",
		"booking_id", "TEXT",
		"user_id", "TEXT",
		"is_online", "BOOLEAN",
		"client_id", "TEXT",
		"ip_address", "TEXT",
		"updated_at", "NUMERIC", "SORTABLE",
	).Result()
	return err
}

func (r *Repo) createUserIndex(ctx context.Context) error {
	client := r.redis

	client.Do(ctx, "FT.DROPINDEX", "useridx")

	_, err := client.Do(ctx,
		"FT.CREATE", "useridx", "ON", "HASH", "PREFIX", "1", "user:", "SCHEMA",
		"location", "GEO",
		"name", "TEXT",
		"phone", "TEXT",
		"booking_id", "TEXT",
		"captain_id", "TEXT",
		"status", "TAG",
		"is_online", "TAG",
		"client_id", "TEXT",
		"ip_address", "TEXT",
		"updated_at", "NUMERIC", "SORTABLE",
	).Result()
	return err
}

func (r *Repo) CreateSchemaFtQuery(ctx context.Context) error {
	if err := r.createCaptainIndex(ctx); err != nil {
		return fmt.Errorf("failed to create captain index: %w", err)
	}
	if err := r.createUserIndex(ctx); err != nil {
		return fmt.Errorf("failed to create user index: %w", err)
	}
	return nil
}
