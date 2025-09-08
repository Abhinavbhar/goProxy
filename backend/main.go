package main

import (
	"backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize MongoDB
	routes.InitMongo()

	// Create new Fiber app
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // change later to your extension/frontend
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber 🚀")
	})
	app.Post("/login", routes.Login)
	app.Post("/auth/verify", routes.VerifyAuth)

	// Start server on port 3000
	app.Listen(":3000")
}
