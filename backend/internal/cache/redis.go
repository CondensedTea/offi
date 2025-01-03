package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
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

	return &Redis{client: client}, nil
}

func (r Redis) GetAllKeys(hashKey string) ([]string, error) {
	var keys []string

	switch hashKey {
	case logsKey, matchKey:
		break
	default:
		return nil, fmt.Errorf("unknown hash key: %s", hashKey)
	}

	res, err := r.client.HGetAll(hashKey).Result()
	if err != nil {
		return nil, err
	}
	for key := range res {
		keys = append(keys, key)
	}
	return keys, nil
}
