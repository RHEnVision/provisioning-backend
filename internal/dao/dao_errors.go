package dao

import (
	"context"
	"fmt"
)

// Error represents a common DAO error.
type Error struct {
	Message string
	Context context.Context
	Err     error
}

// NoRowsError is returned when no rows were returned.
type NoRowsError struct {
	Message string
	Context context.Context
	Err     error
}

// MismatchAffectedError is returned when affected rows do not match expectation (e.g. create/delete).
type MismatchAffectedError struct {
	Message string
	Context context.Context
}

func (e *Error) Error() string {
	return fmt.Sprintf("DAO error: %s: %s", e.Message, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *NoRowsError) Error() string {
	return fmt.Sprintf("DAO no rows returned: %s", e.Message)
}

func (e *NoRowsError) Unwrap() error {
	return e.Err
}

func (e *MismatchAffectedError) Error() string {
	return fmt.Sprintf("DAO mismatch affected rows: %s", e.Message)
}
