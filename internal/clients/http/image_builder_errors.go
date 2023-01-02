package http

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var (
	CloneNotFoundErr        = fmt.Errorf("image clone not found: %w", clients.NotFoundErr)
	ComposeNotFoundErr      = fmt.Errorf("image compose not found: %w", clients.NotFoundErr)
	ImageStatusErr          = errors.New("build of requested image has not finished yet")
	UnknownImageTypeErr     = errors.New("unknown image type")
	AmiNotFoundInStatusErr  = fmt.Errorf("AMI not found in image status: %w", clients.NotFoundErr)
	NameNotFoundInStatusErr = fmt.Errorf("image name not found in image status: %w", clients.NotFoundErr)
	IDNotFoundInStatusErr   = fmt.Errorf("project id not found in image status: %w", clients.NotFoundErr)
)
