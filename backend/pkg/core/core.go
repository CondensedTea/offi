package core

import (
	"offi/pkg/cache"
	"offi/pkg/etf2l"
	"offi/pkg/logstf"
)

type Core struct {
	cache  cache.Cache
	etf2l  *etf2l.Client
	logsTf *logstf.Client
}

func New(cache cache.Cache, etf2l *etf2l.Client, logsTf *logstf.Client) *Core {
	return &Core{
		cache:  cache,
		etf2l:  etf2l,
		logsTf: logsTf,
	}
}
