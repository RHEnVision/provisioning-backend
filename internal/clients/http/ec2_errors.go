package http

import "errors"

var DuplicatePubkeyErr = errors.New("pubkey already exists in the target account")
