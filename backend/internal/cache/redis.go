package cache

import (
	"errors"
	"os"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var ErrCached = errors.New("cached error, retry later")

const (
	LeagueRGL   = "rgl"
	LeagueETF2L = "etf2l"
)

type Redis struct {
	client *redis.Client

	enableErrorCaching bool
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

	_, disableErrCaching := os.LookupEnv("DISABLE_ERROR_CACHING")

	return &Redis{client: client, enableErrorCaching: !disableErrCaching}, nil
}
