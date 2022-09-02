package clients

import "errors"

var UnauthorizedErr = errors.New("operation not permitted by client, check the provided account permissions")
var DuplicatePubkeyErr = errors.New("pubkey already exists in the target account")
