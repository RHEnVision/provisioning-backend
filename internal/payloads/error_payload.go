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
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 400,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
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
		RequestId:      ctxval.RequestId(ctx),
		Context:        ctx,
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
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
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
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
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
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewNotFoundError(ctx context.Context, err error) *ResponseError {
	msg := fmt.Sprintf("not found: %v", err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 404,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewInitializeDAOError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("DAO initialization error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewEnqueueTaskError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("error when enqueing task: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewDAOError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("DAO error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewRenderError(ctx context.Context, message string, err error) *ResponseError {
	msg := fmt.Sprintf("render error: %s: %v", message, err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
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
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
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
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}

func NewUnknownError(ctx context.Context, err error) *ResponseError {
	msg := fmt.Sprintf("unknown error: %v", err)
	if logger := ctxval.Logger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &ResponseError{
		HTTPStatusCode: 500,
		Message:        msg,
		RequestId:      ctxval.RequestId(ctx),
		Err:            err,
		Context:        ctx,
	}
}
