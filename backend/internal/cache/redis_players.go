package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Player struct {
	DoesntExists bool `json:"doesnt_exists"`

	ID      int         `json:"id"`
	SteamID string      `json:"steam_id"`
	Name    string      `json:"name"`
	Bans    []PlayerBan `json:"bans"`
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

func (r Redis) GetPlayer(ctx context.Context, playerID int) (Player, error) {
	key := fmt.Sprintf("player-%d", playerID)

	var player Player
	err := r.client.Get(ctx, key).Scan(&player)
	return player, err
}

func (r Redis) SetPlayer(ctx context.Context, playerID int, player Player) error {
	const (
		knownPlayerExpiration   = 5 * 24 * time.Hour
		unknownPlayerExpiration = 10 * 24 * time.Hour
	)

	expire := knownPlayerExpiration
	if player.DoesntExists {
		expire = unknownPlayerExpiration
	}

	key := fmt.Sprintf("player-%d", playerID)

	return r.client.Set(ctx, key, player, expire).Err()
}
