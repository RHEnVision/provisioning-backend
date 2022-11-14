package http

import (
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type SourceError struct {
	err error
}

func (e SourceError) Error() string {
	return e.err.Error()
}

func (e SourceError) Unwrap() error {
	return e.err
}

func WrapSourceError(err error) error {
	return SourceError{err}
}

var (
	MoreThanOneAuthenticationForSourceErr = WrapSourceError(fmt.Errorf("more than one authentication"))
	ApplicationNotFoundErr                = WrapSourceError(fmt.Errorf("application not found is sources app: %w", clients.NotFoundErr))
	ApplicationTypeNotFoundErr            = WrapSourceError(fmt.Errorf("application type 'provisioning' not found: %w", clients.NotFoundErr))
	SourceNotFoundErr                     = WrapSourceError(fmt.Errorf("source not found: %w", clients.NotFoundErr))
	AuthenticationSourceAssociationErr    = WrapSourceError(fmt.Errorf("authentication associated to source id not found"))
	AuthenticationForSourcesNotFoundErr   = WrapSourceError(fmt.Errorf("authentications for source weren't found in sources app: %w", clients.NotFoundErr))
)
