package core

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (c Core) handleGetMatch(ctx *fiber.Ctx) error {
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
}
