package http

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var MoreThenOneAuthenticationForSourceErr = errors.New("more then one authentication")
var AuthenticationForSourcesNotFoundErr = fmt.Errorf("authentications for source weren't found in sources app: %w", clients.NotFoundErr)
var ApplicationNotFoundErr = fmt.Errorf("application not found is sources app: %w", clients.NotFoundErr)
var ApplicationTypeNotFoundErr = fmt.Errorf("application type 'provisioning' not found: %w", clients.NotFoundErr)
var SourceNotFoundErr = fmt.Errorf("source not found: %w", clients.NotFoundErr)
var AuthenticationSourceAssociationErr = errors.New("authentication associated to source id not found")
