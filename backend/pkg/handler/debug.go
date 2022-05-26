package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) Debug(ctx *fiber.Ctx) error {
	hashKey := ctx.Params("hashKey")

	keys, err := h.core.GetKeys(hashKey)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to get keys: %v", err))
	}
	return ctx.JSON(fiber.Map{"keys": keys, "count": len(keys)})
}
