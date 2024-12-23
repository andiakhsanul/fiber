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
}

// ------------------------------
// Fungsi CUD yang Dipisah
// ------------------------------

// Fungsi CREATE
func CreateUserModule(c *fiber.Ctx) error {
	var req UserModuleRequest

	// Parsing dan validasi
	if err := parseAndValidateRequest(c, &req); err != nil {
		return err
	}

	// Proses pembuatan modul untuk setiap user
	for _, userID := range req.UserIDs {
		// Ambil data user
		var user model.User
		oid, _ := primitive.ObjectIDFromHex(userID)
		err := UserCollection.FindOne(c.Context(), bson.M{"_id": oid}).Decode(&user)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		// Insert atau Update modul
		_, err = UserModulCollection.UpdateOne(
			c.Context(),
			bson.M{"jenis_user": user.JenisUser, "user_id": bson.M{"$in": []primitive.ObjectID{oid}}},
			bson.M{
				"$addToSet": bson.M{"modul_id": bson.M{"$each": parseObjectIDs(req.ModulIDs)}},
				"$setOnInsert": bson.M{
					"user_id":    []primitive.ObjectID{oid},
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

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User modules successfully created"})
}

// Fungsi UPDATE
func UpdateUserModule(c *fiber.Ctx) error {
	var req UserModuleRequest

	// Parsing dan validasi
	if err := parseAndValidateRequest(c, &req); err != nil {
		return err
	}

	// Update modul untuk user
	_, err := UserModulCollection.UpdateMany(
		c.Context(),
		bson.M{"user_id": bson.M{"$in": parseObjectIDs(req.UserIDs)}},
		bson.M{"$addToSet": bson.M{"modul_id": bson.M{"$each": parseObjectIDs(req.ModulIDs)}}},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user modules"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User modules successfully updated"})
}

// Fungsi DELETE
func DeleteUserModule(c *fiber.Ctx) error {
	var req UserModuleRequest

	// Parsing dan validasi
	if err := parseAndValidateRequest(c, &req); err != nil {
		return err
	}

	// Hapus modul dari user
	_, err := UserModulCollection.UpdateMany(
		c.Context(),
		bson.M{"user_id": bson.M{"$in": parseObjectIDs(req.UserIDs)}},
		bson.M{"$pull": bson.M{"modul_id": bson.M{"$in": parseObjectIDs(req.ModulIDs)}}},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user modules"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User modules successfully deleted"})
}

// ------------------------------
// Fungsi Utility (Helper)
// ------------------------------

// Fungsi untuk parsing dan validasi request
func parseAndValidateRequest(c *fiber.Ctx, req *UserModuleRequest) error {
	// Parsing request
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validasi User IDs
	for _, id := range req.UserIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID: " + id})
		}
	}

	// Validasi Modul IDs
	for _, id := range req.ModulIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid modul ID: " + id})
		}
	}

	return nil
}

// Fungsi untuk konversi array string ke array ObjectID
func parseObjectIDs(ids []string) []primitive.ObjectID {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		oid, _ := primitive.ObjectIDFromHex(id)
		objectIDs = append(objectIDs, oid)
	}
	return objectIDs
}