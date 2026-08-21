package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode defines application error codes
type ErrorCode string

const (
	// Authentication errors
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeInvalidToken     ErrorCode = "INVALID_TOKEN"
	ErrCodeExpiredToken     ErrorCode = "EXPIRED_TOKEN"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"

	// User errors
	ErrCodeUserNotFound     ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidEmail     ErrorCode = "INVALID_EMAIL"
	ErrCodeInvalidPassword  ErrorCode = "INVALID_PASSWORD"

	// Wallet/Transaction errors
	ErrCodeInsufficientBalance ErrorCode = "INSUFFICIENT_BALANCE"
	ErrCodeInvalidAmount       ErrorCode = "INVALID_AMOUNT"
	ErrCodeTransactionFailed   ErrorCode = "TRANSACTION_FAILED"
	ErrCodeWalletNotFound      ErrorCode = "WALLET_NOT_FOUND"

	// Order errors
	ErrCodeOrderNotFound      ErrorCode = "ORDER_NOT_FOUND"
	ErrCodeProductNotInStock  ErrorCode = "PRODUCT_NOT_IN_STOCK"
	ErrCodeInvalidOrder       ErrorCode = "INVALID_ORDER"
	ErrCodeCannotCancelOrder  ErrorCode = "CANNOT_CANCEL_ORDER"

	// Product errors
	ErrCodeProductNotFound    ErrorCode = "PRODUCT_NOT_FOUND"
	ErrCodeInvalidProduct     ErrorCode = "INVALID_PRODUCT"

	// Validation errors
	ErrCodeValidationFailed   ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidInput       ErrorCode = "INVALID_INPUT"

	// Database errors
	ErrCodeDatabaseError      ErrorCode = "DATABASE_ERROR"
	ErrCodeIntegrityViolation ErrorCode = "INTEGRITY_VIOLATION"

	// Internal errors
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotImplemented     ErrorCode = "NOT_IMPLEMENTED"
)

// AppError represents an application error with context
type AppError struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Details    interface{} `json:"details,omitempty"`
	Err        error       `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewAppErrorWithDetails creates an error with additional details
func NewAppErrorWithDetails(code ErrorCode, message string, statusCode int, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

// NewAppErrorWithErr wraps an existing error
func NewAppErrorWithErr(code ErrorCode, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common error constructors
func UnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

func NotFoundError(message string) *AppError {
	return NewAppError(ErrCodeUserNotFound, message, http.StatusNotFound)
}

func BadRequestError(message string) *AppError {
	return NewAppError(ErrCodeValidationFailed, message, http.StatusBadRequest)
}

func InternalServerError(message string, err error) *AppError {
	return NewAppErrorWithErr(ErrCodeInternalError, message, http.StatusInternalServerError, err)
}

func ConflictError(message string) *AppError {
	return NewAppError(ErrCodeUserAlreadyExists, message, http.StatusConflict)
}

func InsufficientBalanceError(balance string) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeInsufficientBalance,
		"Insufficient wallet balance",
		http.StatusPaymentRequired,
		map[string]string{"balance": balance},
	)
}

func ProductNotInStockError(quantity int) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeProductNotInStock,
		"Product not in stock",
		http.StatusBadRequest,
		map[string]int{"requested": quantity},
	)
}

func ValidationFailedError(message string, details interface{}) *AppError {
	return NewAppErrorWithDetails(ErrCodeValidationFailed, message, http.StatusBadRequest, details)
}
