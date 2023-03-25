package handler

import (
	"offi/pkg/core"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Handler struct {
	app          *fiber.App
	sessionStore *session.Store

	core *core.Core
}

func New(c *core.Core) *Handler {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		GETOnly:      true,
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{}))

	handler := &Handler{
		app:          app,
		sessionStore: session.New(),
		core:         c,
	}

	handler.app.Get("/ready", adaptor.HTTPHandlerFunc(c.Ready()))

	handler.app.Get("/match/:matchId", handler.GetMatch)
	handler.app.Get("/log/:logId", handler.GetLog)
	handler.app.Get("/team/:id", handler.GetTeam)
	handler.app.Get("/players", handler.GetPlayers)

	handler.app.Get("/debug/keys/:hashKey", handler.DebugKeys)

	return handler
}

func (h Handler) Run() error {
	return h.app.Listen(":8080")
}
