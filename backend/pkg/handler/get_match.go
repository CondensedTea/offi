package handler

import (
	"errors"
	"fmt"
	"offi/pkg/cache"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetMatch(ctx *fiber.Ctx) error {
	matchId, err := ctx.ParamsInt("matchId")
	if err != nil {
		return err
	}
	if matchId == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "match id is required")
	}
	logs, err := h.core.GetLogs(matchId)
	if errors.Is(err, cache.ErrCached) {
		return fiber.NewError(fiber.StatusTooEarly, err.Error())
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get logs: %v", err))
	}

	sess, err := h.sessionStore.Get(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get session: %v", err))
	}

	if err = sess.Save(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save session: %v", err))
	}

	views, err := h.core.CountViews("match", matchId, sess.Fresh())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get view count for match: %v", err))
	}
	return ctx.JSON(fiber.Map{"logs": logs, "views": views})
}
