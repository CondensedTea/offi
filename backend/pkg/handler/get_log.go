package handler

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetLog(ctx *fiber.Ctx) error {
	logId, err := ctx.ParamsInt("logId")
	if err != nil {
		return err
	}
	if logId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "log id is required")
	}
	match, err := h.core.GetMatch(logId)
	if err == redis.Nil {
		return fiber.NewError(fiber.StatusNotFound, "this log does not have linked match yet")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get match: %v", err))
	}

	sess, err := h.sessionStore.Get(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get session: %v", err))
	}
	if err = sess.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save session: %v", err))
	}
	views, err := h.core.CountViews("logs", logId, sess.Fresh())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update view count for log: %v", err))
	}
	return ctx.JSON(fiber.Map{"match": match, "views": views})
}
