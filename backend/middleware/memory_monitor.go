package middleware

import (
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MemoryStats represents memory usage statistics
type MemoryStats struct {
	Alloc        uint64 `json:"alloc"`         // Bytes allocated and in use
	TotalAlloc   uint64 `json:"total_alloc"`   // Total bytes allocated
	Sys          uint64 `json:"sys"`           // Bytes obtained from system
	NumGC        uint32 `json:"num_gc"`        // Number of garbage collections
	HeapAlloc    uint64 `json:"heap_alloc"`    // Bytes allocated and in use on heap
	HeapSys      uint64 `json:"heap_sys"`      // Bytes obtained from system for heap
	HeapIdle     uint64 `json:"heap_idle"`     // Bytes in idle spans
	HeapInuse    uint64 `json:"heap_inuse"`    // Bytes in non-idle spans
	HeapReleased uint64 `json:"heap_released"` // Bytes released to the OS
}

// GetMemoryStats returns current memory statistics
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		NumGC:        m.NumGC,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
	}
}

// BytesToMB converts bytes to megabytes
func BytesToMB(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

// MemoryMonitor middleware logs memory usage for each request
func MemoryMonitor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		startStats := GetMemoryStats()

		// Process request
		err := c.Next()

		// Log memory usage after request
		endStats := GetMemoryStats()
		duration := time.Since(start)

		// Calculate memory delta
		allocDelta := int64(endStats.Alloc) - int64(startStats.Alloc)

		log.Printf(
			"[MEMORY] %s %s | Duration: %v | Memory: %.2fMB | Delta: %+.2fMB | GC: %d",
			c.Method(),
			c.Path(),
			duration,
			BytesToMB(endStats.Alloc),
			BytesToMB(uint64(allocDelta)),
			endStats.NumGC-startStats.NumGC,
		)

		return err
	}
}

// MemoryMonitorDetailed provides more detailed logging for specific routes
func MemoryMonitorDetailed(routeName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		startStats := GetMemoryStats()

		log.Printf(
			"[MEMORY-START] %s | Route: %s | Memory: %.2fMB | Heap: %.2fMB",
			c.Method(),
			routeName,
			BytesToMB(startStats.Alloc),
			BytesToMB(startStats.HeapAlloc),
		)

		// Process request
		err := c.Next()

		// Detailed logging after request
		endStats := GetMemoryStats()
		duration := time.Since(start)

		allocDelta := int64(endStats.Alloc) - int64(startStats.Alloc)
		heapDelta := int64(endStats.HeapAlloc) - int64(startStats.HeapAlloc)

		log.Printf(
			"[MEMORY-END] %s | Route: %s | Duration: %v | Memory: %.2fMB (+%.2fMB) | Heap: %.2fMB (+%.2fMB) | GC: %d | Sys: %.2fMB",
			c.Method(),
			routeName,
			duration,
			BytesToMB(endStats.Alloc),
			BytesToMB(uint64(allocDelta)),
			BytesToMB(endStats.HeapAlloc),
			BytesToMB(uint64(heapDelta)),
			endStats.NumGC-startStats.NumGC,
			BytesToMB(endStats.Sys),
		)

		return err
	}
}