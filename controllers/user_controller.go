package controllers

import (
	"context"
	"demoapp/config"
	"demoapp/model"
	"demoapp/responses"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

// CreateUser - Create a new user
func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User

	// Validate request body
	if err := c.BodyParser(&user); err != nil {
		fmt.Println("Error parsing body:", err)
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data1": "Invalid request body format. " + err.Error()},
		})
	}

	// Validate required fields using validator
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data2": validationErr.Error()},
		})
	}


	// Check if the username already exists in the collection
	var existingUser model.User
	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		// Username already exists
		return c.Status(http.StatusConflict).JSON(responses.UserResponse{
			Status:  http.StatusConflict,
			Message: "error",
			Data:    &fiber.Map{"data": "Username is already taken"},
		})
	} else if err != mongo.ErrNoDocuments {
		// Other error occurred during the lookup
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Error checking username uniqueness"},
		})
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data3": err.Error()},
		})
	}

	// Create new user object
	newUser := model.User{
		ID:       primitive.NewObjectID(),
		Username: user.Username,
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
		Password: string(hashedPassword), // Store hashed password
	}

	// Insert the new user into the MongoDB collection
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data4": err.Error()},
		})
	}

	// Return the result of the insertion (could be user ID or other relevant info)
	return c.Status(http.StatusCreated).JSON(responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data5": result},
	})
}

// GetAUser - Get a single user by ID
func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": user},
	})
}

// EditAUser - Edit a single user by ID
func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	// Validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()},
		})
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	// Get updated user details
	var updatedUser model.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": updatedUser},
	})
}

// DeleteAUser - Delete a single user by ID
func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &fiber.Map{"data": "User with specified ID not found!"},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": "User successfully deleted!"},
	})
}

// GetUsers - Get all users
func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []model.User
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		// Log error if there's an issue with the MongoDB query
		fmt.Println("Error fetching users:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch users from DB"})
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil {
		// Log error if there's an issue decoding the results
		fmt.Println("Error decoding users:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to decode users"})
	}

	return c.Status(http.StatusOK).JSON(users)
}

// LoginHandler - User login handler
func LoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Find user by username
	var user model.User
	err := userCollection.FindOne(context.TODO(), bson.M{
		"username": loginReq.Username,
	}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password"})
	}

	// Generate token (function `generateJWT` should be implemented)
	token := generateJWT(loginReq.Username)
	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}



func EditPassword(c *fiber.Ctx) error {
	// Ambil userId dari parameter URL
	userId := c.Params("userId")
	
	// Bind body JSON ke struct untuk request
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "Invalid request body: " + err.Error()},
		})
	}

	// Convert userId ke ObjectID
	objId, _ := primitive.ObjectIDFromHex(userId)

	// Ambil user dari database
	var user model.User
	err := userCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &fiber.Map{"data": "User not found"},
		})
	}

	// Verifikasi apakah old password sesuai dengan password yang ada di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "error",
			Data:    &fiber.Map{"data": "Incorrect old password"},
		})
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Failed to hash password: " + err.Error()},
		})
	}

	// Update password baru ke database
	update := bson.M{"password": string(hashedPassword)}
	_, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Failed to update password: " + err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Password updated successfully",
		Data:    &fiber.Map{"data": "Password changed successfully"},
	})
}

func UploadPhoto(c *fiber.Ctx) error {
	// Context untuk database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ambil userId dari parameter URL
	userId := c.Params("userId")

	// Parse userId ke ObjectID
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "Invalid User ID"},
		})
	}

	// Ambil file dari form-data dengan key "photo"
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": "Failed to retrieve file"},
		})
	}

	// Buat folder ./storage/images jika belum ada
	imageDir := "./storage/images"
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Failed to create image directory"},
		})
	}

	// Dapatkan ekstensi file
	fileExtension := filepath.Ext(file.Filename)

	// Generate nama file baru dengan format timestamp
	timestamp := time.Now().Format("20060102150405999") // Format YYYYMMDDHHmmSSsss
	newFileName := fmt.Sprintf("%s%s", timestamp, fileExtension)
	filePath := filepath.Join(imageDir, newFileName)

	// Simpan file ke storage/images
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Failed to save file"},
		})
	}

	// Update field photo pada user document
	update := bson.M{"$set": bson.M{"photo": filePath}}
	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": "Failed to update user photo"},
		})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": "Photo uploaded successfully", "photo_path": filePath},
	})
}