package redis

import (
	"context"
	"fmt"
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
	return fmt.Errorf("%w: %v", ErrCached, val)
}
