package booking_common

import (
	"context"
	"fmt"
	"time"
)

func (r *Repo) CreateUser(ctx context.Context, data UserData) error {
	key := r.GetUserKey(data.Id)

	fields := make([]interface{}, 0)

	if data.Name != "" {
		fields = append(fields, "name", data.Name)
	}
	if data.Phone != "" {
		fields = append(fields, "phone", data.Phone)
	}
	if data.Latitude != 0 || data.Longitude != 0 {
		fields = append(fields, "location", fmt.Sprintf("%f,%f", data.Longitude, data.Latitude))
	}
	if data.Status != "" {
		fields = append(fields, "status", data.Status)
	}
	isBooked := "0"
	isAvailable := "1"
	if data.IsBooked {
		isBooked = "1"
		isAvailable = "0"
	}
	fields = append(fields, "is_booked", isBooked)
	fields = append(fields, "is_available", isAvailable)
	fields = append(fields, "is_online", "1")
	fields = append(fields, "status", "available")
	fields = append(fields, "booking_id", "")
	fields = append(fields, "captain_id", "")
	fields = append(fields, "updated_at", time.Now().Unix())

	_, err := r.redis.HSet(ctx, key, fields...).Result()

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
