package sqlx

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
)

type NamedForError interface {
	// NameForError returns DAO implementation name that is passed in the error message (e.g. "account").
	NameForError() string
}

func newError(ctx context.Context, msg string, err error) dao.Error {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return dao.Error{
		Message: msg,
		Context: ctx,
		Err:     err,
	}
}

func newMismatchAffectedError(ctx context.Context, msg string) dao.MismatchAffectedError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return dao.MismatchAffectedError{
		Message: msg,
		Context: ctx,
	}
}

func newNoRowsError(ctx context.Context, msg string, noRowsErr error) dao.NoRowsError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Debug().Err(noRowsErr).Msg(msg)
	}
	return dao.NoRowsError{
		Err:     noRowsErr,
		Message: msg,
		Context: ctx,
	}
}

func NewPrepareStatementError(context context.Context, daoName NamedForError, sql string, err error) dao.Error {
	msg := fmt.Sprintf("sqlx %s prepare statement error: %s: %v", daoName.NameForError(), sql, err)
	return newError(context, msg, err)
}

func NewTransactionError(context context.Context, err error) dao.Error {
	msg := fmt.Sprintf("transaction: %v", err)
	return newError(context, msg, err)
}

func NewGetError(context context.Context, daoName NamedForError, sql string, err error) dao.Error {
	msg := fmt.Sprintf("sqlx %s get error: %s: %v", daoName.NameForError(), sql, err)
	return newError(context, msg, err)
}

func NewCreateError(context context.Context, daoName NamedForError, sql string, err error) dao.Error {
	msg := fmt.Sprintf("sqlx %s exec create error: %s: %v", daoName.NameForError(), sql, err)
	return newError(context, msg, err)
}

func NewSelectError(context context.Context, daoName NamedForError, sql string, err error) dao.Error {
	msg := fmt.Sprintf("sqlx %s select error: %s: %v", daoName.NameForError(), sql, err)
	return newError(context, msg, err)
}

func NewExecUpdateError(context context.Context, daoName NamedForError, sql string, err error) dao.Error {
	msg := fmt.Sprintf("sqlx %s exec update error: %s: %v", daoName.NameForError(), sql, err)
	return newError(context, msg, err)
}

func NewExecDeleteError(context context.Context, daoName NamedForError, sql string, err error) dao.Error {
	msg := fmt.Sprintf("sqlx %s exec delete error: %s: %v", daoName.NameForError(), sql, err)
	return newError(context, msg, err)
}

func NewDeleteMismatchAffectedError(context context.Context, daoName NamedForError, expected, was int64) dao.MismatchAffectedError {
	msg := fmt.Sprintf("sqlx %s delete expected: %d rows, was: %d rows", daoName.NameForError(), expected, was)
	return newMismatchAffectedError(context, msg)
}

func NewUpdateMismatchAffectedError(context context.Context, daoName NamedForError, expected, was int64) dao.MismatchAffectedError {
	msg := fmt.Sprintf("sqlx %s update expected: %d rows, was: %d rows", daoName.NameForError(), expected, was)
	return newMismatchAffectedError(context, msg)
}

func NewNoRowsError(context context.Context, daoName NamedForError, sql string, noRowsErr error) dao.NoRowsError {
	msg := fmt.Sprintf("sqlx %s no rows returned from: %s", daoName.NameForError(), sql)
	return newNoRowsError(context, msg, noRowsErr)
}
