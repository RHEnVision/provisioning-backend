package sqlx

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
)

func NewPrepareStatementError(context context.Context, sql string, err error) *dao.Error {
	if logger := ctxval.GetLogger(context); logger != nil {
		logger.Error().Msgf("sqlx prepare statement error: %s: %v", sql, err)
	}
	return &dao.Error{
		Message: sql,
		Context: context,
		Err:     err,
	}
}

func NewGetError(context context.Context, msg string, err error) *dao.Error {
	if logger := ctxval.GetLogger(context); logger != nil {
		logger.Error().Msgf("sqlx get error: %s: %v", msg, err)
	}
	return &dao.Error{
		Message: msg,
		Context: context,
		Err:     err,
	}
}
