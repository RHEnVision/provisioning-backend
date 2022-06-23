package payloads

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"

	"github.com/go-chi/render"
)

type APIError interface {
	Error() string
	GetStatus() int
	SetMessage(string)
	Unwrap() error
	Render(_ http.ResponseWriter, r *http.Request) error
}

func (e *ResponseError) Error() string       { return fmt.Sprintf("%s: %s", e.Message, e.Err.Error()) }
func (e *ResponseError) GetStatus() int      { return e.HTTPStatusCode }
func (e *ResponseError) SetMessage(m string) { e.Message = m }
func (e *ResponseError) Unwrap() error       { return e.Err }
func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

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

type InternalServerError struct {
	ResponseError
}

func NewInvalidRequestError(ctx context.Context, err error) APIError {
	msg := fmt.Sprintf("invalid request: %v", err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}

	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func New3rdPartyClientError(ctx context.Context, message string, err error) APIError {
	msg := fmt.Sprintf("3rd Party Client error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}

	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func NewClientInitializationError(ctx context.Context, message string, err error) APIError {
	msg := fmt.Sprintf("HTTP client initialization error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}

	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

type NotFoundError struct {
	ResponseError
}

func NewNotFoundError(ctx context.Context, err error) APIError {
	msg := fmt.Sprintf("not found: %v", err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Warn().Msg(msg)
	}
	return &NotFoundError{
		ResponseError: ResponseError{
			HTTPStatusCode: 404,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func NewInitializeDAOError(ctx context.Context, message string, err error) APIError {
	msg := fmt.Sprintf("DAO initialization error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func NewDAOError(ctx context.Context, message string, err error) APIError {
	msg := fmt.Sprintf("DAO error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func NewRenderError(ctx context.Context, message string, err error) APIError {
	msg := fmt.Sprintf("render error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

type BadRequestError struct {
	ResponseError
}

func NewURLParsingError(ctx context.Context, paramName string, err error) APIError {
	msg := fmt.Sprintf("URL parsing error for param '%s': %v", paramName, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &BadRequestError{
		ResponseError: ResponseError{
			HTTPStatusCode: 400,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func NewAWSError(ctx context.Context, message string, err error) APIError {
	msg := fmt.Sprintf("AWS error: %s: %v", message, err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}

func NewUnknownError(ctx context.Context, err error) APIError {
	msg := fmt.Sprintf("unknown error: %v", err)
	if logger := ctxval.GetLogger(ctx); logger != nil {
		logger.Error().Msg(msg)
	}
	return &InternalServerError{
		ResponseError: ResponseError{
			HTTPStatusCode: 500,
			Message:        msg,
			RequestId:      ctxval.GetRequestId(ctx),
			Err:            err,
			Context:        ctx,
		},
	}
}
