package routes

import (
	"demoapp/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
    // Endpoint login tanpa JWTMiddleware
    app.Post("/login", controllers.LoginHandler)
	// app.Post("/create", controllers.CreateUser)
	

    // Grup user dengan JWTMiddleware
    userGroup := app.Group("/user", controllers.JWTMiddleware)
    userGroup.Get("/users", controllers.GetUsers)
    userGroup.Post("/create", controllers.CreateUser)
    userGroup.Get("/:userId", controllers.GetAUser)
    userGroup.Put("/:userId", controllers.EditAUser)
    userGroup.Delete("/:userId", controllers.DeleteAUser)
	
    userGroup.Put("/:userId/password", controllers.EditPassword)
}
