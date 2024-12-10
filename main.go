package main

import (
	"demoapp/config"
	"demoapp/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)


func main() {
	//  Initialize a new Fiber app
	app := fiber.New()

	config.ConnectDB()
	routes.UserRoute(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}