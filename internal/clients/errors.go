package clients

import (
	"github.com/RHEnVision/provisioning-backend/internal/usrerr"
)

var (
	ErrBadRequest   = usrerr.New(401, "backend service returned bad request", "")
	ErrNotFound     = usrerr.New(404, "backend service returned not found or no data", "")
	ErrUnauthorized = usrerr.New(401, "backend service returned unauthorized", "")

	ErrUnknownAuthenticationType  = usrerr.New(500, "unknown authentication type", "sources backend error")
	ErrUnknownProvider            = usrerr.New(500, "unknown provider type", "sources backend error")
	ErrMissingProvisioningSources = usrerr.New(500, "missing provisioning source authentication", "sources backend error")

	ErrUnexpectedBackendResponse = usrerr.New(500, "backend service returned unexpected HTTP code", "unexpected backend data")
	ErrNoResponseData            = usrerr.New(500, "no data in response", "missing backend data")
)
