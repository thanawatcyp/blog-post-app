package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	// Clear the auth cookie by setting it to expire in the past
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "logout successful",
	})
}