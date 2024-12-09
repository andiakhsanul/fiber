package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole != role {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Role not authorized.",
			})
		}
		return c.Next()
	}
}

func CheckJenis_user(ju string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userJenis := c.Locals("jenis_user")
		if userJenis != ju {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Jenis_user not authorized.",
			})
		}
		return c.Next()
	}
}
