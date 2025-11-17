package util

import (
	"strconv"

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

// ParsePagination reads page/page_size from query params, falling back to defaults if absent/invalid.
type PaginationParams interface {
	Query(key string, defaultValue ...string) string
}

func ParsePagination(q PaginationParams) (int, int) {
	page, pageSize := defaultMinPage, 10
	if v := q.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= defaultMinPage {
			page = p
		}
	}
	if v := q.Query("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps >= defaultMinPageSize {
			pageSize = ps
		}
	}
	return page, pageSize
}
