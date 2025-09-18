package handlers

import (
	"blog-app-backend/middleware"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DetailedMemoryStats includes additional system information
type DetailedMemoryStats struct {
	middleware.MemoryStats
	Timestamp    time.Time `json:"timestamp"`
	Goroutines   int       `json:"goroutines"`
	CPUCount     int       `json:"cpu_count"`
	GoVersion    string    `json:"go_version"`
	AllocMB      float64   `json:"alloc_mb"`
	SysMB        float64   `json:"sys_mb"`
	HeapAllocMB  float64   `json:"heap_alloc_mb"`
	HeapSysMB    float64   `json:"heap_sys_mb"`
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats(c *fiber.Ctx) error {
	stats := middleware.GetMemoryStats()

	detailed := DetailedMemoryStats{
		MemoryStats: stats,
		Timestamp:   time.Now(),
		Goroutines:  runtime.NumGoroutine(),
		CPUCount:    runtime.NumCPU(),
		GoVersion:   runtime.Version(),
		AllocMB:     middleware.BytesToMB(stats.Alloc),
		SysMB:       middleware.BytesToMB(stats.Sys),
		HeapAllocMB: middleware.BytesToMB(stats.HeapAlloc),
		HeapSysMB:   middleware.BytesToMB(stats.HeapSys),
	}

	return c.Status(http.StatusOK).JSON(detailed)
}

// ForceGC forces garbage collection and returns memory stats
func ForceGC(c *fiber.Ctx) error {
	beforeStats := middleware.GetMemoryStats()

	// Force garbage collection
	runtime.GC()
	runtime.GC() // Call twice for more thorough collection

	afterStats := middleware.GetMemoryStats()

	response := fiber.Map{
		"message": "Garbage collection forced",
		"before": fiber.Map{
			"alloc_mb":      middleware.BytesToMB(beforeStats.Alloc),
			"heap_alloc_mb": middleware.BytesToMB(beforeStats.HeapAlloc),
			"sys_mb":        middleware.BytesToMB(beforeStats.Sys),
		},
		"after": fiber.Map{
			"alloc_mb":      middleware.BytesToMB(afterStats.Alloc),
			"heap_alloc_mb": middleware.BytesToMB(afterStats.HeapAlloc),
			"sys_mb":        middleware.BytesToMB(afterStats.Sys),
		},
		"freed_mb": middleware.BytesToMB(beforeStats.Alloc - afterStats.Alloc),
	}

	return c.Status(http.StatusOK).JSON(response)
}