package core

import (
	"offi/pkg/cache"
)

func (c Core) GetMatch(logId int) (cache.MatchPage, error) {
	return c.cache.GetMatch(logId)
}
