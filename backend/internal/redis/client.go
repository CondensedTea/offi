package redis

import (
	"errors"
	"os"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var ErrCached = errors.New("cached error, retry later")

const (
	// LeagueRGL represents the rgl.gg league.
	LeagueRGL = "rgl"
	// LeagueETF2L represents the etf2l.org league.
	LeagueETF2L = "etf2l"
)

type Client struct {
	client *redis.Client

	enableErrorCaching bool
}

func NewClient(url string) (*Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	if err = redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}

	_, disableErrCaching := os.LookupEnv("DISABLE_ERROR_CACHING")

	return &Client{client: client, enableErrorCaching: !disableErrCaching}, nil
}
