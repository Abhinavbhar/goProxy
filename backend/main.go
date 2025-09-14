package main

import (
	"backend/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize MongoDB
	routes.InitMongo()
	fmt.Println("0")
	// Create new Fiber app
	app := fiber.New()
	fmt.Println("1")
	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // change later to your extension/frontend
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	fmt.Println("2")
	// Define routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber 🚀")
	})
	app.Post("/login", routes.Login)
	app.Post("/auth/verify", routes.VerifyAuth)
	app.Get("/allowedips", routes.AllowedIp)

	// Start server on port 3000
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting server:", err)
	}
	fmt.Println("3")
}
