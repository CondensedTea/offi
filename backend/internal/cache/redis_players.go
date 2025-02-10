package cache

import (
	"fmt"
	"time"
)

const playerExpiration = 24 * time.Hour

func (r Redis) GetPlayer(playerID int) (Player, error) {
	key := fmt.Sprintf("player-%d", playerID)

	var player Player
	err := r.client.Get(key).Scan(&player)
	return player, err
}

func (r Redis) SetPlayer(playerID int, player Player) error {
	key := fmt.Sprintf("player-%d", playerID)

	return r.client.Set(key, player, playerExpiration).Err()
}
