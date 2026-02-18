package booking_common

import (
	"context"
	"fmt"
	"time"
)

func (r *Repo) CreateUser(ctx context.Context, data UserData) error {
	key := r.GetUserKey(data.Id)
	isBooked := "0"
	isAvailable := "1"
	if data.IsBooked {
		isBooked = "1"
		isAvailable = "0"
	}

	_, err := r.redis.HSet(ctx, key,
		"name", data.Name,
		"phone", data.Phone,
		"status", "available",
		"location", fmt.Sprintf("%f,%f", data.Longitude, data.Latitude), // longitude, latitude
		"is_booked", isBooked,
		"is_available", isAvailable,
		"is_online", "1",
		"booking_id", "",
		"captain_id", "",
		"updated_at", time.Now().Unix(),
	).Result()

	return err
}

func (r *Repo) GetUser(ctx context.Context, captainId string) (*UserData, error) {
	key := r.GetUserKey(captainId)
	result, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("captain not found")
	}
	lat, long := parseLocation(result["location"])
	userData := &UserData{
		Id:        captainId,
		Name:      result["name"],
		Phone:     result["phone"],
		IsBooked:  result["is_booked"] == "1",
		Latitude:  lat,
		Longitude: long,
	}
	return userData, nil
}

func (r *Repo) UpdateUserLocation(ctx context.Context, userId string, lat, long float64) error {
	key := r.GetUserKey(userId)
	_, err := r.redis.HSet(ctx, key,
		"location", fmt.Sprintf("%f,%f", long, lat), // longitude first
		"updated_at", time.Now().Unix(),
	).Result()
	return err
}

func (r *Repo) UpdateUserStatus(ctx context.Context, userId string, status string) error {
	key := r.GetUserKey(userId)
	_, err := r.redis.HSet(ctx, key,
		"status", status,
		"updated_at", time.Now().Unix(),
	).Result()
	return err
}

func (r *Repo) LockUser(ctx context.Context, userId string) error {
	key := r.GetUserKey(userId)
	isBooked, err := r.redis.SetNX(ctx, key, "booking", time.Second*30).Result()
	if err != nil {
		return fmt.Errorf("failed to check user booking status: %v", err)
	}
	if isBooked {
		return fmt.Errorf("user is already booked")
	}
	return nil
}

func (r *Repo) UnlockUser(ctx context.Context, userId string) error {
	key := r.GetUserKey(userId)
	_, err := r.redis.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to unlock user: %v", err)
	}
	return nil
}
