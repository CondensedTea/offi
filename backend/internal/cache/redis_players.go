package cache

import (
	"context"
	"fmt"
	"time"
)

const playerExpiration = 72 * time.Hour

func (r Redis) GetPlayer(ctx context.Context, playerID int) (Player, error) {
	key := fmt.Sprintf("player-%d", playerID)

	var player Player
	err := r.client.Get(ctx, key).Scan(&player)
	return player, err
}

func (r Redis) SetPlayer(ctx context.Context, playerID int, player Player) error {
	key := fmt.Sprintf("player-%d", playerID)

	return r.client.Set(ctx, key, player, playerExpiration).Err()
}
