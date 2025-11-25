package services

/**
 * ValidationError represents a validation error
 * @description
 * - Contains field name and error message
 * - Used for input validation failures
 */
type ValidationError struct {
	Field   string
	Message string
}

/**
 * Error returns the error message
 * @returns {string} Error message
 */
func (e *ValidationError) Error() string {
	return e.Message
}

/**
 * ConflictError represents a conflict error
 * @description
 * - Used for duplicate resource conflicts
 * - Contains error message
 */
type ConflictError struct {
	Message string
}

/**
 * Error returns the error message
 * @returns {string} Error message
 */
func (e *ConflictError) Error() string {
	return e.Message
}

/**
 * NotFoundError represents a not found error
 * @description
 * - Used for resource not found scenarios
 * - Contains error message
 */
type NotFoundError struct {
	Message string
}

/*
*
  - Error returns the error message
  - @returns {string} Error message
*/
func (e *NotFoundError) Error() string {
	return e.Message
}