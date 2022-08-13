package core

import (
	"errors"
	"fmt"
	"offi/pkg/cache"
	"strconv"

	"github.com/go-redis/redis"
)

var ErrTeamNotFound = errors.New("team info is not found")

func (c Core) GetTeam(teamIdString string) (cache.Team, error) {
	teamID, err := strconv.Atoi(teamIdString)
	if err != nil {
		return cache.Team{}, fmt.Errorf("failed to parse team ID: %v", err)
	}

	team, err := c.cache.GetTeam(teamID)
	switch {
	case err == redis.Nil:
		return cache.Team{}, ErrTeamNotFound
	case err != nil:
		return cache.Team{}, fmt.Errorf("failed to get team from cache: %v", err)
	}
	return team, nil
}
