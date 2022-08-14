package core

import (
	"context"
	"errors"
	"fmt"
	"offi/pkg/cache"
	"offi/pkg/etf2l"

	"github.com/go-redis/redis"
)

// todo: tests this method

func (c Core) GetPlayers(ctx context.Context, playerIDs []int) ([]cache.Player, error) {
	var players []cache.Player

	for _, playerID := range playerIDs {
		player, err := c.cache.GetPlayer(playerID)
		switch {
		case err == redis.Nil:
			etf2lPlayer, etf2lErr := etf2l.GetPlayer(ctx, playerID)
			if errors.Is(err, etf2l.ErrPlayerNotFound) {
				// skipping player if etf2l api could not find them
				continue
			}
			if etf2lErr != nil {
				return nil, fmt.Errorf("failed to get player page from etf2l: %v", etf2lErr)
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
