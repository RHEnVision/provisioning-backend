package sources

import (
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var MoreThenOneAuthenticationForSourceErr = errors.New("more then one authentication")
var AuthenticationForSourcesNotFoundErr = errors.New("authentications for source weren't found in sources app")
var ApplicationNotFoundErr = errors.New("application not found is sources app")
var ApplicationTypeNotFoundErr = errors.New("application type 'provisioning' not found")
var SourceNotFoundErr = fmt.Errorf("source not found: %w", clients.NotFoundError)
var AuthenticationSourceAssociationErr = errors.New("authentication associated to source id not found")
