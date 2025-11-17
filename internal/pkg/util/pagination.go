package util

import (
	"github.com/ryuyb/fusion/internal/pkg/errors"
)

const (
	defaultMinPage     = 1
	defaultMinPageSize = 1
	defaultMaxPageSize = 200
)

func ValidatePagination(page, pageSize int) error {
	if page < defaultMinPage {
		return errors.BadRequest("page must be greater than zero")
	}
	if pageSize < defaultMinPageSize || pageSize > defaultMaxPageSize {
		return errors.BadRequest("page size must be between 1 and 200")
	}
	return nil
}
