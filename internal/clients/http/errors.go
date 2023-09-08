package http

import (
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/usrerr"
)

// Sources
var (
	ApplicationNotFoundErr              = usrerr.New(404, "application not found is sources", "")
	ApplicationTypeNotFoundErr          = usrerr.New(404, "application type 'provisioning' not found in sources", "")
	SourceNotFoundErr                   = usrerr.New(404, "source not found", "")
	AuthenticationSourceAssociationErr  = usrerr.New(404, "authentication associated to source id not found", "")
	AuthenticationForSourcesNotFoundErr = usrerr.New(404, "authentications for source weren't found in sources", "")
	ApplicationReadErr                  = usrerr.New(500, "application read returned no application type in sources", "")
	SourceTypeNameNotFoundErr           = usrerr.New(404, "source type name not found", "")
	NotEvenErr                          = errors.New("number of keys and values is not even when building a query")
	SourcesInvalidAuthentication        = usrerr.New(400, "insufficient data for authentication", "")
)

// Image Builder
var (
	CloneNotFoundErr        = usrerr.New(404, "image clone not found", "")
	ComposeNotFoundErr      = usrerr.New(404, "image compose not found", "")
	ImageStatusErr          = usrerr.New(500, "build of requested image has not finished yet", "image still building")
	UnknownImageTypeErr     = usrerr.New(500, "unknown image type", "")
	UploadStatusErr         = usrerr.New(500, "cannot get image status", "")
	ImageRequestNotFoundErr = usrerr.New(500, "image compose request not found", "")
)

// EC2
var (
	DuplicatePubkeyErr             = usrerr.New(406, "public key already exists in target cloud provider account and region", "")
	PubkeyNotFoundErr              = usrerr.New(404, "pubkey not found in AWS account", "")
	ServiceAccountUnsupportedOpErr = usrerr.New(500, "unsupported operation on service account", "")
	ARNParsingError                = usrerr.New(500, "ARN parsing error", "")
	NoReservationErr               = usrerr.New(404, "no reservation has found in AWS response", "")
)
