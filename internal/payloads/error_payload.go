package payloads

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"

	httpClients "github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/go-chi/render"
)

// ResponseError is used as a payload for all errors
type ResponseError struct {
	// HTTP status code
	HTTPStatusCode int `json:"-"`

	// user facing error message
	Message string `json:"msg"`

	// trace id from context (if provided)
	TraceId string `json:"trace_id,omitempty"`

	// full root cause
	Error error `json:"error"`
}

func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewInvalidRequestError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("invalid request: %s", message)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewMissingRequestParameterError(ctx context.Context, message string) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
	}
}

func PubkeyDuplicateError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 422,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewUnknownAuthenticationType(ctx context.Context, err error) *ResponseError {
	message := "unknown authentication type"
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func ClientError(ctx context.Context, err error) *ResponseError {
	var status int
	var message string
	if errors.Is(err, clients.NotFoundErr) {
		status = 404
		message = "service returned not found or no data"
	} else if errors.Is(err, clients.UnauthorizedErr) {
		status = 401
		message = "service returned unauthorized"
	} else if errors.Is(err, clients.ForbiddenErr) {
		status = 403
		message = "service returned forbidden"
	} else if errors.Is(err, clients.Non2xxResponseErr) {
		status = 500
		message = "service did not return 2xx"
	} else {
		status = 500
		message = "unknown client error"
	}

	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: status,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func SourcesError(ctx context.Context, err error) *ResponseError {
	var status int
	var message string
	if errors.Is(err, httpClients.MoreThanOneAuthenticationForSourceErr) {
		status = 500
		message = "sources HTTP client: more than one authentication"
	} else if errors.Is(err, httpClients.ApplicationNotFoundErr) {
		status = 404
		message = "sources HTTP client: application not found"
	} else if errors.Is(err, httpClients.ApplicationTypeNotFoundErr) {
		status = 404
		message = "sources HTTP client: unexpected source type"
	} else if errors.Is(err, httpClients.SourceNotFoundErr) {
		status = 404
		message = "sources HTTP client: source not found"
	} else if errors.Is(err, httpClients.AuthenticationSourceAssociationErr) {
		status = 500
		message = "sources HTTP client: authentication associated to source id not found"
	} else if errors.Is(err, httpClients.AuthenticationForSourcesNotFoundErr) {
		status = 404
		message = "sources HTTP client: authentication for source not found"
	} else {
		status = 500
		message = "sources HTTP client: unknown error"
	}
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: status,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewImageBuilderError(ctx context.Context, err error) *ResponseError {
	var status int
	var message string
	if errors.Is(err, httpClients.ComposeNotFoundErr) {
		status = 404
		message = "image builder client: image compose not found"
	} else if errors.Is(err, httpClients.ImageStatusErr) {
		status = 500
		message = "image builder client: build of requested image has not finished yet"
	} else if errors.Is(err, httpClients.UnknownImageTypeErr) {
		status = 500
		message = "image builder client: unknown image type"
	} else if errors.Is(err, httpClients.AmiNotFoundInStatusErr) {
		status = 404
		message = "image builder client: AMI not found in status"
	} else if errors.Is(err, httpClients.NameNotFoundInStatusErr) {
		status = 404
		message = "image builder client: image name not found in image status"
	} else if errors.Is(err, httpClients.IDNotFoundInStatusErr) {
		status = 404
		message = "image builder client: project id not found in image status"
	} else {
		status = 500
		message = "image builder client: unknown error"
	}

	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: status,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func GeneralError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewNotFoundError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("not found: %s", message)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 404,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewEnqueueTaskError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewDAOError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("dao error: %s", message)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewRenderError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewURLParsingError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewStatusError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewAWSError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewAzureError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewGCPError(ctx context.Context, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}
