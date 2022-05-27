package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseUint64(r *http.Request, param string) (uint64, error) {
	i, err := strconv.ParseUint(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting URL param to uint64: %w", err)
	}
	return i, nil
}
