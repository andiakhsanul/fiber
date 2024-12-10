package routes

import (
	"demoapp/controllers"
	"demoapp/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AdminRoute(app *fiber.App) {
	// Route login tidak memerlukan autentikasi JWT
	app.Post("/login", controllers.LoginHandler)

	// Grup pengguna dengan autentikasi JWT
	adminGroup := app.Group("/admin", middlewares.JWTMiddleware, middlewares.CheckRole("admin"))
	adminGroup.Get("/users", controllers.GetUsers)
	adminGroup.Post("/create", controllers.CreateUser)
	adminGroup.Get("/:userId", controllers.GetAUser)
	adminGroup.Put("/:userId", controllers.EditAUser)
	adminGroup.Delete("/:userId", controllers.DeleteAUser)
}



