package dao

import (
	"context"
	"fmt"
)

// Error type for all DAO errors.
type Error struct {
	Message string
	Context context.Context
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("DAO error: %s: %s", e.Message, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
