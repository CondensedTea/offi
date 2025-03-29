package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	matchKey = "matches"
	logsKey  = "logs"
)

const errorMatchExpire = 3 * time.Hour

var ErrCached = errors.New("cached error, retry later")

type Redis struct {
	client *redis.Client
}

func New(url string) (*Redis, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err = redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}

	if err = redisotel.InstrumentMetrics(client); err != nil {
		return nil, err
	}

	return &Redis{client: client}, nil
}

func (r Redis) GetAllKeys(ctx context.Context, hashKey string) ([]string, error) {
	var keys []string

	switch hashKey {
	case logsKey, matchKey:
		break
	default:
		return nil, fmt.Errorf("unknown hash key: %s", hashKey)
	}

	res, err := r.client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}
	for key := range res {
		keys = append(keys, key)
	}
	return keys, nil
}
