package main

import (
	"demoapp/config"
	"demoapp/routes"
	"log"


	"github.com/gofiber/fiber/v2/middleware/cache"
	"time"
	"github.com/gofiber/fiber/v2"
)


func main() {
	//  Initialize a new Fiber app
	app := fiber.New()

	cacheMiddleware := cache.New(cache.Config{
		Expiration: 30 * time.Second,
	})
	app.Use(cacheMiddleware)

	config.ConnectDB()
	routes.UserRoute(app)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))

	// cacheStore := cacheMiddleware.Store
	// cacheStore.Reset() // Clear all cache
	// return c.SendString("Cache cleared!")
}