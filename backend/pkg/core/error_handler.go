package core

import "github.com/gofiber/fiber/v2"

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return ctx.Status(code).JSON(err.Error())
}
