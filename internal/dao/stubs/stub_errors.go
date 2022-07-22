package stubs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/go-playground/validator/v10"
)

var ContextReadError = errors.New("missing variable in context")
var ContextSecondInitializationError = errors.New("trying to initialize context twice, please avoid that")
var RSAGenerationError = errors.New("rsa key generation failed")

func NewRecordNotFoundError(ctx context.Context, stubName string) dao.NoRowsError {
	return dao.NoRowsError{
		Message: fmt.Sprintf("%s DAO record does not exist", stubName),
		Context: ctx,
	}
}

func NewCreateError(ctx context.Context, stubName string) dao.Error {
	return dao.Error{
		Message: fmt.Sprintf("create of %s failed", stubName),
		Context: ctx,
	}
}

func newValidationError(ctx context.Context, stubName string, model interface{}, validationErr validator.ValidationErrors) dao.ValidationError {
	errors := []string{fmt.Sprintf("Validation of %s failed: ", stubName)}
	for _, ve := range validationErr {
		errors = append(errors, ve.Error())
	}
	msg := strings.Join(errors, ", ")

	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Info().Msg(msg)
	}
	return dao.ValidationError{Context: ctx, Message: msg, Err: validationErr, Model: model}
}
