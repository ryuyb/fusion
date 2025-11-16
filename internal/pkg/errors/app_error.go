package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ryuyb/fusion/internal/infrastructure/provider/validator"
	"github.com/samber/lo"
)

type ErrorCode string

const (
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"

	ErrCodeDatabaseError       ErrorCode = "DATABASE_ERROR"
	ErrCodeConstraintViolation ErrorCode = "CONSTRAINT_VIOLATION"

	ErrCodeStreamingPlatformError ErrorCode = "STREAMING_PLATFORM_ERROR"
)

type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	HTTPStatus int            `json:"-"`
	Details    map[string]any `json:"details,omitempty"`
	Err        error          `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

func (e *AppError) WithDetail(key string, value any) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

func (e *AppError) WithDetails(details map[string]any) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

func NewAppError(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

func NotFound(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

func ValidationError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

func CustomValidationError(validationErrors []validator.ValidationError) *AppError {
	details := make(map[string]any)
	details["errors"] = validationErrors
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    "validation errors",
		HTTPStatus: http.StatusUnprocessableEntity,
		Details:    details,
	}
}

func Internal(err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func DatabaseError(err error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseError,
		Message:    "Database operation failed",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func StreamingPlatformError(platform string, message string, err error) *AppError {
	if message == "" {
		message = "Streaming platform error"
	}
	return &AppError{
		Code:       ErrCodeStreamingPlatformError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
		Details: map[string]any{
			"platform": platform,
		},
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func GetAppError(err error) *AppError {
	if appErr, ok := lo.ErrorsAs[*AppError](err); ok {
		return appErr
	}
	return nil
}

func IsNotFoundError(err error) bool {
	appErr := GetAppError(err)
	return appErr != nil && appErr.HTTPStatus == http.StatusNotFound
}
