package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func (h Handler) GetPlayers(ctx *fiber.Ctx) error {
	idsString := ctx.Query("id")
	if idsString == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you must specify player IDs"))
	}

	ids := strings.Split(idsString, ",")

	parsedIDs := lo.FilterMap[string, int](ids, func(id string, _ int) (int, bool) {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return 0, false
		}
		return idInt, true
	})

	players, err := h.core.GetPlayers(parsedIDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get player info: %v", err))
	}

	resp := GetPlayersResponse{
		Players: players,
	}

	return ctx.JSON(resp)
}
