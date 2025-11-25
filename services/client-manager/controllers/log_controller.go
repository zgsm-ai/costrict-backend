package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/zgsm-ai/client-manager/services"
)

/**
 * LogController handles HTTP requests for log operations
 * @description
 * - Implements RESTful API endpoints for log management
 * - Handles request validation and response formatting
 * - Integrates with LogService for business logic
 */
type LogController struct {
	logService *services.LogService
	log        *logrus.Logger
}

/**
 * NewLogController creates a new LogController instance
 * @param {logrus.Logger} log - Logger instance
 * @returns {*LogController} New LogController instance
 */
func NewLogController(log *logrus.Logger) *LogController {
	// Initialize DAOs and services here
	logService := services.NewLogService(nil, log) // Will be properly initialized later

	return &LogController{
		logService: logService,
		log:        log,
	}
}
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	default:
		return ""
	}
}

func getUserId(header http.Header) string {
	// Get Authorization header
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if the header has Bearer prefix
	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[7:] // Remove "Bearer " prefix
	}

	// Parse token without verification (for now)
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return ""
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Extract user_id from claims
		if userID, exists := claims["id"]; exists {
			// Set user_id in request header
			return toString(userID)
		}
	}
	return ""
}

// PostLog handles POST /logs request
// @Summary Create log
// @Description Create a new log record
// @Tags Log
// @Accept json
// @Produce json
// @Param log body map[string]interface{} true "Log data"
// @Success 201 {object} map[string]interface{} "Created log"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /client-manager/api/v1/logs [post]
func (lc *LogController) PostLog(c *gin.Context) {
	// 获取上传的文件
	file, fileHead, err := c.Request.FormFile("logfile")
	if err != nil {
		lc.log.Errorf("get FormFile('logfile') error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()
	var args services.UploadLogArgs
	s := c.Request.FormValue("args")
	if err := json.Unmarshal([]byte(s), &args); err != nil {
		lc.log.Errorf("get FormValue('args') error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := getUserId(c.Request.Header)
	if userId != args.UserID {
		lc.log.Errorf("validate user_id error: args.user_id: %s, token.user_id: %s", args.UserID, userId)
		c.JSON(http.StatusForbidden, gin.H{"error": "userID is invalid"})
		return
	}

	if _, err := lc.logService.CreateLog(context.Background(), &args); err != nil {
		lc.handleError(c, err)
		return
	}

	destPath := filepath.Join("/data", args.ClientID, fileHead.Filename)
	if err := os.MkdirAll(filepath.Join("/data", args.ClientID), 0755); err != nil {
		lc.log.Errorf("Failed to create file: %s, error: %s", destPath, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	destFile, err := os.Create(destPath)
	if err != nil {
		lc.log.Errorf("Failed to create file: %s, error: %s", destPath, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer destFile.Close()
	// 将上传的文件内容复制到目标文件
	if _, err := io.Copy(destFile, file); err != nil {
		lc.log.Errorf("Failed to save file: %s, error: %s", destPath, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    "success",
		"message": fmt.Sprintf("File uploaded successfully: %s", destPath),
	})
}

// GetLogs handles GET /logs/{client_id}/{file_name} request
// @Summary Get logs by client
// @Description Retrieve logs for a specific client with pagination
// @Tags Log
// @Accept json
// @Produce json
// @Param client_id path string true "Client ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(20)
// @Success 200 {object} map[string]interface{} "Logs list with pagination"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /client-manager/api/v1/logs/{client_id}/{file_name} [get]
func (lc *LogController) GetLogs(c *gin.Context) {
	clientID := c.Param("client_id")
	fileName := c.Param("file_name")

	filePath, err := lc.logService.GetLogs(c.Request.Context(), clientID, fileName)
	if err != nil {
		lc.handleError(c, err)
		return
	}

	c.File(filePath)
}

// ListLogs handles GET /logs request
// @Summary Get log statistics
// @Description Retrieve log statistics for a given time period
// @Tags Log
// @Accept json
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "Log statistics"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /client-manager/api/v1/logs [get]
func (lc *LogController) ListLogs(c *gin.Context) {
	// Get query parameters
	var args services.ListLogsArgs
	if err := c.ShouldBindQuery(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "argument.invalid",
			"message": err.Error(),
		})
		return
	}
	// Get log statistics
	logs, paging, err := lc.logService.ListLogs(c.Request.Context(), &args)
	if err != nil {
		lc.handleError(c, err)
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"code":    "success",
		"message": "Log statistics retrieved successfully",
		"data":    logs,
		"paging":  paging,
	})
}

/**
 * handleError handles errors and returns appropriate HTTP responses
 * @param {gin.Context} c - Gin context
 * @param {error} err - Error to handle
 * @description
 * - Maps different error types to appropriate HTTP status codes
 * - Returns standardized error response format
 * - Logs errors for debugging
 */
func (lc *LogController) handleError(c *gin.Context, err error) {
	// Log error
	lc.log.WithError(err).Error("Request processing failed")

	// Handle different error types
	switch e := err.(type) {
	case *services.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "validation.error",
			"message": e.Message,
			"field":   e.Field,
		})
	case *services.ConflictError:
		c.JSON(http.StatusConflict, gin.H{
			"code":    "conflict.error",
			"message": e.Message,
		})
	case *services.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "notfound.error",
			"message": e.Message,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "internal.error",
			"message": "Internal server error",
		})
	}
}

/**
 * SetLogService sets the log service (used for dependency injection)
 * @param {services.LogService} logService - Log service instance
 * @description
 * - Allows setting the log service after controller creation
 * - Used for proper dependency injection
 */
func (lc *LogController) SetLogService(logService *services.LogService) {
	lc.logService = logService
}
