package http

import "errors"

var (
	DuplicatePubkeyErr                    = errors.New("public key already exists in target cloud provider account and region")
	PubkeyNotFoundErr                     = errors.New("pubkey not found in AWS account")
	ServiceAccountUnsupportedOperationErr = errors.New("unsupported operation on service account")
	ARNParsingError                       = errors.New("ARN parsing error")
	NoReservationErr                      = errors.New("no reservation has found in AWS response")
)
