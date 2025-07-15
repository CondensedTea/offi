package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r Client) SaveMatchExists(ctx context.Context, matchID int) error {
	err := r.client.Set(ctx, strconv.Itoa(matchID), 1, 7*24*time.Hour).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (r Client) MatchExists(ctx context.Context, matchID int) (bool, error) {
	return r.client.Get(ctx, strconv.Itoa(matchID)).Bool()
}
