package cache

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

const (
	matchKey  = "matches"
	playerKey = "players"
	logsKey   = "logs"
)

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

func (r Redis) GetPlayer(playerID string) (string, error) {
	return r.client.HGet(playerKey, playerID).Result()
}

func (r Redis) SetPlayer(playerID, steamID string) error {
	return r.client.HSet(playerKey, playerID, steamID).Err()
}

func (r Redis) GetMatch(logId int) (MatchPage, error) {
	var mp MatchPage
	if err := r.client.HGet(logsKey, strconv.Itoa(logId)).Scan(&mp); err != nil {
		return MatchPage{}, err
	}
	return mp, nil
}

func (r Redis) SetMatch(logIds []int, matchPage *MatchPage) error {
	var err error

	for _, id := range logIds {
		if err = r.client.HSet(logsKey, strconv.Itoa(id), matchPage).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r Redis) GetAllKeys(hashKey string) ([]string, error) {
	var keys []string

	switch hashKey {
	case logsKey, playerKey, matchKey:
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
