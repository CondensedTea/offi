package cache

import (
	"encoding/json"
	"time"
)

type LogSet struct {
	Logs []Log `json:"logs" redis:"logs"`
}

type Log struct {
	ID          int       `json:"id" redis:"id"`
	Title       string    `json:"title"`
	Map         string    `json:"map" redis:"map"`
	PlayedAt    time.Time `json:"played_at" redis:"played_at"`
	IsSecondary bool      `json:"is_secondary" redis:"is_secondary"`
}

func (ls *LogSet) MarshalBinary() ([]byte, error) {
	return json.Marshal(ls)
}

func (ls *LogSet) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &ls)
}

type MatchPage struct {
	Id          int    `json:"match_id" redis:"match_id"`
	Competition string `json:"competition" redis:"competition"`
	Stage       string `json:"stage" redis:"stage"`
}

func (m MatchPage) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MatchPage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}
