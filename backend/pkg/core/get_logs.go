package core

import (
	"errors"
	"fmt"
	"offi/pkg/cache"

	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

var ErrTooManyPlayers = errors.New("too many players, 18 or less allowed")

const maxPlayers = 18

func (c Core) GetLogs(matchId int) ([]cache.Log, error) {
	logSet, err := c.cache.GetLogs(matchId)
	switch {
	case err == redis.Nil:
		if c.enableErrorCaching {
			if storedErr := c.cache.CheckLogError(matchId); storedErr != nil {
				return nil, storedErr
			}
		}
		logs, saveErr := c.saveNewMatch(matchId)
		if saveErr != nil {
			if c.enableErrorCaching {
				if cacheErr := c.cache.SetLogError(matchId, saveErr); cacheErr != nil {
					logrus.Errorf("failed to cache log error: %v", cacheErr)
				}
			}
			return nil, saveErr
		}
		return logs, nil
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

	if len(match.Players) > maxPlayers {
		return nil, ErrTooManyPlayers
	}

	players, err := c.GetPlayers(match.Players)
	if err != nil {
		return nil, err
	}

	steamIDs := lo.Map(players, func(player cache.Player, _ int) string {
		return player.SteamID
	})

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

func (c Core) CountViews(object string, id int, freshSession bool) (int64, error) {
	if freshSession {
		return c.cache.IncrementViews(object, id)
	}
	count, err := c.cache.GetViews(object, id)
	switch {
	case err == redis.Nil:
		return c.cache.IncrementViews(object, id)
	case err != nil:
		return 0, err
	}
	return count, nil
}
