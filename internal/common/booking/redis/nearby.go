package booking_common

import (
	"context"
	"fmt"
)

func (r *Repo) FindNearbyCaptains(ctx context.Context, lat, long float64, radiusKm int) ([]map[string]string, error) {
	result, err := r.redis.Do(
		ctx,
		"FT.SEARCH", "captainidx",
		fmt.Sprintf("@location:[%f %f %d km] @is_available:{1} @is_online:{1}", long, lat, radiusKm),
	).Result()

	if err != nil {
		return nil, err
	}

	captains, _ := parseFTSearchResult(result)
	return captains, nil
}

func (r *Repo) FindNearbyUsers(ctx context.Context, lat, long float64, radiusKm int) ([]map[string]string, error) {
	result, err := r.redis.Do(
		ctx,
		"FT.SEARCH", "useridx",
		fmt.Sprintf("@location:[%f %f %d km] @is_available:{1} @is_online:{1}", long, lat, radiusKm),
	).Result()

	if err != nil {
		return nil, err
	}

	users, _ := parseFTSearchResult(result)
	return users, nil
}
