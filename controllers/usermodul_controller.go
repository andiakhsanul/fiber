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
func GetUserModulByID(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var userModul model.UserModul
	err = userModulCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&userModul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "UserModul not found"})
	}

	return c.JSON(userModul)
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