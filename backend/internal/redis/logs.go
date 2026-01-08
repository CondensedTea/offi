package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-faster/errors"
	"github.com/redis/go-redis/v9"
)

func (r Client) SetLogError(ctx context.Context, matchID int, err error) error {
	const errorCacheDuration = 3 * time.Hour

	if !r.enableErrorCaching {
		return nil
	}

	match := fmt.Sprintf("error-match-%d", matchID)
	return r.client.Set(ctx, match, err.Error(), errorCacheDuration).Err()
}

func (r Client) CheckLogError(ctx context.Context, matchID int) error {
	if !r.enableErrorCaching {
		return nil
	}

	match := fmt.Sprintf("error-match-%d", matchID)
	val, err := r.client.Get(ctx, match).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	cachedErrorTTL, err := r.client.TTL(ctx, match).Result()
	if err != nil {
		slog.Error("ttl lookup didn't find any keys after cached error was found", "component", "redis")
		return nil
	}

	return fmt.Errorf("%w in %s: %v", ErrCached, cachedErrorTTL, val)
}
