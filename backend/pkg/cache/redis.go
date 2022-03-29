package cache

import (
	"strconv"

	"github.com/go-redis/redis"
)

const (
	matchKey  = "matches"
	playerKey = "players"
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

func (r Redis) GetMatch(matchId int) (Match, error) {
	var match Match

	if err := r.client.HGet(matchKey, strconv.Itoa(matchId)).Scan(&match); err != nil {
		return Match{}, err
	}
	return match, nil
}

func (r Redis) SetMatch(matchId int, match *Match) error {
	return r.client.HSet(matchKey, strconv.Itoa(matchId), match).Err()
}

func (r Redis) FlushMatch(matchId int) error {
	return r.client.HDel(matchKey, strconv.Itoa(matchId)).Err()
}

func (r Redis) GetPlayer(playerID string) (string, error) {
	return r.client.HGet(playerKey, playerID).Result()
}

func (r Redis) SetPlayer(playerID, steamID string) error {
	return r.client.HSet(playerKey, playerID, steamID).Err()
}
