package handlers

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"blog-app-backend/config"
	"blog-app-backend/middleware"
	"blog-app-backend/models"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

var loginValidator = validator.New()

func Login(c *fiber.Ctx) error {
	startTime := time.Now()
	startStats := middleware.GetMemoryStats()

	log.Printf("[LOGIN-START] Memory: %.2fMB, Goroutines: %d",
		middleware.BytesToMB(startStats.Alloc), runtime.NumGoroutine())

	// 1) parse JSON into DTO
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	// 2) validate
	if err := loginValidator.Struct(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// 3) find user by email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}

	// 4) compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}

	// 5) generate JWT
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not login"})
	}

	// 6) set httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		HTTPOnly: true,
		Secure:   true, // Only send over HTTPS in production
		SameSite: "Strict",
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
	})

	// 7) log final memory stats
	endStats := middleware.GetMemoryStats()
	duration := time.Since(startTime)
	memoryDelta := int64(endStats.Alloc) - int64(startStats.Alloc)

	log.Printf("[LOGIN-END] Duration: %v, Memory: %.2fMB (+%.2fMB), Goroutines: %d, User: %s",
		duration,
		middleware.BytesToMB(endStats.Alloc),
		middleware.BytesToMB(uint64(memoryDelta)),
		runtime.NumGoroutine(),
		user.Email)

	// 8) return response (without token)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "login successful",
		"user": fiber.Map{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
		},
	})
}
