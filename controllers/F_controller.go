// TambahModulUserMenentukan menambahkan modul ke user tertentu


package controllers
import(

	"context"
	"demoapp/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func TambahModulUserMenentukan(c *fiber.Ctx) error {
	type Request struct {
		UserID  primitive.ObjectID `json:"user_id" validate:"required"`
		ModulID primitive.ObjectID `json:"modul_id" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Cari UserModul terkait
	var userModul model.UserModul
	err := userModulCollection.FindOne(context.TODO(), bson.M{"user_id": req.UserID}).Decode(&userModul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "UserModul not found"})
	}

	// Cek jika ModulID sudah ada
	for _, modul := range userModul.ModulID {
		if modul == req.ModulID {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Modul already exists for this user"})
		}
	}

	// Tambahkan ModulID baru
	userModul.ModulID = append(userModul.ModulID, req.ModulID)

	// Update ke MongoDB
	update := bson.M{"$set": bson.M{"modul_id": userModul.ModulID}}
	_, err = userModulCollection.UpdateOne(context.TODO(), bson.M{"_id": userModul.ID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul added successfully"})
}

// UpdateModulUserMenentukan memperbarui modul tertentu di user
func UpdateModulUserMenentukan(c *fiber.Ctx) error {
	type Request struct {
		UserID    primitive.ObjectID `json:"user_id" validate:"required"`
		OldModul  primitive.ObjectID `json:"old_modul" validate:"required"`
		NewModul  primitive.ObjectID `json:"new_modul" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Cari UserModul
	var userModul model.UserModul
	err := userModulCollection.FindOne(context.TODO(), bson.M{"user_id": req.UserID}).Decode(&userModul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "UserModul not found"})
	}

	// Cari dan ganti ModulID
	modulUpdated := false
	for i, modul := range userModul.ModulID {
		if modul == req.OldModul {
			userModul.ModulID[i] = req.NewModul
			modulUpdated = true
			break
		}
	}

	if !modulUpdated {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Old modul not found"})
	}

	// Update MongoDB
	update := bson.M{"$set": bson.M{"modul_id": userModul.ModulID}}
	_, err = userModulCollection.UpdateOne(context.TODO(), bson.M{"_id": userModul.ID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul updated successfully"})
}

// HapusModulUserMenentukan menghapus modul dari user
func HapusModulUserMenentukan(c *fiber.Ctx) error {
	type Request struct {
		UserID  primitive.ObjectID `json:"user_id" validate:"required"`
		ModulID primitive.ObjectID `json:"modul_id" validate:"required"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Cari UserModul
	var userModul model.UserModul
	err := userModulCollection.FindOne(context.TODO(), bson.M{"user_id": req.UserID}).Decode(&userModul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "UserModul not found"})
	}

	// Hapus ModulID
	modulDeleted := false
	var updatedModulID []primitive.ObjectID

	for _, modul := range userModul.ModulID {
		if modul != req.ModulID {
			updatedModulID = append(updatedModulID, modul)
		} else {
			modulDeleted = true
		}
	}

	if !modulDeleted {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Modul not found"})
	}

	// Update MongoDB
	update := bson.M{"$set": bson.M{"modul_id": updatedModulID}}
	_, err = userModulCollection.UpdateOne(context.TODO(), bson.M{"_id": userModul.ID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul deleted successfully"})
}