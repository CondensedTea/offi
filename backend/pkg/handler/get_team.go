package handler

import (
	"fmt"
	"offi/pkg/core"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) GetTeam(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you must specify post type and entity ID"))
	}

	team, err := h.core.GetTeam(id)
	switch {
	case err == core.ErrTeamNotFound:
		return fiber.NewError(fiber.StatusNotFound, core.ErrTeamNotFound.Error())
	case err != nil:
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get player recruitment status: %v", err))
	}

	resp := GetTeamResponse{
		Team: team,
	}

	return ctx.JSON(resp)
}
