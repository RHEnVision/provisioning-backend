package clients

import "errors"

// TODO rename to *Err and 2xx!!
var NotFoundError = errors.New("backend service returned not found (404)")
var UnauthorizedErr = errors.New("backend service returned unauthorized (401)")
var ForbiddenErr = errors.New("backend service returned forbidden (403)")
var Non200ResponseErr = errors.New("backend service did not return 2xx")

// TODO this should not be here, it's AWS specific
var DuplicatePubkeyErr = errors.New("pubkey already exists in the target account")
