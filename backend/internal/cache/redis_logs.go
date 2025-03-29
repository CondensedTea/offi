package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/redis/go-redis/v9"
)

func (r Redis) GetLogs(ctx context.Context, matchId int) (LogSet, error) {
	var logSet LogSet

	if err := r.client.HGet(ctx, matchKey, strconv.Itoa(matchId)).Scan(&logSet); err != nil {
		return LogSet{}, err
	}
	return logSet, nil
}

func (r Redis) SetLogs(ctx context.Context, matchId int, logSet *LogSet) error {
	return r.client.HSet(ctx, matchKey, strconv.Itoa(matchId), logSet).Err()
}

func (r Redis) SetLogError(ctx context.Context, matchId int, err error) error {
	match := fmt.Sprintf("match-%d", matchId)
	return r.client.Set(ctx, match, err.Error(), errorMatchExpire).Err()
}

func (r Redis) CheckLogError(ctx context.Context, matchId int) error {
	match := fmt.Sprintf("match-%d", matchId)
	val, err := r.client.Get(ctx, match).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return fmt.Errorf("%w: %v", ErrCached, val)
}

func (r Redis) DeleteLogs(ctx context.Context, matchId int) (*LogSet, error) {
	logSet, err := r.GetLogs(ctx, matchId)
	if err != nil {
		return nil, err
	}
	if err = r.client.HDel(ctx, matchKey, strconv.Itoa(matchId)).Err(); err != nil {
		return nil, err
	}
	return &logSet, nil
}
