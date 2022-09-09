package http

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

func IsHTTPStatus2xx(status int) bool {
	return status >= 200 && status <= 299
}

func IsHTTPNotFound(status int) bool {
	return status == 404
}

func IsHTTPUnauthorized(status int) bool {
	return status == 401
}

func IsHTTPForbidden(status int) bool {
	return status == 403
}

// HandleHTTPResponses parses HTTP status code and returns on of the errors
// defined in the client package (NotFoundErr, UnauthorizedErr, Non2xxResponseErr)
// or nil when the response was 2xx (200-299 range).
func HandleHTTPResponses(ctx context.Context, statusCode int) error {
	if IsHTTPNotFound(statusCode) {
		return clients.NotFoundErr
	}
	if IsHTTPUnauthorized(statusCode) {
		return clients.UnauthorizedErr
	}
	if IsHTTPForbidden(statusCode) {
		return clients.ForbiddenErr
	}
	if !IsHTTPStatus2xx(statusCode) {
		ctxval.Logger(ctx).Warn().Msgf("Non-200 HTTP response seen: %v", statusCode)
		return clients.Non2xxResponseErr
	}
	return nil
}
