package controllers

import (
	"context"
	"demoapp/config"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"demoapp/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Global MongoDB Collection untuk Modul
var modulCollection *mongo.Collection = config.GetCollection(config.DB, "modul")
var validater = validator.New()

// Utility untuk upload file
func saveFile(fileHeader *multipart.FileHeader) (string, error) {
	uploadDir := "./uploads/"
	filePath := uploadDir + fileHeader.Filename

	// Cek dan buat direktori jika belum ada
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	// Simpan file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fileHeader.Filename, nil
}

// Create Modul
func CreateModul(c *fiber.Ctx) error {
	modul := new(model.Modul)

	// Parse form data
	modul.NmModul = c.FormValue("nm_modul")
	modul.KetModul = c.FormValue("ket_modul")
	modul.Alamat = c.FormValue("alamat")
	modul.IsAktif = c.FormValue("is_aktif") == "true"

	urutan, err := strconv.Atoi(c.FormValue("urutan"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'urutan' value"})
	}
	modul.Urutan = urutan

	// Handle file upload
	file, err := c.FormFile("gbr_icon")
	if err == nil { // File hanya diproses jika ada
		filename, fileErr := saveFile(file)
		if fileErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": fileErr.Error()})
		}
		modul.Gbr_Icon = filename
	}

	// Set timestamps
	modul.CreatedAt = time.Now()
	modul.UpdatedAt = time.Now()

	// Validasi data
	if err := validater.Struct(modul); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Insert modul ke MongoDB
	result, err := modulCollection.InsertOne(context.TODO(), modul)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create modul"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Modul created successfully",
		"id":      result.InsertedID,
	})
}

// Get All Moduls
func GetAllModuls(c *fiber.Ctx) error {
	cursor, err := modulCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch moduls"})
	}
	defer cursor.Close(context.TODO())

	var moduls []model.Modul
	for cursor.Next(context.TODO()) {
		var modul model.Modul
		if err := cursor.Decode(&modul); err != nil {
			continue
		}
		moduls = append(moduls, modul)
	}

	return c.JSON(moduls)
}

// Get Modul by ID
func GetModulByID(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var modul model.Modul
	err = modulCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&modul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Modul not found"})
	}

	return c.JSON(modul)
}

// Update Modul
func UpdateModul(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	modul := new(model.Modul)

	// Parse form data
	modul.NmModul = c.FormValue("nm_modul")
	modul.KetModul = c.FormValue("ket_modul")
	modul.Alamat = c.FormValue("alamat")
	modul.IsAktif = c.FormValue("is_aktif") == "true"

	urutan, err := strconv.Atoi(c.FormValue("urutan"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'urutan' value"})
	}
	modul.Urutan = urutan

	// Handle file upload jika ada
	file, err := c.FormFile("gbr_icon")
	if err == nil {
		filename, fileErr := saveFile(file)
		if fileErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": fileErr.Error()})
		}
		modul.Gbr_Icon = filename
	}

	modul.UpdatedAt = time.Now()

	update := bson.M{"$set": modul}
	_, err = modulCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update modul"})
	}

	return c.JSON(fiber.Map{"message": "Modul updated successfully"})
}

// Delete Modul
func DeleteModul(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var modul model.Modul
	err = modulCollection.FindOneAndDelete(context.TODO(), bson.M{"_id": objectID}).Decode(&modul)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Modul not found"})
	}

	if modul.Gbr_Icon != "" {
		filePath := "./uploads/" + modul.Gbr_Icon
		os.Remove(filePath)
	}

	return c.JSON(fiber.Map{"message": "Modul deleted successfully"})
}
