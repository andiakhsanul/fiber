package routes

import (
	"demoapp/controllers"
	"demoapp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Post("/login", controllers.LoginHandler)

	// Grup pengguna dengan autentikasi JWT
	// userGroup := app.Group("/user", middlewares.JWTMiddleware)
	// userGroup.Get("/users", controllers.GetUsers)
	// userGroup.Post("/create", controllers.CreateUser)
	// userGroup.Get("/:userId", controllers.GetAUser)
	// userGroup.Put("/:userId", controllers.EditAUser)
	// userGroup.Delete("/:userId", controllers.DeleteAUser)
}

func AdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin", middlewares.JWTMiddleware, middlewares.CheckRole("admin"))

	adminGroup.Get("/users", controllers.GetUsers)
	adminGroup.Get("/users/:userId", controllers.GetAUser)
	adminGroup.Post("/users", controllers.CreateUser)
	adminGroup.Put("/users/:userId", controllers.EditAUser)
	adminGroup.Delete("/users/:userId", controllers.DeleteAUser)
}