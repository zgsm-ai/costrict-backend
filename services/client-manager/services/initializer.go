package services

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zgsm-ai/client-manager/dao"
	"github.com/zgsm-ai/client-manager/internal"
	"github.com/zgsm-ai/client-manager/utils"
)

// AppContext holds all the core application objects
type AppContext struct {
	DB         *gorm.DB
	Logger     *logrus.Logger
	LogDAO     *dao.LogDAO
	LogService *LogService
}

// InitializeApp initializes all core application objects and returns AppContext
/**
 * Initialize application core objects
 * @returns {*AppContext, error} Application context and error if initialization fails
 * @description
 * - Initializes database connection
 * - Initializes Prometheus metrics
 * - Creates all DAO objects
 * - Creates all service objects
 * - Creates all controller objects
 * @throws
 * - Database initialization error
 */
func InitializeApp() (*AppContext, error) {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	// Initialize database
	db, err := internal.InitDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Initialize Prometheus metrics
	internal.InitMetrics()

	// Initialize DAOs
	logDAO := dao.NewLogDAO(db, logger)

	// Initialize services
	logService := NewLogService(logDAO, logger)

	// Create and return app context
	appContext := &AppContext{
		DB:         db,
		Logger:     logger,
		LogDAO:     logDAO,
		LogService: logService,
	}

	return appContext, nil
}

// StartServer starts the HTTP server
/**
 * Start HTTP server
 * @param {*gin.Engine} r - Gin engine
 * @param {*logrus.Logger} logger - Application logger
 * @description
 * - Gets server port from configuration
 * - Records startup time
 * - Starts the HTTP server
 * @throws
 * - Server start error
 */
func StartServer(r *gin.Engine, logger *logrus.Logger) error {
	// Get port from configuration
	listenAddr := internal.GetListenAddr()

	// Start server
	logger.Infof("Starting server on %s", listenAddr)

	// Record startup time
	utils.SetStartupTime(time.Now())

	return r.Run(listenAddr)
}
