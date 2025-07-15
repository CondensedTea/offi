package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r Client) SaveMatchExists(ctx context.Context, matchID int) error {
	return r.client.Set(ctx, strconv.Itoa(matchID), 1, 7*24*time.Hour).Err()
}

func (r Client) MatchExists(ctx context.Context, matchID int) (bool, error) {
	res, err := r.client.Get(ctx, strconv.Itoa(matchID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	return strconv.ParseBool(res)
}
