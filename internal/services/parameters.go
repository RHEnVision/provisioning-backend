package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseInt64(r *http.Request, param string) (int64, error) {
	i, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting URL param to int64: %w", err)
	}
	return i, nil
}
