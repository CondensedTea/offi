package core

import (
	"fmt"
	"offi/pkg/cache"
	"offi/pkg/etf2l"

	"github.com/go-redis/redis"
)

func (c Core) GetLogs(matchId int) ([]cache.Log, error) {
	logSet, err := c.cache.GetLogs(matchId)
	switch {
	case err == redis.Nil:
		return c.saveNewMatch(matchId)
	case err != nil:
		return nil, fmt.Errorf("failed to get match from cache: %v", err)
	}
	return logSet.Logs, nil
}

func (c Core) saveNewMatch(matchId int) ([]cache.Log, error) {
	logIDs := make([]int, 0)

	match, err := c.etf2l.ParseMatchPage(matchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get players for etf2l match: %v", err)
	}

	steamIDs, err := c.GetSteamIDs(match)
	if err != nil {
		return nil, err
	}

	var cacheLogs []cache.Log

	matchLogs, secondaryLogs, err := c.logsTf.SearchLogs(steamIDs, match.Maps, match.PlayedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to search logs: %v", err)
	}
	for _, log := range matchLogs {
		cacheLog := log.ToCache(false)
		logIDs = append(logIDs, log.Id)
		cacheLogs = append(cacheLogs, cacheLog)
	}
	for _, log := range secondaryLogs {
		cacheLog := log.ToCache(true)
		logIDs = append(logIDs, log.Id)
		cacheLogs = append(cacheLogs, cacheLog)
	}
	if err = c.cache.SetLogs(matchId, &cache.LogSet{Logs: cacheLogs}); err != nil {
		return nil, fmt.Errorf("failed to set match in cache: %v", err)
	}
	if err = c.cache.SetMatch(logIDs, &cache.MatchPage{
		Id:          match.ID,
		Competition: match.Competition,
		Stage:       match.Stage,
	}); err != nil {
		return nil, fmt.Errorf("failed to set logs in cache: %v", err)
	}
	return cacheLogs, nil
}

func (c Core) GetSteamIDs(match *etf2l.Match) ([]string, error) {
	var steamIDs []string

	for _, playerURL := range match.Players {
		steamID, err := c.cache.GetPlayer(playerURL)
		if err == redis.Nil {
			steamID, err = c.etf2l.ResolvePlayerSteamID(playerURL)
			if err != nil {
				return nil, fmt.Errorf("failed to get player page from etf2l: %v", err)
			}
			if err = c.cache.SetPlayer(playerURL, steamID); err != nil {
				return nil, fmt.Errorf("failed to save player in cache: %v", err)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get player from cache: %v\n", err)
		}
		if steamID != "" {
			steamIDs = append(steamIDs, steamID)
		}
	}
	return steamIDs, nil
}
