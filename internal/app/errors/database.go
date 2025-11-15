package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ryuyb/fusion/internal/infrastructure/database/ent"
)

func ConvertDatabaseError(err error, resource string) error {
	if err == nil {
		return nil
	}
	if ent.IsNotFound(err) {
		return NotFound(resource).Wrap(err)
	}
	if ent.IsConstraintError(err) {
		return handleConstraintError(err, resource)
	}
	if ent.IsValidationError(err) {
		return ValidationError(err.Error()).Wrap(err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return NotFound(resource).Wrap(err)
	}
	if isConnectionError(err) {
		return &AppError{
			Code:       ErrCodeDatabaseError,
			Message:    "Database connection failed",
			HTTPStatus: 503, // Service Unavailable
			Err:        err,
		}
	}

	return DatabaseError(err)
}

// handleConstraintError 处理约束错误
func handleConstraintError(err error, resource string) *AppError {
	errMsg := err.Error()

	// 唯一约束冲突
	if strings.Contains(errMsg, "unique") || strings.Contains(errMsg, "duplicate") {
		return &AppError{
			Code:       ErrCodeConflict,
			Message:    fmt.Sprintf("%s already exists", resource),
			HTTPStatus: http.StatusConflict,
			Err:        err,
		}
	}

	// 外键约束
	if strings.Contains(errMsg, "foreign key") {
		return &AppError{
			Code:       ErrCodeConflict,
			Message:    "Referenced resource does not exist",
			HTTPStatus: http.StatusConflict,
			Err:        err,
		}
	}

	// 检查约束
	if strings.Contains(errMsg, "check constraint") {
		return ValidationError("Data validation failed").Wrap(err)
	}

	return &AppError{
		Code:       ErrCodeConstraintViolation,
		Message:    "Database constraint violation",
		HTTPStatus: http.StatusConflict,
		Err:        err,
	}
}

// isConnectionError 判断是否为连接错误
func isConnectionError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	keywords := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"no such host",
		"timeout",
	}

	for _, keyword := range keywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}
	return false
}
