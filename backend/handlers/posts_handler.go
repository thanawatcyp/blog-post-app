package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"blog-app-backend/config"
	"blog-app-backend/models"
	"blog-app-backend/services"
)

type ListPostsResponse struct {
	Items    []models.Post `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author" validate:"required"`
}

var postValidator = validator.New()

// ListPublicPosts → GET /posts
func ListPublicPosts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	q := strings.TrimSpace(c.Query("q", ""))

	db := config.DB.Model(&models.Post{}).Where("published = ?", true)

	if q != "" {
		like := "%" + q + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	var posts []models.Post
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	return c.Status(http.StatusOK).JSON(ListPostsResponse{
		Items:    posts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// CreatePost → POST /posts
func CreatePost(c *fiber.Ctx) error {
	var req CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}

	if err := postValidator.Struct(req); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// Check content for inappropriate language using AI
	contentFilter := services.NewContentFilterService()
	isClean, err := contentFilter.CheckContent(req.Title, req.Content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Content filtering service unavailable. Please try again later."})
	}

	if !isClean {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Your post contains inappropriate content or offensive language. Please review and modify your content before posting."})
	}

	post := models.Post{
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		Published: true,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not create post"})
	}

	return c.Status(http.StatusCreated).JSON(post)
}
