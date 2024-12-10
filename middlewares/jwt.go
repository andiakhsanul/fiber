package middlewares

import (
	"time"
	"os"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// Fungsi untuk membuat token JWT
func GenerateJWT(username, role, jenisUser string) (string, error) {
	claims := jwt.MapClaims{
		"username":  username,
		"role":      role,
		"jenisUser": jenisUser,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}

func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization") // Ambil token dari header
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	// Hilangkan prefix "Bearer " jika ada
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parsing token dan verifikasi
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Mengambil klaim dari token
	claims := token.Claims.(jwt.MapClaims)
	c.Locals("username", claims["username"])
	c.Locals("role", claims["role"])

	//buatlah pengecekan jika role tidak ada maka akan mengeprint kosong
	if claims["role"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Role not found"})
	}

	return c.Next() // Lanjutkan ke handler berikutnya
}