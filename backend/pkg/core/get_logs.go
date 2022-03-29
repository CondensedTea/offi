package core

import (
	"fmt"
	"offi/pkg/cache"
	"time"

	"github.com/go-redis/redis"
)

func (c Core) GetLogs(matchId int) ([]cache.Log, error) {
	match, err := c.cache.GetMatch(matchId)
	switch {
	case err == redis.Nil:
		return c.saveNewMatch(matchId)
	case err != nil:
		return nil, fmt.Errorf("failed to get match from cache: %v", err)
	}
	return match.Logs, nil
}

func (c Core) saveNewMatch(matchId int) ([]cache.Log, error) {
	cacheLogs := make([]cache.Log, 0)

	match, err := c.etf2l.ParseMatchPage(matchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get players for etf2l match: %v", err)
	}

	var (
		steamID  string
		steamIDs []string
	)

	for _, playerID := range match.Players {
		steamID, err = c.cache.GetPlayer(playerID)
		if err == redis.Nil {
			steamID, err = c.etf2l.ResolvePlayerSteamID(playerID)
			if err != nil {
				return nil, fmt.Errorf("failed to get player page from etf2l: %v\n", err)
			}
			if err = c.cache.SetPlayer(playerID, steamID); err != nil {
				return nil, fmt.Errorf("failed to set player in cache: %v\n", err)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get player from cache: %v\n", err)
		}
		if steamID != "" {
			steamIDs = append(steamIDs, steamID)
		}
	}

	logsMetadata, err := c.logsTf.SearchLogs(steamIDs, match.Maps, match.PlayedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to search logs: %v", err)
	}
	for _, log := range logsMetadata {
		cacheLog := cache.Log{
			ID:       log.Id,
			Map:      log.Map,
			PlayedAt: time.Unix(int64(log.Date), 0),
		}
		cacheLogs = append(cacheLogs, cacheLog)
	}
	if err = c.cache.SetMatch(matchId, &cache.Match{Logs: cacheLogs}); err != nil {
		return nil, fmt.Errorf("failed to set match in cache: %v", err)
	}
	return cacheLogs, nil
}
