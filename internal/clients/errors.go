package clients

import (
	"errors"
)

var (
	// Common errors
	NotFoundErr       = errors.New("backend service returned not found (404) or no data")
	UnauthorizedErr   = errors.New("backend service returned unauthorized (401)")
	ForbiddenErr      = errors.New("backend service returned forbidden (403)")
	Non2xxResponseErr = errors.New("backend service did not return 2xx")

	// Sources errors (some others are defined in http package too)
	UnknownAuthenticationTypeErr = errors.New("unknown authentication type")
)
