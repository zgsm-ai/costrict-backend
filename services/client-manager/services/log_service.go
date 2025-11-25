package services

import (
	"context"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zgsm-ai/client-manager/dao"
	"github.com/zgsm-ai/client-manager/models"
)

/**
 * LogService handles business logic for log operations
 * @description
 * - Implements log processing business rules
 * - Validates log data
 * - Handles different log types
 */
type LogService struct {
	logDAO *dao.LogDAO
	log    *logrus.Logger
}

type UploadLogArgs struct {
	ClientID    string `json:"client_id"`
	UserID      string `json:"user_id"`
	FileName    string `json:"file_name"`
	FirstLineNo int64  `json:"first_line_no"`
	LastLineNo  int64  `json:"end_line_no"`
}

type ListLogsArgs struct {
	ClientId string `form:"client_id"`
	UserId   string `form:"user_id"`
	FileName string `form:"file_name"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

type GetLogArgs struct {
	ClientID string `form:"client_id"`
	UserID   string `form:"user_id"`
	FileName string `form:"file_name"`
}

type LogStats struct {
	FirstLineNo int64 //首行编号
	LastLineNo  int64 //尾行编号
}

type Paginated struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

/**
 * NewLogService creates a new LogService instance
 * @param {dao.LogDAO} logDAO - Log data access object
 * @param {logrus.Logger} log - Logger instance
 * @returns {*LogService} New LogService instance
 */
func NewLogService(logDAO *dao.LogDAO, log *logrus.Logger) *LogService {
	return &LogService{
		logDAO: logDAO,
		log:    log,
	}
}

/**
 * CreateLog creates a new log record
 * @param {context.Context} ctx - Context for request cancellation
 * @param {map[string]interface{}} data - Log data
 * @returns {*models.Log, error} Created log and error if any
 * @description
 * - Validates log data
 * - Creates log record
 * - Logs creation operation
 * @throws
 * - Validation errors for invalid data
 * - Database creation errors
 */
func (s *LogService) CreateLog(ctx context.Context, args *UploadLogArgs) (*models.Log, error) {
	// Validate and extract log data
	err := s.validate(args)
	if err != nil {
		s.log.WithError(err).WithFields(logrus.Fields{
			"client_id": args.ClientID,
			"user_id":   args.UserID,
			"file_name": args.FileName,
		}).Error("Invalid arguments")
		return nil, err
	}

	// Create log
	log := &models.Log{
		ClientID:  args.ClientID,
		UserID:    args.UserID,
		FileName:  args.FileName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Create log
	err = s.logDAO.Upsert(ctx, log)
	if err != nil {
		s.log.WithError(err).WithFields(logrus.Fields{
			"client_id": log.ClientID,
			"user_id":   log.UserID,
			"file_name": log.FileName,
		}).Error("Failed to create log")
		return nil, err
	}

	s.log.WithFields(logrus.Fields{
		"client_id": log.ClientID,
		"user_id":   log.UserID,
		"file_name": log.FileName,
	}).Info("Log created successfully")

	return log, nil
}

/**
 * GetLogs retrieves logs for a specific client
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} clientID - Client identifier
 * @param {int} page - Page number
 * @param {int} pageSize - Number of items per page
 * @returns {map[string]interface{}, error} Response containing logs and pagination info
 * @description
 * - Validates client ID and pagination parameters
 * - Retrieves logs from database
 * - Returns structured response with pagination metadata
 * @throws
 * - Validation errors for invalid parameters
 * - Database query errors
 */
func (s *LogService) GetLogs(ctx context.Context, clientID, fname string) (string, error) {
	if clientID == "" {
		return "", &ValidationError{Field: "client_id", Message: "client_id is required"}
	}
	if fname == "" {
		return "", &ValidationError{Field: "file_name", Message: "file_name is required"}
	}

	_, _, err := s.logDAO.ListLogs(ctx, clientID, "", fname, 1, 10)
	if err != nil {
		s.log.WithError(err).WithFields(logrus.Fields{
			"client_id": clientID,
			"file_name": fname,
		}).Error("Failed to get logs by client")
		return "", err
	}

	return filepath.Join("/data", clientID, fname), nil
}

func (s *LogService) ListLogs(ctx context.Context, args *ListLogsArgs) (logs []models.Log, paging Paginated, err error) {
	if args.Page < 1 {
		args.Page = 1
	}
	if args.PageSize < 1 || args.PageSize > 100 {
		args.PageSize = 20
	}
	var total int64
	logs, total, err = s.logDAO.ListLogs(ctx, args.ClientId, args.UserId, args.FileName, args.Page, args.PageSize)
	if err != nil {
		s.log.WithError(err).WithFields(logrus.Fields{
			"page":      args.Page,
			"page_size": args.PageSize,
		}).Error("Failed to get logs by user")
		return
	}
	paging.Page = int64(args.Page)
	paging.PageSize = int64(args.PageSize)
	paging.Total = total
	paging.TotalPages = (total + int64(args.PageSize) - 1) / int64(args.PageSize)

	s.log.WithFields(logrus.Fields{
		"user_id":   args.UserId,
		"page":      args.Page,
		"page_size": args.PageSize,
		"total":     total,
	}).Info("Logs retrieved successfully by user")
	return
}

/**
 * DeleteOldLogs deletes logs older than specified date
 * @param {context.Context} ctx - Context for request cancellation
 * @param {string} beforeDate - Delete logs before this date
 * @returns {int64, error} Number of deleted records and error if any
 * @description
 * - Validates date parameter
 * - Performs cleanup of old log records
 * - Returns count of deleted records
 * @throws
 * - Validation errors for invalid date
 * - Database deletion errors
 */
func (s *LogService) DeleteOldLogs(ctx context.Context, beforeDate string) (int64, error) {
	// Validate date parameter
	if beforeDate == "" {
		return 0, &ValidationError{Field: "before_date", Message: "before_date is required"}
	}

	// Delete old logs
	count, err := s.logDAO.DeleteOldLogs(ctx, beforeDate)
	if err != nil {
		s.log.WithError(err).WithField("before_date", beforeDate).Error("Failed to delete old logs")
		return 0, err
	}

	s.log.WithFields(logrus.Fields{
		"before_date":   beforeDate,
		"deleted_count": count,
	}).Info("Old logs deleted successfully")

	return count, nil
}

/**
 * validateAndExtractLog validates and extracts log data
 * @param {map[string]interface{}} data - Log data
 * @returns {*models.Log, error} Validated log and error if any
 * @description
 * - Validates required log fields
 * - Extracts log data
 * - Creates log object
 * @throws
 * - Validation errors for missing required fields
 */
func (s *LogService) validate(args *UploadLogArgs) error {
	if args.ClientID == "" {
		return &ValidationError{Field: "client_id", Message: "client_id is required and must be a string"}
	}
	if args.UserID == "" {
		return &ValidationError{Field: "user_id", Message: "user_id is required and must be a string"}
	}
	if args.FileName == "" {
		return &ValidationError{Field: "file_name", Message: "file_name is required and must be a string"}
	}

	return nil
}
