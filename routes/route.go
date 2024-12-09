package routes

import (
	"demoapp/controllers"
	"demoapp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	// Endpoint login tanpa JWTMiddleware
	app.Post("/login", controllers.LoginHandler)
	// app.Post("/create", controllers.CreateUser)

	// Grup user dengan JWTMiddleware
	// userGroup := app.Group("/user", controllers.JWTMiddleware)
	// userGroup.Get("/users", controllers.GetUsers)
	// userGroup.Post("/create", controllers.CreateUser)
	// userGroup.Get("/:userId", controllers.GetAUser)
	// userGroup.Put("/:userId", controllers.EditAUser)
	// userGroup.Delete("/:userId", controllers.DeleteAUser)

	// userGroup.Put("/:userId/password", controllers.EditPassword)
	// userGroup.Post("/:userId/upload", controllers.UploadPhoto)
}

func AdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin", middlewares.CheckRole("admin"))

	adminGroup.Get("/users", controllers.GetUsers)
	adminGroup.Get("/users/:userId", controllers.GetAUser)
	adminGroup.Post("/users", controllers.CreateUser)
	adminGroup.Put("/users/:userId", controllers.EditAUser)
	adminGroup.Delete("/users/:userId", controllers.DeleteAUser)
}
