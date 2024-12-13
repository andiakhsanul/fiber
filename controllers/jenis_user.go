package controllers

import (
	"demoapp/config"
	"demoapp/model"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var UserModulCollection = config.GetCollection(config.DB, "usermodul")
// Fungsi untuk memindahkan jenis_user
func ChangeUserType(c *fiber.Ctx) error {
    type RequestBody struct {
        UserID    string   `json:"user_id" validate:"required"`
        NewType   string   `json:"new_type" validate:"required"`
        NewModuls []string `json:"new_moduls" validate:"required"` // Modul baru untuk jenis_user baru
    }

    var body RequestBody
    if err := c.BodyParser(&body); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    userID, err := primitive.ObjectIDFromHex(body.UserID)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    // Kosongkan array `modul_id` untuk `user_id` tertentu
    _, err = UserModulCollection.UpdateMany(
        c.Context(),
        bson.M{"user_id": userID},
        bson.M{"$set": bson.M{"modul_id": []primitive.ObjectID{}}},
    )
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to clear old modules"})
    }

    // Tambahkan modul baru berdasarkan `NewModuls`
    var newModulIDs []primitive.ObjectID
    for _, modulID := range body.NewModuls {
        id, err := primitive.ObjectIDFromHex(modulID)
        if err != nil {
            return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid modul ID"})
        }
        newModulIDs = append(newModulIDs, id)
    }

    userModul := model.UserModul{
        ID:        primitive.NewObjectID(),
        JenisUser: body.NewType,
        UserID:    []primitive.ObjectID{userID},
        ModulID:   newModulIDs,
        CreatedAt: time.Now(),
    }

    // Simpan data baru ke koleksi `usermodul`
    _, err = UserModulCollection.InsertOne(c.Context(), userModul)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to assign new modules"})
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User type updated successfully"})
}