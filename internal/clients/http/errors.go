package http

import (
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/usrerr"
)

// Sources
var (
	ErrApplicationNotFound              = usrerr.New(404, "application not found is sources", "")
	ErrApplicationTypeNotFound          = usrerr.New(404, "application type 'provisioning' not found in sources", "")
	ErrSourceNotFound                   = usrerr.New(404, "source not found", "")
	ErrAuthenticationSourceAssociation  = usrerr.New(404, "authentication associated to source id not found", "")
	ErrAuthenticationForSourcesNotFound = usrerr.New(404, "authentications for source weren't found in sources", "")
	ErrApplicationRead                  = usrerr.New(500, "application read returned no application type in sources", "")
	ErrSourceTypeNameNotFound           = usrerr.New(404, "source type name not found", "")
	ErrNotEven                          = errors.New("number of keys and values is not even when building a query")
	ErrSourcesInvalidAuthentication     = usrerr.New(400, "insufficient data for authentication", "")
)

// Image Builder
var (
	ErrCloneNotFound        = usrerr.New(404, "image clone not found", "")
	ErrComposeNotFound      = usrerr.New(404, "image compose not found", "")
	ErrImageStatus          = usrerr.New(500, "build of requested image has not finished yet", "image still building")
	ErrUnknownImageType     = usrerr.New(500, "unknown image type", "")
	ErrUploadStatus         = usrerr.New(500, "cannot get image status", "")
	ErrImageRequestNotFound = usrerr.New(500, "image compose request not found", "")
)

// EC2
var (
	ErrDuplicatePubkey             = usrerr.New(406, "public key already exists in target cloud provider account and region", "")
	ErrPubkeyNotFound              = usrerr.New(404, "pubkey not found in AWS account", "")
	ErrServiceAccountUnsupportedOp = usrerr.New(500, "unsupported operation on service account", "")
	ErrARNParsing                  = usrerr.New(500, "ARN parsing error", "")
	ErrNoReservation               = usrerr.New(404, "no reservation has found in AWS response", "")
)
