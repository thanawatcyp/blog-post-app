package main

import (
	"blog-app-backend/config"
	"blog-app-backend/models"
	"blog-app-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	// Connect to database
	config.ConnectDB()

	// Auto-migrate the schema
	err := config.DB.AutoMigrate(
		&models.Post{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS for frontend communication
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// API routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Blog App API is running!",
			"status":  "success",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
		})
	})

	// Register API routes
	routes.Register(app)

	// Start server
	log.Println("Server starting on :4000")
	app.Listen(":4000")
}
