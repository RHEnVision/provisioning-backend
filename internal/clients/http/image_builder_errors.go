package http

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type ImageBuilderError struct {
	err error
}

func (e ImageBuilderError) Error() string {
	return e.err.Error()
}

func (e ImageBuilderError) Unwrap() error {
	return e.err
}

func WrapImageBuilderError(err error) error {
	return ImageBuilderError{err}
}

var (
	ComposeNotFoundErr      = ImageBuilderError{fmt.Errorf("image compose not found: %w", clients.NotFoundErr)}
	ImageStatusErr          = ImageBuilderError{errors.New("build of requested image has not finished yet")}
	UnknownImageTypeErr     = ImageBuilderError{errors.New("unknown image type")}
	AmiNotFoundInStatusErr  = ImageBuilderError{fmt.Errorf("AMI not found in image status: %w", clients.NotFoundErr)}
	NameNotFoundInStatusErr = ImageBuilderError{fmt.Errorf("image name not found in image status: %w", clients.NotFoundErr)}
	IDNotFoundInStatusErr   = ImageBuilderError{fmt.Errorf("project id not found in image status: %w", clients.NotFoundErr)}
)
