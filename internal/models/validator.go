package models

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/ssh"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	_ = validate.RegisterValidation("sshPubkey", func(fl validator.FieldLevel) bool {
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(fl.Field().String()))
		return err == nil
	})
}

func Validate(ctx context.Context, model interface{}) validator.ValidationErrors {
	err := validate.StructCtx(ctx, model)
	var validationError validator.ValidationErrors
	if err != nil && !errors.As(err, &validationError) {
		// dependency cycle
		//ctxval.Logger(ctx).Fatal().Msgf("Invalid model passed for validation: %v", model)
		panic(err)
	}
	return validationError
}
