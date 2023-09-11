package stubs

import (
	"errors"
)

var (
	ErrNotImplemented               = errors.New("stub not yet implemented")
	ErrMissingInstanceID            = errors.New("instance id is not present")
	ErrSourceAuthenticationNotFound = errors.New("stubbed authentication for source not found")
	ErrContextRead                  = errors.New("failed to find or convert dao stored in testing context")
)
