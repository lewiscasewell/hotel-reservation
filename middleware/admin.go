package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lewiscasewell/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.Role != "admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to view this resource",
		})
	}

	return c.Next()
}
