package payloads

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"

	"github.com/go-chi/render"
)

// ResponseError implements Go standard error interface as well as Wrapper and Renderer
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

func NewInvalidRequestError(ctx context.Context, err error) *ResponseError {
	msg := "invalid request error"
	if logger := ctxval.Logger(ctx); logger != nil {
		// TODO we should also call .Err(err) to log error
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewMissingRequestParameterError(ctx context.Context, param string) *ResponseError {
	msg := fmt.Sprintf("missing parameter: %s", param)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
	}
}

func PubkeyAlreadyExistsError(ctx context.Context, err error) *ResponseError {
	msg := "pubkey with such name or fingerprint already exists for this account"
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 422,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func ClientError(ctx context.Context, client string, message string, err error, status int) *ResponseError {
	msg := fmt.Sprintf("%v error: %s: %v", client, message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: status,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewClientInitializationError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("HTTP client initialization error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewNotFoundError(ctx context.Context, err error) *ResponseError {
	msg := "not found error"
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 404,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewEnqueueTaskError(ctx context.Context, message string, err error) *ResponseError {
	msg := "enqueue task error"
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewDAOError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("DAO error: %s", message)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewRenderError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("render error: %s", message)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(message)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewURLParsingError(ctx context.Context, paramName string, err error) *ResponseError {
	msg := fmt.Sprintf("URL parsing error for param '%s': %v", paramName, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewStatusError(ctx context.Context, err error) *ResponseError {
	name := "status error"
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(name)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        "status error",
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewAWSError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("AWS error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewAzureError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("Azure error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}

func NewGCPError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("GCP error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		TraceId:        ctxval.TraceId(ctx),
		Error:          err,
	}
}
