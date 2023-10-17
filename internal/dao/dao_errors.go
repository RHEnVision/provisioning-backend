package dao

import (
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/usrerr"

	"github.com/jackc/pgx/v5"
)

var (
	// ErrNoRows is returned when there are no rows in the result
	// Typically, REST requests should end up with 404 error
	ErrNoRows = pgx.ErrNoRows

	// ErrAffectedMismatch is returned when unexpected number of affected rows
	// was returned for INSERT, UPDATE and DELETE queries.
	// Typically, REST requests should end up with 409 error
	ErrAffectedMismatch = errors.New("unexpected affected rows")

	// ErrValidation is returned when model does not validate
	ErrValidation = usrerr.New(400, "validation error", "invalid input")

	// ErrTransformation is returned when model transformation fails
	ErrTransformation = errors.New("transformation error")

	// ErrWrongAccount is returned on DAO operations with not matching account id in the context
	ErrWrongAccount = usrerr.New(403, "wrong account", "incorrect user account")

	// ErrStubGeneric is a generic error returned for test-related cases
	ErrStubGeneric = errors.New("generic stub error")

	// ErrStubMissingContext is returned when stub object is missing from the context
	ErrStubMissingContext = errors.New("missing variable in context")

	// ErrStubContextAlreadySet is returned when stub object was already added to the context
	ErrStubContextAlreadySet = errors.New("context object already set")

	// ErrReservationRateExceeded is returned when SQL constraint does not allow to insert more reservations
	ErrReservationRateExceeded = usrerr.New(429, "rate limit exceeded", "too many reservations, wait and retry")

	// ErrPubkeyNotFound is returned when a nil pointer to a pubkey is used for reservation detail
	ErrPubkeyNotFound = usrerr.New(404, "pubkey not found", "no pubkey found, it may have been already deleted")
)
