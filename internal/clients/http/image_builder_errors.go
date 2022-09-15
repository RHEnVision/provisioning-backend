package http

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var (
	ComposeNotFoundErr     = fmt.Errorf("image compose not found: %w", clients.NotFoundErr)
	ImageStatusErr         = errors.New("build of requested image has not finished yet")
	UnknownImageTypeErr    = errors.New("unknown image type")
	AmiNotFoundInStatusErr = errors.New("AMI not found in image status")
)
