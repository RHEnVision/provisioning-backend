package clients

import (
	"errors"
	"fmt"
)

var (
	// Common errors
	HttpClientErr     = errors.New("HTTP client")
	BadRequestErr     = fmt.Errorf("%w: backend service returned bad request (400)", HttpClientErr)
	NotFoundErr       = fmt.Errorf("%w: backend service returned not found (404) or no data", HttpClientErr)
	UnauthorizedErr   = fmt.Errorf("%w: backend service returned unauthorized (401)", HttpClientErr)
	ForbiddenErr      = fmt.Errorf("%w: backend service returned forbidden (403)", HttpClientErr)
	Non2xxResponseErr = fmt.Errorf("%w: backend service did not return 2xx", HttpClientErr)

	// Sources errors (some others are defined in http package too)
	UnknownAuthenticationTypeErr = errors.New("unknown authentication type")
	UnknownProviderErr           = errors.New("unknown provider type")
	MissingProvisioningSources   = errors.New("missing provisioning source authentication")

	// Client errors
	UnexpectedBackendResponse = errors.New("backend service returned unexpected HTTP code")
	NoResponseData            = errors.New("no data in response")
)
