package dao

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	// ErrNoRows is returned when there are no rows in the result
	ErrNoRows = pgx.ErrNoRows

	// ErrAffectedMismatch is returned when unexpected number of affected rows
	// was returned for INSERT, UPDATE and DELETE queries.
	ErrAffectedMismatch = errors.New("unexpected affected rows")

	// ErrValidation is returned when model does not validate
	ErrValidation = errors.New("validation error")

	// ErrValidation is returned when model transformation fails
	ErrTransformation = errors.New("transformation error")

	// ErrWrongAccount is returned on DAO operations with not matching account id in the context
	ErrWrongAccount = errors.New("wrong account")

	// ErrStubGeneric is a generic error returned for test-related cases
	ErrStubGeneric = errors.New("generic stub error")

	// ErrStubMissingContext is returned when stub object is missing from the context
	ErrStubMissingContext = errors.New("missing variable in context")

	// ErrStubContextAlreadySet is returned when stub object was already added to the context
	ErrStubContextAlreadySet = errors.New("context object already set")
)
