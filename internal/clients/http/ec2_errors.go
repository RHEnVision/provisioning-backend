package http

import "errors"

var (
	DuplicatePubkeyErr                    = errors.New("pubkey already exists in the target account")
	ServiceAccountUnsupportedOperationErr = errors.New("unsupported operation on service account")
)
