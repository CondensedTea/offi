package logstf

import (
	"offi/pkg/cache"
	"time"
)

type Response struct {
	Logs []Log `json:"logs"`
}

type Log struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Map     string `json:"map"`
	Date    int64  `json:"date"`
	Views   int    `json:"views"`
	Players int    `json:"players"`
}

func (l Log) ToCache(isSecondary bool) cache.Log {
	return cache.Log{
		ID:          l.Id,
		Title:       l.Title,
		Map:         l.Map,
		PlayedAt:    time.Unix(int64(l.Date), 0),
		IsSecondary: isSecondary,
	}
}

type FullLog struct {
	Info Log `json:"info"`
}
