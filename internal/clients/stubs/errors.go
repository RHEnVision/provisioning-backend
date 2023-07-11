package stubs

import (
	"errors"
)

var (
	NotImplementedErr            = errors.New("stub not yet implemented")
	MissingInstanceIDErr         = errors.New("instance id is not present")
	SourceAuthenticationNotFound = errors.New("stubbed authentication for source not found")
	ContextReadError             = errors.New("failed to find or convert dao stored in testing context")
)
