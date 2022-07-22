package models

import (
	"context"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/ssh"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	_ = validate.RegisterValidation("sshPubkey", func(fl validator.FieldLevel) bool {
		_, err := ssh.ParsePublicKey(fl.Field().Bytes())
		return err == nil
	})
}

func Validate(ctx context.Context, model interface{}) error {
	return validate.StructCtx(ctx, model)
}
