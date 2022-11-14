package clients

import (
	"errors"
	"fmt"
)

var (
	// Common errors
	ClientErr         = errors.New("HTTP client")
	NotFoundErr       = fmt.Errorf("%w: backend service returned not found (404) or no data", ClientErr)
	UnauthorizedErr   = fmt.Errorf("%w: backend service returned unauthorized (401)", ClientErr)
	ForbiddenErr      = fmt.Errorf("%w: backend service returned forbidden (403)", ClientErr)
	Non2xxResponseErr = fmt.Errorf("%w: backend service did not return 2xx", ClientErr)

	// Sources errors (some others are defined in http package too)
	UnknownAuthenticationTypeErr = errors.New("unknown authentication type")
)
