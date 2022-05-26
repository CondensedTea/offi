package core

import (
	"offi/pkg/cache"
	"offi/pkg/etf2l"

	"github.com/go-redis/redis"
)

func (c Core) GetPlayerRecruitmentStatus(id string) (*cache.Entry, error) {
	entry, err := c.cache.GetRecruitmentPost(etf2l.PlayerPost, id)
	if err == redis.Nil {
		return &cache.Entry{Empty: true}, nil
	} else if err != nil {
		return nil, err
	}
	return entry, nil
}

func (c Core) GetTeamRecruitmentStatus(id string) (*cache.Entry, error) {
	entry, err := c.cache.GetRecruitmentPost(etf2l.TeamPost, id)
	if err == redis.Nil {
		return &cache.Entry{Empty: true}, nil
	} else if err != nil {
		return nil, err
	}
	return entry, nil
}
