package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetTeam(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you must specify post type and entity ID"))
	}

	entry, err := h.core.GetTeamRecruitmentStatus(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get player recruitment status: %v", err))
	}
	return ctx.JSON(fiber.Map{"status": entry})
}
