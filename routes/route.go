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
	adminGroup.Get("/allmoduls", controllers.GetAllModuls)
	adminGroup.Get("/modul/:modulId", controllers.GetModulByID)
	adminGroup.Get("/usermodul", controllers.GetAllUserModuls)

	adminGroup.Get("/usermodul/:user_id", controllers.GetUserModules)

	adminGroup.Put("/changeusertype", controllers.ChangeUserType)


	adminGroup.Post("/create", controllers.CreateUser)
	adminGroup.Get("/:userId", controllers.GetAUser)
	adminGroup.Put("/:userId", controllers.EditAUser)
	adminGroup.Delete("/:userId", controllers.DeleteAUser)

	// Route khusus untuk upload foto
	adminGroup.Put("/:userId/upload-photo", controllers.UploadPhoto)
	// Route khusu untuk edit password
	adminGroup.Put("/:userId/edit-password", controllers.EditPassword)

	// Grup untuk modul
	
	adminGroup.Post("/modul", controllers.CreateModul)
	adminGroup.Get("/modul/:modulId", controllers.GetModulByID)
	adminGroup.Put("/modul/:modulId", controllers.UpdateModul)
	adminGroup.Delete("/modul/:modulId", controllers.DeleteModul)


	//group untuk usermodul
	
	adminGroup.Post("/usermodul", controllers.CreateUserModul)
	adminGroup.Put("/usermodul/:usermodulId", controllers.UpdateUserModul)
	adminGroup.Delete("/usermodul/:usermodulId", controllers.DeleteUserModul)

	//group untuk ganti jenis user
	

	//group usermodul untuk user tertentu yaitu cud
	adminGroup.Post("/usermodul/manage", controllers.ManageUserModule)

}

