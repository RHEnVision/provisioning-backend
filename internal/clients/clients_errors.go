package clients

import "errors"

var NotFoundErr = errors.New("backend service returned not found (404)")
var UnauthorizedErr = errors.New("backend service returned unauthorized (401)")
var ForbiddenErr = errors.New("backend service returned forbidden (403)")
var Non2xxResponseErr = errors.New("backend service did not return 2xx")
