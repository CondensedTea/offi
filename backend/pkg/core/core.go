package core

import (
	"offi/pkg/cache"
	"offi/pkg/etf2l"
	"offi/pkg/logstf"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Core struct {
	cache  cache.Cache
	etf2l  *etf2l.ETF2L
	logsTf *logstf.Client
}

func CreateApp(cache cache.Cache, etf2l *etf2l.ETF2L, logsTf *logstf.Client) *fiber.App {
	c := &Core{
		cache:  cache,
		etf2l:  etf2l,
		logsTf: logsTf,
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(logger.New())

	app.Get("/match/:matchId", c.handleGetMatch)

	return app
}
