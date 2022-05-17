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

func NewInvalidRequestError(ctx context.Context, err error) *ResponseError {
	msg := fmt.Sprintf("invalid request: %v", err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewNotFoundError(ctx context.Context, err error) *ResponseError {
	msg := fmt.Sprintf("not found: %v", err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 404,
		Message:        msg,
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewInitializeDAOError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("DAO initialization error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewDAOError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("DAO error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewRenderError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("render error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewURLParsingError(ctx context.Context, paramName string, err error) *ResponseError {
	msg := fmt.Sprintf("URL parsing error for param '%s': %v", paramName, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        msg,
		RequestId:      ctxval.GetRequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}
