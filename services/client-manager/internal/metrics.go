package internal

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"

	"github.com/zgsm-ai/client-manager/utils"
)

// Prometheus metrics
var (
	// HTTP request counter
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP request duration histogram
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP error counter
	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Active connections gauge
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	// Database operations counter
	dbOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "table"},
	)

	// Database operation duration histogram
	dbOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "Database operation duration in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 1, 10},
		},
		[]string{"operation", "table"},
	)

	// Cache operations counter
	cacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "cache"},
	)

	// Cache hit/miss counter
	cacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits and misses",
		},
		[]string{"cache", "result"},
	)

	// Feedback received counter
	feedbackReceivedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "feedback_received_total",
			Help: "Total number of feedback received",
		},
		[]string{"type"},
	)

	// Logs received counter
	logsReceivedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "logs_received_total",
			Help: "Total number of logs received",
		},
		[]string{"client_id", "module"},
	)

	// Configuration access counter
	configurationAccessTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "configuration_access_total",
			Help: "Total number of configuration accesses",
		},
		[]string{"namespace", "key"},
	)
)

/**
 * InitMetrics initializes Prometheus metrics
 * @description
 * - Initializes all Prometheus metrics
 * - Registers metrics with Prometheus registry
 * - Sets default values for gauges
 * @throws
 * - Metrics registration errors
 */
func InitMetrics() {
	// Initialize active connections gauge
	activeConnections.Set(0)

	// Log metrics initialization
	logrus.Info("Prometheus metrics initialized")
}

/**
 * IncrementRequestCount increments the total request counter
 * @description
 * - Increments the global request counter
 * - Updates the active connections gauge
 * - Used by the request middleware
 */
func IncrementRequestCount() {
	// Increment utils counter
	utils.IncrementRequestCount()

	// Increment active connections
	activeConnections.Inc()
}

/**
 * DecrementActiveConnections decrements the active connections gauge
 * @description
 * - Decrements the active connections gauge
 * - Should be called when request processing completes
 */
func DecrementActiveConnections() {
	activeConnections.Dec()
}

/**
 * RecordHTTPRequest records HTTP request metrics
 * @param {string} method - HTTP method
 * @param {string} endpoint - Request endpoint
 * @param {int} statusCode - HTTP status code
 * @param {time.Duration} duration - Request duration
 * @description
 * - Records HTTP request count and duration
 * - Updates both total counter and histogram
 * - Formats status code as string for labels
 */
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration) {
	statusStr := strconv.Itoa(statusCode)

	// Increment request counter
	httpRequestsTotal.WithLabelValues(method, endpoint, statusStr).Inc()

	// Record request duration
	httpRequestDuration.WithLabelValues(method, endpoint, statusStr).Observe(duration.Seconds())

	// Record error if status code indicates error
	if statusCode >= 400 {
		httpErrorsTotal.WithLabelValues(method, endpoint, statusStr).Inc()
	}
}

/**
 * RecordDBOperation records database operation metrics
 * @param {string} operation - Database operation type
 * @param {string} table - Database table name
 * @param {time.Duration} duration - Operation duration
 * @description
 * - Records database operation count and duration
 * - Updates both counter and histogram
 * - Used by DAO layer for performance monitoring
 */
func RecordDBOperation(operation, table string, duration time.Duration) {
	// Increment operation counter
	dbOperationsTotal.WithLabelValues(operation, table).Inc()

	// Record operation duration
	dbOperationDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

/**
 * RecordCacheOperation records cache operation metrics
 * @param {string} operation - Cache operation type
 * @param {string} cache - Cache name
 * @description
 * - Records cache operation count
 * - Updates the cache operations counter
 * - Used by cache layer for monitoring
 */
func RecordCacheOperation(operation, cache string) {
	cacheOperationsTotal.WithLabelValues(operation, cache).Inc()
}

/**
 * RecordCacheHit records cache hit metrics
 * @param {string} cache - Cache name
 * @param {bool} hit - Whether it was a hit or miss
 * @description
 * - Records cache hit/miss count
 * - Updates the cache hits counter
 * - Used for cache effectiveness monitoring
 */
func RecordCacheHit(cache string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	cacheHitsTotal.WithLabelValues(cache, result).Inc()
}

/**
 * RecordFeedbackReceived records feedback received metrics
 * @param {string} feedbackType - Type of feedback
 * @description
 * - Records feedback received count
 * - Updates the feedback counter
 * - Used for feedback analytics
 */
func RecordFeedbackReceived(feedbackType string) {
	feedbackReceivedTotal.WithLabelValues(feedbackType).Inc()
}

/**
 * RecordLogsReceived records logs received metrics
 * @param {string} clientID - Client identifier
 * @param {string} module - Module name
 * @description
 * - Records logs received count
 * - Updates the logs counter
 * - Used for logging analytics
 */
func RecordLogsReceived(clientID, module string) {
	logsReceivedTotal.WithLabelValues(clientID, module).Inc()
}

/**
 * RecordConfigurationAccess records configuration access metrics
 * @param {string} namespace - Configuration namespace
 * @param {string} key - Configuration key
 * @description
 * - Records configuration access count
 * - Updates the configuration access counter
 * - Used for configuration usage analytics
 */
func RecordConfigurationAccess(namespace, key string) {
	configurationAccessTotal.WithLabelValues(namespace, key).Inc()
}

/**
 * GetMetricsSummary returns a summary of all metrics
 * @returns {map[string]interface{}} Summary of metrics
 * @description
 * - Collects current values of all metrics
 * - Returns a structured summary for reporting
 * - Used for health checks and monitoring
 */
func GetMetricsSummary() map[string]interface{} {
	summary := make(map[string]interface{})

	// Add request counts from utils
	summary["requests"] = map[string]interface{}{
		"total":  utils.GetRequestCount(),
		"errors": utils.GetErrorCount(),
	}

	// Get startup time from utils
	startupTime := utils.GetStartupTime()
	if !startupTime.IsZero() {
		summary["uptime"] = time.Since(startupTime).String()
	}

	return summary
}
