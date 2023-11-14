package validation

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrInvalidId = errors.New("invalid id")

func DigitsOnly(id string) error {
	// Checking for out of range error
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidId, err.Error())
	}

	return nil
}
