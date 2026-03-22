package main

import (
	"fmt"
	"log"
	"strconv"

	"gorps/handlers"
	"gorps/libs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Load domain from .env file
	domain := libs.LoadEnv("DOMAIN")

	// Initialize SQLite database
	libs.InitDB("cargo.db")

	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS for frontend integration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://localhost:5177",
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
			req.Domain = domain
		}

		if req.Iterations <= 0 {
			req.Iterations = 1
		}

		// Start RPS test in goroutine
		go handlers.TestRPSWithIterations(req.Domain, "urls.txt", req.Iterations)

		return c.JSON(fiber.Map{
			"message":       "RPS test started",
			"domain":        req.Domain,
			"iterations":    req.Iterations,
			"totalRequests": req.Iterations * 99,
		})
	})

	// API route to get results for the current test run
	app.Get("/api/results", func(c *fiber.Ctx) error {
		sinceStr := c.Query("since", "0")
		since, err := strconv.Atoi(sinceStr)
		if err != nil {
			since = 0
		}

		runID, running, total := handlers.GetCurrentRunState()
		if runID == 0 {
			return c.JSON(fiber.Map{
				"results": []interface{}{},
				"running": false,
				"total":   0,
			})
		}

		results, err := libs.GetResultsSinceDB(runID, since)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch results"})
		}

		if results == nil {
			results = []libs.TestRunResult{}
		}

		return c.JSON(fiber.Map{
			"results": results,
			"running": running,
			"total":   total,
			"runId":   runID,
		})
	})

	// API route to list all test runs
	app.Get("/api/runs", func(c *fiber.Ctx) error {
		runs, err := libs.GetAllTestRuns()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch test runs"})
		}
		return c.JSON(fiber.Map{"runs": runs})
	})

	// API route to get results for a specific test run
	app.Get("/api/runs/:id/results", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid run ID"})
		}

		run, err := libs.GetTestRun(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Test run not found"})
		}

		results, err := libs.GetResultsSinceDB(id, 0)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch results"})
		}

		return c.JSON(fiber.Map{
			"run":     run,
			"results": results,
		})
	})

	// Start the server
	fmt.Println("Server is running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
