package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"strings"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

func GenerateJWT(username, role, jenisUser string) (string, error) {
	claims := jwt.MapClaims{
		"username":  username,
		"role":      role,
		"jenisUser": jenisUser,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("username", claims["username"])
	c.Locals("role", claims["role"])

	return c.Next()
}