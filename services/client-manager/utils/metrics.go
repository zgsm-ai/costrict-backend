package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

// Global counters for metrics
var (
	requestCount uint64
	errorCount   uint64
	startupTime  time.Time
	timeMutex    sync.Mutex
)

/**
 * IncrementRequestCount increments the global request counter
 * @description
 * - Thread-safe increment of request counter
 * - Used by middleware to track total requests
 * - Supports concurrent access
 */
func IncrementRequestCount() {
	atomic.AddUint64(&requestCount, 1)
}

/**
 * GetRequestCount returns the total number of requests
 * @returns {uint64} Total request count
 * @description
 * - Thread-safe read of request counter
 * - Used for health checks and metrics
 * - Returns current atomic value
 */
func GetRequestCount() uint64 {
	return atomic.LoadUint64(&requestCount)
}

/**
 * IncrementErrorCount increments the global error counter
 * @description
 * - Thread-safe increment of error counter
 * - Used by middleware to track errors
 * - Supports concurrent access
 */
func IncrementErrorCount() {
	atomic.AddUint64(&errorCount, 1)
}

/**
 * GetErrorCount returns the total number of errors
 * @returns {uint64} Total error count
 * @description
 * - Thread-safe read of error counter
 * - Used for health checks and metrics
 * - Returns current atomic value
 */
func GetErrorCount() uint64 {
	return atomic.LoadUint64(&errorCount)
}

/**
 * SetStartupTime sets the application startup time
 * @param {time.Time} t - Startup time
 * @description
 * - Thread-safe setting of startup time
 * - Called during application initialization
 * - Used for uptime calculation
 */
func SetStartupTime(t time.Time) {
	timeMutex.Lock()
	defer timeMutex.Unlock()
	startupTime = t
}

/**
 * GetStartupTime returns the application startup time
 * @returns {time.Time} Startup time
 * @description
 * - Thread-safe read of startup time
 * - Used for uptime calculation
 * - Returns zero time if not set
 */
func GetStartupTime() time.Time {
	timeMutex.Lock()
	defer timeMutex.Unlock()
	return startupTime
}

/**
 * RecordRequestDuration records request duration metrics
 * @param {string} method - HTTP method
 * @param {string} path - Request path
 * @param {int} statusCode - HTTP status code
 * @param {time.Duration} duration - Request duration
 * @description
 * - Records request duration for metrics
 * - Can be extended to store in database or send to monitoring system
 * - Supports performance analysis
 */
func RecordRequestDuration(method, path string, statusCode int, duration time.Duration) {
	// This function can be extended to store metrics in database or send to monitoring system
	// For now, it's a placeholder for future implementation
}

/**
 * GetUptime returns the application uptime
 * @returns {time.Duration} Application uptime
 * @description
 * - Calculates time since startup
 * - Returns 0 if startup time not set
 * - Used for health checks
 */
func GetUptime() time.Duration {
	startup := GetStartupTime()
	if startup.IsZero() {
		return 0
	}
	return time.Since(startup)
}

/**
 * ResetMetrics resets all counters
 * @description
 * - Thread-safe reset of all counters
 * - Used for testing or periodic resets
 * - Should be used with caution in production
 */
func ResetMetrics() {
	atomic.StoreUint64(&requestCount, 0)
	atomic.StoreUint64(&errorCount, 0)

	timeMutex.Lock()
	defer timeMutex.Unlock()
	startupTime = time.Time{}
}

/**
 * GetMetricsSummary returns a summary of all metrics
 * @returns {map[string]interface{}} Summary of metrics
 * @description
 * - Collects current values of all metrics
 * - Returns structured data for reporting
 * - Used for health checks and monitoring
 */
func GetMetricsSummary() map[string]interface{} {
	return map[string]interface{}{
		"requests": GetRequestCount(),
		"errors":   GetErrorCount(),
		"uptime":   GetUptime().String(),
	}
}
