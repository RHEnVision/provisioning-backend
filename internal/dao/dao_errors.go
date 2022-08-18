package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/go-playground/validator/v10"
)

type NamedForError interface {
	// NameForError returns DAO implementation name that is passed in the error message (e.g. "account").
	NameForError() string
}

// Error represents a common DAO error.
type Error struct {
	Message string
	Context context.Context
	Err     error
}

// ValidationError is returned when validation on model fails
type ValidationError struct {
	Message string
	Context context.Context
	Err     error
	Model   interface{}
}

// TransformationError is returned when model transformation fails
type TransformationError Error

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

var WrongTenantError = errors.New("trying to manipulate data of different tenant")

func (e Error) Error() string {
	return fmt.Sprintf("DAO error: %s: %s", e.Message, e.Err.Error())
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("DAO error: %s: %s", e.Message, e.Err.Error())
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

func (e NoRowsError) Error() string {
	return fmt.Sprintf("DAO no rows returned: %s", e.Message)
}

func (e NoRowsError) Unwrap() error {
	return e.Err
}

func (e MismatchAffectedError) Error() string {
	return fmt.Sprintf("DAO mismatch affected rows: %s", e.Message)
}

func (e TransformationError) Error() string {
	return fmt.Sprintf("DAO error: %s: %s", e.Message, e.Err.Error())
}

func (e TransformationError) Unwrap() error {
	return e.Err
}

func NewValidationError(ctx context.Context, dao NamedForError, model interface{}, validationErr validator.ValidationErrors) ValidationError {
	errors := []string{fmt.Sprintf("Validation of %s failed: ", dao.NameForError())}
	for _, ve := range validationErr {
		errors = append(errors, ve.Error())
	}
	msg := strings.Join(errors, ", ")

	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Info().Msg(msg)
	}
	return ValidationError{Context: ctx, Message: msg, Err: validationErr, Model: model}
}
