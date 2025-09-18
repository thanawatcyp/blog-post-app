package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"blog-app-backend/config"
	"blog-app-backend/models"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
}

var validate = validator.New()

func Register(c *fiber.Ctx) error {
	// 1) parse JSON into our DTO (JSON transform to struct do every time)
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	// 2) validate fields (length, email format, etc.)
	if err := validate.Struct(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// 3) check duplicates (username or email must be unique)
	var cnt int64
	if err := config.DB.Model(&models.User{}).
		Where("username = ? OR email = ?", req.Username, req.Email).
		Count(&cnt).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	if cnt > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "username or email already in use"})
	}

	// 4) hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "hash error"})
	}

	// 5) build the User model (store hash, never plain password)
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		FullName: req.FullName,
	}

	// 6) insert into DB
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "insert error"})
	}

	// 7) return safe JSON (Password is hidden by json:"-")
	return c.Status(http.StatusCreated).JSON(user)
}
