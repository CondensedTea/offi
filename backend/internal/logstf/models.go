package logstf

import (
	"offi/internal/cache"
	"time"
)

type Response struct {
	Logs []Log `json:"logs"`
}

type Log struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Map   string `json:"map"`
	Date  int64  `json:"date"`
}

func (l Log) ToCache(isSecondary bool) cache.Log {
	return cache.Log{
		ID:          l.Id,
		Title:       l.Title,
		Map:         l.Map,
		PlayedAt:    time.Unix(l.Date, 0),
		IsSecondary: isSecondary,
	}
}
