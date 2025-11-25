package models

import (
	"time"

	"gorm.io/gorm"
)

/**
 * Configuration model represents a configuration item in the system
 * @description
 * - Stores configuration data with namespace and key
 * - Supports caching for performance optimization
 * - Includes metadata for configuration management
 */
type Configuration struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Namespace   string         `json:"namespace" gorm:"index;not null"`
	Key         string         `json:"key" gorm:"index;not null"`
	Value       string         `json:"value" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

/**
 * Feedback model represents user feedback data
 * @description
 * - Stores various types of user feedback
 * - Includes conversation tracking
 * - Supports different feedback types (completion, copy, evaluate, etc.)
 */
type Feedback struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Type           string    `json:"type" gorm:"index;not null"`
	ConversationID string    `json:"conversation_id" gorm:"index"`
	UserID         string    `json:"user_id" gorm:"index"`
	Content        string    `json:"content" gorm:"type:text"`
	Metadata       string    `json:"metadata" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

/**
 * Log model represents client log entries
 * @description
 * - Stores log data from clients
 * - Includes client and user identification
 * - Supports structured logging with module information
 */
type Log struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ClientID    string    `json:"client_id" gorm:"index;not null"`
	UserID      string    `json:"user_id" gorm:"index"`
	FileName    string    `json:"file_name" gorm:"index;not null"`
	FirstLineNo int64     `json:"first_line_no"`
	LastLineNo  int64     `json:"end_line_no"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

/**
 * TableName returns the table name for Configuration model
 * @returns {string} Database table name
 */
func (Configuration) TableName() string {
	return "configurations"
}

/**
 * TableName returns the table name for Feedback model
 * @returns {string} Database table name
 */
func (Feedback) TableName() string {
	return "feedbacks"
}

/**
 * TableName returns the table name for Log model
 * @returns {string} Database table name
 */
func (Log) TableName() string {
	return "logs"
}
