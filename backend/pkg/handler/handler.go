package handler

import (
	"fmt"
	"offi/pkg/core"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CreateApp(c *core.Core) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		GETOnly:      true,
	})

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://etf2l.org, https://logs.tf",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/match/:matchId", func(ctx *fiber.Ctx) error {
		matchId, err := ctx.ParamsInt("matchId")
		if err != nil {
			return err
		}
		if matchId == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "match id is required")
		}
		logs, err := c.GetLogs(matchId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get logs: %v", err))
		}
		return ctx.JSON(fiber.Map{"logs": logs})
	})

	app.Get("/log/:logId", func(ctx *fiber.Ctx) error {
		logId, err := ctx.ParamsInt("logId")
		if err != nil {
			return err
		}
		if logId == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "log id is required")
		}
		match, err := c.GetMatch(logId)
		if err == redis.Nil {
			return fiber.NewError(fiber.StatusNotFound, "this log does not have linked match yet")
		}
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get match: %v", err))
		}
		return ctx.JSON(fiber.Map{"match": match})
	})

	app.Get("/debug/keys/:hashKey", func(ctx *fiber.Ctx) error {
		hashKey := ctx.Params("hashKey")

		keys, err := c.GetKeys(hashKey)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to get keys: %v", err))
		}
		return ctx.JSON(fiber.Map{"keys": keys})
	})

	return app
}
