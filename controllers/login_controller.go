package controllers

import (
	"context"
	"demoapp/middlewares"
	"demoapp/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

)





func LoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Mencari pengguna di database
	var user model.User
	err := userCollection.FindOne(context.TODO(), bson.M{
		"username": loginReq.Username,
	}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	// Generate token
	token, err := middlewares.GenerateJWT(user.Username, user.Role, user.JenisUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Kirim token di response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}