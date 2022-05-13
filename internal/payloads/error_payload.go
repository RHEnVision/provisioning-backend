package payloads

import (
	"context"
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"net/http"

	"github.com/go-chi/render"
)

// ResponseError implements Go standard error interface as well as Wrapper and Renderer
type ResponseError struct {
	// HTTP status code
	HTTPStatusCode int `json:"-"`

	// error message including message from wrapped error
	Message string `json:"msg"`

	// request id from context (if provided)
	RequestId string `json:"request_id,omitempty"`

	// root cause
	Err error `json:"-"`

	// context (when provided)
	Context context.Context `json:"-"`
}

func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
}

func (e *ResponseError) Unwrap() error {
	return e.Err
}

func NewInitializeDAOError(ctx context.Context, msg string, err error) *ResponseError {
	format := "DAO initialization error: %s: %v"
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msgf(format, msg, err)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        fmt.Sprintf(format, msg, err),
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewDAOError(ctx context.Context, msg string, err error) *ResponseError {
	format := "DAO error: %s: %v"
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msgf(format, msg, err)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        fmt.Sprintf(format, msg, err),
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewRenderError(ctx context.Context, msg string, err error) *ResponseError {
	format := "render error: %s: %v"
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msgf(format, msg, err)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        fmt.Sprintf(format, msg, err),
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewURLParsingError(ctx context.Context, paramName string, err error) *ResponseError {
	format := "URL parsing error for param '%s': %v"
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msgf(format, paramName, err)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        fmt.Sprintf(format, paramName, err),
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}
