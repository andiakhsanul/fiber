package routes

import (
	"demoapp/controllers"
	"demoapp/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	// Route login tidak memerlukan autentikasi JWT
	app.Post("/login", controllers.LoginHandler)

	// Grup pengguna dengan autentikasi JWT
	userGroup := app.Group("/admin", middlewares.JWTMiddleware, middlewares.CheckRole("admin"))
	userGroup.Get("/users", controllers.GetUsers)
	userGroup.Post("/create", controllers.CreateUser)
	userGroup.Get("/:userId", controllers.GetAUser)
	userGroup.Put("/:userId", controllers.EditAUser)
	userGroup.Delete("/:userId", controllers.DeleteAUser)
}
