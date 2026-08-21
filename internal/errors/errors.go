package errors

import "fmt"

// APIError represents a standard API error
type APIError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status_code"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code string, message string, statusCode int) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// Common error codes
const (
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeConflict           = "CONFLICT"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeUserNotFound       = "USER_NOT_FOUND"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeProductNotFound    = "PRODUCT_NOT_FOUND"
	ErrCodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrCodeOrderNotFound      = "ORDER_NOT_FOUND"
	ErrCodeSessionExpired     = "SESSION_EXPIRED"
)

// Predefined errors
var (
	ErrUnauthorized = NewAPIError(ErrCodeUnauthorized, "Authentication required", 401)
	ErrForbidden    = NewAPIError(ErrCodeForbidden, "Access denied", 403)
	ErrNotFound     = NewAPIError(ErrCodeNotFound, "Resource not found", 404)
	ErrInternalServer = NewAPIError(ErrCodeInternalError, "Internal server error", 500)
)
