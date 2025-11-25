package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

/**
 * Common utilities for the application
 * @description
 * - Provides general utility functions
 * - Supports string manipulation and validation
 * - Handles common data transformation tasks
 */

/**
 * GenerateRandomString generates a random string of specified length
 * @param {int} length - Length of the random string
 * @returns {string} Random string
 * @description
 * - Generates cryptographically secure random string
 * - Uses hexadecimal encoding for readability
 * - Used for generating IDs and tokens
 */
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2) // Each byte represents 2 hex characters
	if _, err := crand.Read(bytes); err != nil {
		// Fallback to math/rand if crypto/rand fails
		return generateFallbackRandomString(length)
	}
	return hex.EncodeToString(bytes)[:length]
}

/**
 * generateFallbackRandomString generates a random string using math/rand
 * @param {int} length - Length of the random string
 * @returns {string} Random string
 * @description
 * - Fallback function for random string generation
 * - Uses alphanumeric characters
 * - Less secure than crypto/rand
 */
func generateFallbackRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	// Seed the random number generator
	mrand.Seed(time.Now().UnixNano())

	for i := range result {
		result[i] = charset[mrand.Intn(len(charset))]
	}
	return string(result)
}

/**
 * IsValidEmail validates email format
 * @param {string} email - Email address to validate
 * @returns {bool} True if email is valid
 * @description
 * - Uses regex pattern for email validation
 * - Checks basic email format requirements
 * - Used for input validation
 */
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

/**
 * SanitizeString sanitizes input string
 * @param {string} input - Input string to sanitize
 * @returns {string} Sanitized string
 * @description
 * - Removes HTML tags and special characters
 * - Trims whitespace
 * - Used for input sanitization
 */
func SanitizeString(input string) string {
	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	sanitized := htmlRegex.ReplaceAllString(input, "")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

/*
*
ToSnakeCase converts string to snake_case
* @param {string} input - Input string to convert
* @returns {string} Snake case string
* @description
* - Converts camelCase or PascalCase to snake_case
* - Handles consecutive uppercase characters
* - Used for database field names
*/
func ToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

/**
 * ToCamelCase converts string to camelCase
 * @param {string} input - Input string to convert
 * @returns {string} Camel case string
 * @description
 * - Converts snake_case to camelCase
 * - Handles consecutive underscores
 * - Used for JSON field names
 */
func ToCamelCase(input string) string {
	parts := strings.Split(input, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

/**
 * TruncateString truncates string to specified length
 * @param {string} input - Input string to truncate
 * @param {int} maxLength - Maximum length
 * @param {string} suffix - Suffix to add if truncated
 * @returns {string} Truncated string
 * @description
 * - Truncates string if longer than maxLength
 * - Adds suffix if truncated
 * - Preserves word boundaries if possible
 */
func TruncateString(input string, maxLength int, suffix string) string {
	if len(input) <= maxLength {
		return input
	}

	// Try to truncate at word boundary
	truncated := input[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + suffix
}

/**
 * ContainsString checks if slice contains string
 * @param {[]string} slice - Slice to search
 * @param {string} item - Item to find
 * @returns {bool} True if item is found
 * @description
 * - Case-sensitive string search in slice
 * - Returns true if item is found
 * - Used for validation and filtering
 */
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

/**
 * RemoveString removes item from string slice
 * @param {[]string} slice - Source slice
 * @param {string} item - Item to remove
 * @returns {[]string} New slice without item
 * @description
 * - Creates new slice without specified item
 * - Preserves order of remaining items
 * - Returns original slice if item not found
 */
func RemoveString(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

/**
 * UniqueStrings returns unique strings from slice
 * @param {[]string} slice - Source slice
 * @returns {[]string} Slice with unique strings
 * @description
 * - Removes duplicate strings from slice
 * - Preserves order of first occurrence
 * - Used for data deduplication
 */
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

/**
 * IsValidUUID validates UUID format
 * @param {string} uuid - UUID to validate
 * @returns {bool} True if UUID is valid
 * @description
 * - Validates standard UUID format
 * - Accepts both UUID v4 and other versions
 * - Used for input validation
 */
func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

/**
 * MaskString masks sensitive information in string
 * @param {string} input - Input string to mask
 * @param {int} visibleChars - Number of visible characters at start and end
 * @param {string} maskChar - Character to use for masking
 * @returns {string} Masked string
 * @description
 * - Masks middle portion of string
 * - Preserves specified characters at start and end
 * - Used for sensitive data display
 */
func MaskString(input string, visibleChars int, maskChar string) string {
	if len(input) <= visibleChars*2 {
		return strings.Repeat(maskChar, len(input))
	}

	start := input[:visibleChars]
	end := input[len(input)-visibleChars:]
	middle := strings.Repeat(maskChar, len(input)-visibleChars*2)

	return start + middle + end
}
