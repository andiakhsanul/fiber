package controllers

import (
	"demoapp/model"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Struktur request
type UserModuleRequest struct {
	UserID   string   `json:"user_id" validate:"required"`
	ModulIDs []string `json:"modul_ids" validate:"required"`
	Action   string   `json:"action" validate:"required"` // create, update, delete
}

// Fungsi untuk mengelola data usermodul
func ManageUserModule(c *fiber.Ctx) error {
	var req UserModuleRequest

	// Parsing request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi ObjectID
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Konversi ModulID ke ObjectID
	var modulIDs []primitive.ObjectID
	for _, id := range req.ModulIDs {
		modulID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid modul ID: " + id})
		}
		modulIDs = append(modulIDs, modulID)
	}

	// Cari jenis_user dari user
	var user model.User
	err = UserCollection.FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Aksi CUD
	switch req.Action {
	case "create":
		// Tambahkan usermodul jika belum ada
		_, err := UserModulCollection.UpdateOne(
			c.Context(),
			bson.M{"jenis_user": user.JenisUser, "user_id": userID},
			bson.M{
				"$addToSet": bson.M{"modul_id": bson.M{"$each": modulIDs}},
				"$setOnInsert": bson.M{
					"catatan":   "userkhusus",
					"created_at": time.Now(),
				},
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user module"})
		}

	case "update":
		// Tambahkan modul baru
		_, err := UserModulCollection.UpdateOne(
			c.Context(),
			bson.M{"jenis_user": user.JenisUser, "user_id": userID},
			bson.M{"$addToSet": bson.M{"modul_id": bson.M{"$each": modulIDs}}},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user module"})
		}

	case "delete":
		// Hapus modul yang ditentukan
		_, err := UserModulCollection.UpdateOne(
			c.Context(),
			bson.M{"jenis_user": user.JenisUser, "user_id": userID},
			bson.M{"$pull": bson.M{"modul_id": bson.M{"$in": modulIDs}}},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user module"})
		}

	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User module successfully managed"})
}