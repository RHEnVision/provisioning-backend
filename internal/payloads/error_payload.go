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
	msg := fmt.Sprintf("DAO error: %s", message)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Err(err).Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
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
