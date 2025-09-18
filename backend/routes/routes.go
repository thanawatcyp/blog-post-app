package routes

import (
	"blog-app-backend/handlers"
	"blog-app-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	api := app.Group("/api")

	// Public routes (no authentication required)
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Post("/logout", handlers.Logout)

	// Public health check
	api.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })

	// Memory monitoring endpoints (public for testing)
	api.Get("/memory", handlers.GetMemoryStats)
	api.Post("/memory/gc", handlers.ForceGC)

	// Protected routes (authentication required)
	protected := api.Group("/", middleware.JWTProtected())
	// Posts routes (authenticated users only)
	protected.Get("/posts", handlers.ListPublicPosts)
	protected.Post("/posts/create", handlers.CreatePost)
}
