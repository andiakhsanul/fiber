
package controllers

import (
	"context"
	"demoapp/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

)


func RegisterHandler(c *fiber.Ctx) error {
	// Struktur untuk menerima request
	type RegisterRequest struct {
		Username     string `json:"username" validate:"required"`
		Password     string `json:"password" validate:"required"`
		Email        string `json:"email" validate:"required,email"`
		JenisKelamin int    `json:"jenis_kelamin" validate:"required"`
		Phone        string `json:"phone" validate:"required"`
		JenisUser    string `json:"jenis_user" validate:"required"`
		Photo        string `json:"photo,omitempty"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Cek apakah username atau email sudah digunakan
	count, _ := userCollection.CountDocuments(context.TODO(), bson.M{"username": req.Username})
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	count, _ = userCollection.CountDocuments(context.TODO(), bson.M{"email": req.Email})
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Buat user baru
	newUser := model.User{
		ID:           primitive.NewObjectID(),
		Username:     req.Username,
		NmUser:       req.Username, // Nama user sama dengan username
		Password:     string(hashedPassword),
		Email:        req.Email,
		Role:         "user", // Default role
		CreatedAt:    primitive.NewDateTimeFromTime(time.Now()),
		JenisKelamin: req.JenisKelamin,
		Phone:        req.Phone,
		JenisUser:    "pelanggan",
	}

	_, err = userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}