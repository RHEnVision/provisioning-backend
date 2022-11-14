package http

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var (
	ApplicationNotFoundErr              = fmt.Errorf("application not found is sources app: %w", clients.NotFoundErr)
	ApplicationTypeNotFoundErr          = fmt.Errorf("application type 'provisioning' not found: %w", clients.NotFoundErr)
	SourceNotFoundErr                   = fmt.Errorf("source not found: %w", clients.NotFoundErr)
	AuthenticationSourceAssociationErr  = errors.New("authentication associated to source id not found")
	AuthenticationForSourcesNotFoundErr = fmt.Errorf("authentications for source weren't found in sources app: %w", clients.NotFoundErr)
	ApplicationReadErr                  = fmt.Errorf("application read returned no application type in sources: %w", clients.NotFoundErr)
)
