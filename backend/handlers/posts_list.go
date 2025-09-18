package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"blog-app-backend/config"
	"blog-app-backend/models"
)

type ListPostsResponse struct {
	Items    []models.Post `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

func ListPublicPosts(c *fiber.Ctx) error {
	// query params: ?page=1&page_size=10&q=keyword
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	q := strings.TrimSpace(c.Query("q", ""))

	// base query: only published posts
	db := config.DB.Model(&models.Post{}).Where("published = ?", true)

	if q != "" {
		like := "%" + q + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	// count total
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
