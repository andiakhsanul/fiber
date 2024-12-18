package controllers

import (
	"context"
	"demoapp/config"
	"demoapp/model"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Inisialisasi koleksi dan validator
var userModulCollection *mongo.Collection = config.GetCollection(config.DB, "usermodul")
var ModulCollection *mongo.Collection = config.GetCollection(config.DB, "modul")
var validateUserModul = validator.New()

// Create UserModul
func CreateUserModul(c *fiber.Ctx) error {
	userModul := new(model.UserModul)

	// Parsing data dari request body
	if err := c.BodyParser(userModul); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi data
	if err := validateUserModul.Struct(userModul); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set timestamp
	userModul.ID = primitive.NewObjectID()
	userModul.CreatedAt = time.Now()

	// Simpan data ke MongoDB
	result, err := userModulCollection.InsertOne(context.TODO(), userModul)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user modul"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "UserModul created successfully",
		"id":      result.InsertedID,
	})
}

// Get All UserModuls
func GetAllUserModuls(c *fiber.Ctx) error {
	cursor, err := userModulCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user moduls"})
	}
	defer cursor.Close(context.TODO())

	var userModuls []model.UserModul
	if err := cursor.All(context.TODO(), &userModuls); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user moduls"})
	}

	return c.JSON(userModuls)
}

// Get UserModul by ID
func GetUserModules(c *fiber.Ctx) error {
	userIDParam := c.Params("user_id")

	// Validasi ID
	userID, err := primitive.ObjectIDFromHex(userIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Cari semua data usermodul yang berisi user_id
	var userModules []model.UserModul
	cursor, err := UserModulCollection.Find(c.Context(), bson.M{"user_id": userID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user modules"})
	}

	// Decode hasil query
	if err := cursor.All(c.Context(), &userModules); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode user modules"})
	}

	// Jika tidak ada modul
	if len(userModules) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "No modules found for this user"})
	}

	// Ambil nama modul dari ModulCollection
	var moduleNames []string
	for _, userModul := range userModules {
		for _, modulID := range userModul.ModulID {
			var modul model.Modul
			err := ModulCollection.FindOne(c.Context(), bson.M{"_id": modulID}).Decode(&modul)
			if err == nil {
				moduleNames = append(moduleNames, modul.NmModul)
			}
		}
	}

	// Response sukses
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user_id":     userIDParam,
		"modules":     moduleNames,
		"total_count": len(moduleNames),
	})
}

// Update UserModul
func UpdateUserModul(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	userModul := new(model.UserModul)
	if err := c.BodyParser(userModul); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi data
	if err := validateUserModul.Struct(userModul); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Update timestamp
	userModul.CreatedAt = time.Now()

	// Update ke MongoDB
	update := bson.M{
		"$set": bson.M{
			"jenis_user": userModul.JenisUser,
			"user_id":    userModul.UserID,
			"modul_id":   userModul.ModulID,
			"catatan":    userModul.Catatan,
			"created_at": userModul.CreatedAt,
		},
	}

	_, err = userModulCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user modul"})
	}

	return c.JSON(fiber.Map{"message": "UserModul updated successfully"})
}

// Delete UserModul
func DeleteUserModul(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	_, err = userModulCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user modul"})
	}

	return c.JSON(fiber.Map{"message": "UserModul deleted successfully"})
}

func AddModuleToUser(c *fiber.Ctx) error {
	type RequestBody struct {
		UserID  string `json:"user_id" validate:"required"`
		ModulID string `json:"modul_id" validate:"required"`
	}
	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := primitive.ObjectIDFromHex(body.UserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	modulID, err := primitive.ObjectIDFromHex(body.ModulID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid modul ID"})
	}

	// Tambahkan modul ke usermodul
	filter := bson.M{"user_id": userID}
	update := bson.M{"$addToSet": bson.M{"modul_id": modulID}}

	_, err = userModulCollection.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add module"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Module added successfully"})
}