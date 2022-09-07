package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ParseInt64 converts param into int64. If param does not exist, it returns an error.
// TODO: It would be better to move chi.URLParam call out of this function so it can
// be also used for URL params. See below for an examples (MustParseBool/ParseBool).
func ParseInt64(r *http.Request, param string) (int64, error) {
	i, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing URL param '%s' to int64: %w", param, err)
	}
	return i, nil
}

// MustParseBool converts string into bool. If string is empty, it returns an error.
func MustParseBool(str string) (bool, error) {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false, fmt.Errorf("error parsing '%s' to bool: %w", str, err)
	}
	return b, nil
}

// ParseBool converts string into bool. Returns nil when string is empty.
func ParseBool(str string) (*bool, error) {
	if str == "" {
		return nil, nil
	}
	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil, fmt.Errorf("error parsing '%s' to bool: %w", str, err)
	}
	return &b, nil
}
