package cache

import (
	"encoding/json"
	"time"
)

type LogSet struct {
	Logs []Log `json:"logs"`
}

type Log struct {
	ID          int       `json:"id" redis:"id"`
	Title       string    `json:"title"`
	Map         string    `json:"map" redis:"map"`
	PlayedAt    time.Time `json:"played_at" redis:"played_at"`
	IsSecondary bool      `json:"is_secondary" redis:"is_secondary"`
}

func (ls LogSet) MarshalBinary() ([]byte, error) {
	return json.Marshal(ls)
}

func (ls *LogSet) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &ls)
}

type MatchPage struct {
	Id          int    `json:"match_id"`
	Competition string `json:"competition"`
	Stage       string `json:"stage"`
	Tier        string `json:"tier"`
}

func (m MatchPage) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MatchPage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}

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
