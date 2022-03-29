package cache

import (
	"encoding/json"
	"time"
)

type Match struct {
	Logs []Log `json:"logs" redis:"logs"`
}

type Log struct {
	ID       int       `json:"id" redis:"id"`
	Map      string    `json:"map" redis:"map"`
	PlayedAt time.Time `json:"played_at" redis:"played_at"`
}

func (m *Match) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Match) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}
