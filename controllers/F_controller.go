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
	UserIDs  []string `json:"user_ids" validate:"required"`
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

	// Validasi ObjectID untuk UserIDs
	var userIDs []primitive.ObjectID
	for _, id := range req.UserIDs {
		userID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID: " + id})
		}
		userIDs = append(userIDs, userID)
	}

	// Validasi ObjectID untuk ModulIDs
	var modulIDs []primitive.ObjectID
	for _, id := range req.ModulIDs {
		modulID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid modul ID: " + id})
		}
		modulIDs = append(modulIDs, modulID)
	}

	// Aksi CUD
	switch req.Action {
	case "create":
		for _, userID := range userIDs {
			var user model.User
			err := UserCollection.FindOne(c.Context(), bson.M{"_id": userID}).Decode(&user)
			if err != nil {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}

			// Buat atau tambahkan modul ke usermodul
			_, err = UserModulCollection.UpdateOne(
				c.Context(),
				bson.M{
					"jenis_user": user.JenisUser,
					"user_id":    bson.M{"$in": []primitive.ObjectID{userID}},
				},
				bson.M{
					"$addToSet": bson.M{
						"modul_id": bson.M{"$each": modulIDs},
					},
					"$setOnInsert": bson.M{
						"user_id":    []primitive.ObjectID{userID},
						"catatan":    "userkhusus",
						"created_at": time.Now(),
					},
				},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user module"})
			}
		}

	case "update":
		_, err := UserModulCollection.UpdateMany(
			c.Context(),
			bson.M{"user_id": bson.M{"$in": userIDs}},
			bson.M{"$addToSet": bson.M{"modul_id": bson.M{"$each": modulIDs}}},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user modules"})
		}

	case "delete":
		_, err := UserModulCollection.UpdateMany(
			c.Context(),
			bson.M{"user_id": bson.M{"$in": userIDs}},
			bson.M{"$pull": bson.M{"modul_id": bson.M{"$in": modulIDs}}},
		)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user modules"})
		}

	default:
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User modules successfully managed"})
}