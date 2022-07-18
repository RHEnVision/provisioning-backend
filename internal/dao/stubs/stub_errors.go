package stubs

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
)

var ContextReadError = errors.New("missing variable in context")

func NewRecordNotFoundError(ctx context.Context, stubName string) dao.NoRowsError {
	return dao.NoRowsError{
		Message: fmt.Sprintf("%s DAO record does not exist", stubName),
		Context: ctx,
	}
}
