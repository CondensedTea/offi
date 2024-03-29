package core

import (
	"errors"
	"fmt"
	"offi/pkg/cache"
	"offi/pkg/etf2l"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func (c Core) GetPlayers(playerIDs []int) ([]cache.Player, error) {
	var players []cache.Player

	for _, playerID := range playerIDs {
		player, err := c.cache.GetPlayer(playerID)
		switch {
		case err == redis.Nil:
			etf2lPlayer, etf2lErr := c.etf2l.GetPlayer(playerID)
			switch {
			case errors.Is(etf2lErr, etf2l.ErrPlayerNotFound):
				logrus.Warnf("failed to get player %d from etf2l: %v", playerID, etf2lErr)
				continue
			case etf2lErr != nil:
				return nil, fmt.Errorf("failed to get player from etf2l: %v", etf2lErr)
			}

			player = etf2lPlayer.ToCache()
			if cacheErr := c.cache.SetPlayer(playerID, player); cacheErr != nil {
				return nil, fmt.Errorf("failed to save player to cache: %v", cacheErr)
			}
		case err != nil:
			return nil, fmt.Errorf("failed to get player from cache: %v", err)
		}
		players = append(players, player)
	}

	return players, nil
}
