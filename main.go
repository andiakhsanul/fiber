package main

import (
	"demoapp/config"
	"demoapp/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

type JSONText struct {
	Message string
}

func main() {
	//  Initialize a new Fiber app
    app := fiber.New()

	config.ConnectDB()
    routes.UserRoute(app)

    // Start the server on port 3000
    log.Fatal(app.Listen(":3000"))

	// config.ConnectDB()
}