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

	// Grup untuk modul
	adminGroup.Get("/modul", controllers.GetAllModuls)
	adminGroup.Post("/modul", controllers.CreateModul)
	adminGroup.Get("/modul/:modulId", controllers.GetModulByID)
	adminGroup.Put("/modul/:modulId", controllers.UpdateModul)
	adminGroup.Delete("/modul/:modulId", controllers.DeleteModul)


	//gorup unduk usermodul
	adminGroup.Get("/usermodul", controllers.GetAllUserModuls)
	adminGroup.Post("/usermodul", controllers.CreateUserModul)
	adminGroup.Get("/usermodul/:usermodulId", controllers.GetUserModulByID)
	adminGroup.Put("/usermodul/:usermodulId", controllers.UpdateUserModul)
	adminGroup.Delete("/usermodul/:usermodulId", controllers.DeleteUserModul)

	//group untuk ganti jenis user
	adminGroup.Put("/changeusertype", controllers.ChangeUserType)


}



