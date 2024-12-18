package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// CheckRole untuk memverifikasi peran pengguna yang terautentikasi
func CheckRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string) // Ambil peran dari context
		if role != requiredRole {
			// Jika peran tidak sesuai, kembalikan status Forbidden
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}
		return c.Next() // Lanjutkan ke handler berikutnya
	}
}

// CheckJenisUser memeriksa apakah jenis_user sesuai dengan parameter ju
func CheckJenisUser(ju string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil jenis_user dari context
		jenisUser, ok := c.Locals("jenis_user").(string)
		if !ok || jenisUser != ju {
			// Jika tidak sesuai, kembalikan status Forbidden
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: unauthorized user type",
			})
		}
		return c.Next() // Lanjutkan ke handler berikutnya
	}
}
