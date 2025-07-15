package redis

import (
	"context"
	"strconv"
	"time"
)

func (r Client) SaveMatchExists(ctx context.Context, matchID int) error {
	return r.client.Set(ctx, strconv.Itoa(matchID), 1, 7*24*time.Hour).Err()
}

func (r Client) MatchExists(ctx context.Context, matchID int) (bool, error) {
	return r.client.Get(ctx, strconv.Itoa(matchID)).Bool()
}
