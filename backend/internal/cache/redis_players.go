package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Player struct {
	DoesntExists bool `json:"doesnt_exists"`

	SteamID int64  `json:"steam_id,omitempty"`
	Name    string `json:"name,omitempty"`

	// Only for ETF2L
	ID   int         `json:"id,omitempty"`
	Bans []PlayerBan `json:"bans,omitempty"`
}

type PlayerBan struct {
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Reason string `json:"reason"`
}

func (p Player) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

func (p *Player) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &p)
}

func (r Redis) GetPlayer(ctx context.Context, league string, playerID int64) (Player, error) {
	key := fmt.Sprintf("%s-player-%d", league, playerID)

	var player Player
	err := r.client.Get(ctx, key).Scan(&player)
	return player, err
}

func (r Redis) GetPlayers(ctx context.Context, league string, playerIDs []int64) (map[int64]*Player, error) {
	keys := make([]string, len(playerIDs))
	for i, p := range playerIDs {
		keys[i] = fmt.Sprintf("%s-player-%d", league, p)
	}

	players := make(map[int64]*Player, len(playerIDs))
	results, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get players from cache: %w", err)
	}

	for i, result := range results {
		if result == nil {
			players[playerIDs[i]] = nil
		} else {
			var player Player
			if err = json.Unmarshal([]byte(result.(string)), &player); err != nil {
				return nil, fmt.Errorf("unmarshaling player: %w", err)
			}

			players[playerIDs[i]] = &player
		}
	}

	return players, nil
}

func (r Redis) SetPlayer(ctx context.Context, league string, playerID int64, player Player) error {
	const (
		knownPlayerExpiration   = 5 * 24 * time.Hour
		unknownPlayerExpiration = 10 * 24 * time.Hour
	)

	expire := knownPlayerExpiration
	if player.DoesntExists {
		expire = unknownPlayerExpiration
	}

	key := fmt.Sprintf("%s-player-%d", league, playerID)

	return r.client.Set(ctx, key, player, expire).Err()
}
