package main

import (
	"fmt"
	"log"

	"gorps/handlers"
	"gorps/libs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load domain from .env file
	domain := libs.LoadEnv("DOMAIN")

	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS for frontend integration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Testing RPS tool is running!")
	})

	// API route to start RPS test
	app.Post("/api/start-test", func(c *fiber.Ctx) error {
		type TestRequest struct {
			Domain     string `json:"domain"`
			Iterations int    `json:"iterations"`
		}

		var req TestRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if req.Domain == "" {
			req.Domain = domain // Use default domain from .env
		}

		if req.Iterations <= 0 {
			req.Iterations = 1 // Default to 1 iteration
		}

		// Start RPS test in goroutine
		go handlers.TestRPSWithIterations(req.Domain, "urls.txt", req.Iterations)

		return c.JSON(fiber.Map{
			"message":    "RPS test started",
			"domain":     req.Domain,
			"iterations": req.Iterations,
			"totalRequests": req.Iterations * 99, // Assuming 99 URLs
		})
	})

	// Start the server
	fmt.Println("Server is running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
