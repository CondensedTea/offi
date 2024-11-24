package cache

import (
	"fmt"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/go-redis/redis"
)

func (r Redis) GetLogs(matchId int) (LogSet, error) {
	var logSet LogSet

	if err := r.client.HGet(matchKey, strconv.Itoa(matchId)).Scan(&logSet); err != nil {
		return LogSet{}, err
	}
	return logSet, nil
}

func (r Redis) SetLogs(matchId int, logSet *LogSet) error {
	return r.client.HSet(matchKey, strconv.Itoa(matchId), logSet).Err()
}

func (r Redis) SetLogError(matchId int, err error) error {
	match := fmt.Sprintf("match-%d", matchId)
	return r.client.Set(match, err.Error(), errorMatchExpire).Err()
}

func (r Redis) CheckLogError(matchId int) error {
	match := fmt.Sprintf("match-%d", matchId)
	val, err := r.client.Get(match).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return fmt.Errorf("%w: %v", ErrCached, val)
}

func (r Redis) DeleteLogs(matchId int) (*LogSet, error) {
	logSet, err := r.GetLogs(matchId)
	if err != nil {
		return nil, err
	}
	if err = r.client.HDel(matchKey, strconv.Itoa(matchId)).Err(); err != nil {
		return nil, err
	}
	return &logSet, nil
}
