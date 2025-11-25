package utils

import (
	"time"
)

/**
 * Time utilities for the application
 * @description
 * - Provides time-related utility functions
 * - Supports time formatting and parsing
 * - Handles time zone conversions
 */

/**
 * GetCurrentTimeString returns current time as formatted string
 * @param {string} format - Time format (default: RFC3339)
 * @returns {string} Formatted time string
 * @description
 * - Returns current time in specified format
 * - Defaults to RFC3339 format
 * - Used for logging and timestamps
 */
func GetCurrentTimeString(format ...string) string {
	t := time.Now()
	if len(format) > 0 {
		return t.Format(format[0])
	}
	return t.Format(time.RFC3339)
}

/**
 * ParseTimeString parses a time string into time.Time
 * @param {string} timeStr - Time string to parse
 * @param {string} format - Time format (default: RFC3339)
 * @returns {time.Time, error} Parsed time and error
 * @description
 * - Parses time string according to format
 * - Defaults to RFC3339 format
 * - Returns error if parsing fails
 */
func ParseTimeString(timeStr string, format ...string) (time.Time, error) {
	if len(format) > 0 {
		return time.Parse(format[0], timeStr)
	}
	return time.Parse(time.RFC3339, timeStr)
}

/**
 * GetDurationString returns human-readable duration string
 * @param {time.Duration} duration - Duration to format
 * @returns {string} Human-readable duration string
 * @description
 * - Converts duration to readable format
 * - Handles hours, minutes, seconds
 * - Used for uptime and performance metrics
 */
func GetDurationString(duration time.Duration) string {
	if duration < time.Minute {
		return duration.String()
	}
	if duration < time.Hour {
		return formatDuration(duration, time.Minute, "minute")
	}
	return formatDuration(duration, time.Hour, "hour")
}

/**
 * formatDuration formats duration with specified unit
 * @param {time.Duration} duration - Duration to format
 * @param {time.Duration} unit - Time unit for formatting
 * @param {string} unitName - Name of the unit
 * @returns {string} Formatted duration string
 * @description
 * - Helper function for duration formatting
 * - Handles singular and plural forms
 * - Returns formatted string with unit
 */
func formatDuration(duration time.Duration, unit time.Duration, unitName string) string {
	value := int(duration / unit)
	remaining := duration % unit

	if value == 1 {
		unitName = unitName[:len(unitName)-1] // Remove 's' for singular
	}

	result := ""
	if value > 0 {
		result += string(rune(value+'0')) + " " + unitName
		if remaining > 0 && remaining < time.Minute {
			result += " " + remaining.String()
		}
	}

	return result
}

/**
 * GetTimeRange returns start and end time for a given period
 * @param {string} period - Time period ("today", "week", "month", "year")
 * @returns {time.Time, time.Time} Start and end time
 * @description
 * - Calculates time range for specified period
 * - Returns start of period and current time
 * - Used for date filtering and reporting
 */
func GetTimeRange(period string) (time.Time, time.Time) {
	now := time.Now()
	var start time.Time

	switch period {
	case "today":
		year, month, day := now.Date()
		start = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	case "week":
		// Start of week (Monday)
		weekday := now.Weekday()
		daysSinceMonday := int(weekday - time.Monday)
		if daysSinceMonday < 0 {
			daysSinceMonday += 7
		}
		start = now.AddDate(0, 0, -daysSinceMonday)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		// Default to last 24 hours
		start = now.Add(-24 * time.Hour)
	}

	return start, now
}

/**
 * IsExpired checks if a time has expired
 * @param {time.Time} expiryTime - Time to check
 * @returns {bool} True if expired
 * @description
 * - Checks if given time is in the past
 * - Used for cache expiration and token validation
 */
func IsExpired(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}

/**
 * GetExpiryTime returns expiry time from now
 * @param {time.Duration} duration - Duration from now
 * @returns {time.Time} Expiry time
 * @description
 * - Calculates expiry time from current time
 * - Used for setting cache and token expiry
 */
func GetExpiryTime(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

/**
 * GetMidnight returns midnight time for given date
 * @param {time.Time} t - Input time
 * @returns {time.Time} Midnight time
 * @description
 * - Returns start of day (00:00:00) for given time
 * - Used for daily operations and reporting
 */
func GetMidnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

/**
 * GetEndOfDay returns end of day time for given date
 * @param {time.Time} t - Input time
 * @returns {time.Time} End of day time
 * @description
 * - Returns end of day (23:59:59.999999999) for given time
 * - Used for daily operations and reporting
 */
func GetEndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}
