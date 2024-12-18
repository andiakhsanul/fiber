package controllers

import (
	"demoapp/config"
	"demoapp/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserModulCollection = config.GetCollection(config.DB, "usermodul")
var UserCollection = config.GetCollection(config.DB, "users")

// Fungsi untuk mengganti jenis_user dan memperbarui data user
func ChangeUserType(c *fiber.Ctx) error {
	// Struktur request body
	type RequestBody struct {
		UserID  string `json:"user_id" validate:"required"`
		NewType string `json:"new_type" validate:"required"` // jenis_user yang baru
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi format ID
	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Cari modul yang sudah ada dengan user_id yang sesuai
	var userModul model.UserModul
	err = UserModulCollection.FindOne(c.Context(), bson.M{"user_id": userID}).Decode(&userModul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found in any module"})
	}

	// Hapus user_id dari jenis_user yang lama
	_, err = UserModulCollection.UpdateMany(
		c.Context(),
		bson.M{"jenis_user": bson.M{"$ne": body.NewType}, "user_id": userID}, // Cari berdasarkan user_id dan jenis_user yang bukan new_type
		bson.M{"$pull": bson.M{"user_id": userID}}, // Hapus user_id dari modul yang tidak sesuai
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove user from old modules"})
	}

	// Tambahkan user_id ke jenis_user yang baru
	_, err = UserModulCollection.UpdateOne(
		c.Context(),
		bson.M{"jenis_user": body.NewType}, // Cari modul dengan jenis_user yang baru
		bson.M{"$addToSet": bson.M{"user_id": userID}}, // Menambahkan user_id tanpa duplikasi
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add user to the new type"})
	}

	// Perbarui `jenis_user` di koleksi `user`
	_, err = UserCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": userID}, // Cari user berdasarkan _id
		bson.M{"$set": bson.M{"jenis_user": body.NewType}}, // Perbarui jenis_user
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user type in user collection"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User type updated successfully"})
}