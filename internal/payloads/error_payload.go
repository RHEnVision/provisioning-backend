package payloads

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	httpClients "github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/version"
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
	Error string `json:"error"`

	// build commit
	Version string `json:"version"`

	// build time
	BuildTime string `json:"build_time"`
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
	}
}

func ClientErrorHelper(err error) (int, string) {
	if errors.Is(err, clients.NotFoundErr) {
		return 404, "service returned not found or no data"
	} else if errors.Is(err, clients.UnauthorizedErr) {
		return 401, "service returned unauthorized"
	} else if errors.Is(err, clients.ForbiddenErr) {
		return 403, "service returned forbidden"
	} else if errors.Is(err, clients.Non2xxResponseErr) {
		return 500, "service did not return 2xx"
	}
	return 0, ""
}

func SourcesErrorHelper(err error) (int, string) {
	if errors.Is(err, httpClients.ApplicationNotFoundErr) {
		return 404, "sources application not found"
	} else if errors.Is(err, httpClients.ApplicationTypeNotFoundErr) {
		return 404, "unexpected source type"
	} else if errors.Is(err, httpClients.SourceNotFoundErr) {
		return 404, "source not found"
	} else if errors.Is(err, httpClients.AuthenticationSourceAssociationErr) {
		return 500, "authentication associated to source id not found"
	} else if errors.Is(err, httpClients.AuthenticationForSourcesNotFoundErr) {
		return 404, "authentication for source not found"
	}
	return 0, ""
}

func ImageBuilderHelper(err error) (int, string) {
	if errors.Is(err, httpClients.ComposeNotFoundErr) {
		return 404, "image builder did not find image compose"
	} else if errors.Is(err, httpClients.ImageStatusErr) {
		return 500, "image builder has not finished the build of requested image"
	} else if errors.Is(err, httpClients.UnknownImageTypeErr) {
		return 500, "unknown image type"
	} else if errors.Is(err, httpClients.AmiNotFoundInStatusErr) {
		return 404, "image builder did not find AMI in status"
	} else if errors.Is(err, httpClients.NameNotFoundInStatusErr) {
		return 404, "image builder did not find image name in image status"
	} else if errors.Is(err, httpClients.IDNotFoundInStatusErr) {
		return 404, "image builder did not find project id in image status"
	}
	return 0, ""
}

func Response(ctx context.Context, status int, message string, err error) *ResponseError {
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: status,
		Message:        message,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
	}
}

func NewClientError(ctx context.Context, err error) *ResponseError {
	if errors.Is(err, clients.UnknownAuthenticationTypeErr) {
		return Response(ctx, 500, "unknown authentication type", err)
	}
	if status, message := ImageBuilderHelper(err); status != 0 {
		return Response(ctx, status, message, err)
	}
	if status, message := SourcesErrorHelper(err); status != 0 {
		return Response(ctx, status, message, err)
	}
	if status, message := ClientErrorHelper(err); status != 0 {
		return Response(ctx, status, message, err)
	}
	return Response(ctx, 500, "HTTP service returned unknown client error", err)
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
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
		Error:          err.Error(),
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
	}
}
