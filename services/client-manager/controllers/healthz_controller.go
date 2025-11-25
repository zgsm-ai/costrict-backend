package controllers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/zgsm-ai/client-manager/utils"
)

/**
 * HealthController handles HTTP requests for health operations
 * @description
 * - Implements health check endpoints
 * - Provides service status information
 * - Supports Kubernetes health checks
 */
type HealthController struct {
	log *logrus.Logger
}

/**
 * NewHealthController creates a new HealthController instance
 * @param {logrus.Logger} log - Logger instance
 * @returns {*HealthController} New HealthController instance
 */
func NewHealthController(log *logrus.Logger) *HealthController {
	return &HealthController{
		log: log,
	}
}

// GetHealth handles GET /healthz request
// @Summary Health check endpoint
// @Description Check the health status of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status"
// @Failure 500 {object} map[string]interface{} "Service unhealthy"
// @Router /healthz [get]
func (hc *HealthController) GetHealth(c *gin.Context) {
	// Get startup time
	startupTime := utils.GetStartupTime()
	if startupTime.IsZero() {
		startupTime = time.Now()
	}

	// Calculate uptime
	uptime := time.Since(startupTime)

	// Get memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Prepare health status
	healthStatus := map[string]interface{}{
		"status":       "healthy",
		"timestamp":    time.Now().Format(time.RFC3339),
		"version":      "1.0.0",
		"startup_time": startupTime.Format(time.RFC3339),
		"uptime":       uptime.String(),
		"memory": map[string]interface{}{
			"alloc":       memStats.Alloc,
			"total_alloc": memStats.TotalAlloc,
			"sys":         memStats.Sys,
			"num_gc":      memStats.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
	}

	// Get request count from utils
	requestCount := utils.GetRequestCount()
	errorCount := utils.GetErrorCount()

	healthStatus["requests"] = map[string]interface{}{
		"total":  requestCount,
		"errors": errorCount,
	}

	// Log health check
	hc.log.Info("Health check requested")

	// Return health status
	c.JSON(http.StatusOK, gin.H{
		"code":    "success",
		"message": "Service is healthy",
		"data":    healthStatus,
	})
}

// LiveHandler handles GET /live request
// @Summary Liveness check endpoint
// @Description Check if the service is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Liveness status"
// @Failure 500 {object} map[string]interface{} "Service not alive"
// @Router /live [get]
func (hc *HealthController) LiveHandler(c *gin.Context) {
	// Check if service is alive
	isAlive := true

	if isAlive {
		c.JSON(http.StatusOK, gin.H{
			"code":    "success",
			"message": "Service is alive",
			"data": map[string]interface{}{
				"status":    "alive",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "service.dead",
			"message": "Service is not alive",
			"data": map[string]interface{}{
				"status":    "dead",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	}
}

// ReadyHandler handles GET /ready request
// @Summary Readiness check endpoint
// @Description Check if the service is ready to accept traffic
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Readiness status"
// @Failure 503 {object} map[string]interface{} "Service not ready"
// @Router /ready [get]
func (hc *HealthController) ReadyHandler(c *gin.Context) {
	// Check if service is ready
	// This can be extended to check database connections, external services, etc.
	isReady := true

	if isReady {
		c.JSON(http.StatusOK, gin.H{
			"code":    "success",
			"message": "Service is ready",
			"data": map[string]interface{}{
				"status":    "ready",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"code":    "service.unavailable",
			"message": "Service is not ready",
			"data": map[string]interface{}{
				"status":    "not_ready",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	}
}
