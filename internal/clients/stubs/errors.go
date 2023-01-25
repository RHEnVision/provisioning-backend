package stubs

import (
	"errors"
)

var (
	NotImplementedErr            = errors.New("stub not yet implemented")
	SourceAuthenticationNotFound = errors.New("stubbed authentication for source not found")
)
