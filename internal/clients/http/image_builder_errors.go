package http

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var ComposeNotFoundErr = fmt.Errorf("image compose not found: %w", clients.NotFoundErr)
var ImageStatusErr = errors.New("build of requested image has not finished yet")
var UnknownImageTypeErr = errors.New("unknown image type")
var AmiNotFoundInStatusErr = errors.New("AMI not found in image status")
