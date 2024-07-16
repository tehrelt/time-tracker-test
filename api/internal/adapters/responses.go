package adapters

import "github.com/gofiber/fiber/v2"

func internal(c *fiber.Ctx, payload any) error {
	return c.Status(fiber.StatusInternalServerError).JSON(payload)
}
