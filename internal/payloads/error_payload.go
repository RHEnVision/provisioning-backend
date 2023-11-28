package payloads

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/usrerr"
	"github.com/aws/smithy-go"

	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ResponseError is used as a payload for all errors. Use NewResponseError function
// to create new type to set some fields correctly.
type ResponseError struct {
	// HTTP status code
	HTTPStatusCode int `json:"-" yaml:"-"`

	// user facing error message
	Message string `json:"msg,omitempty" yaml:"msg,omitempty"`

	// trace id from context (if provided)
	TraceId string `json:"trace_id,omitempty" yaml:"trace_id"`

	// edge id from context (if provided)
	EdgeId string `json:"edge_id,omitempty" yaml:"edge_id"`

	// full root cause
	Error string `json:"error" yaml:"error"`

	// build commit
	Version string `json:"version" yaml:"version"`

	// build time
	BuildTime string `json:"build_time" yaml:"build_time"`

	// environment (prod or stage or ephemeral)
	Environment string `json:"environment,omitempty" yaml:"environment"`
}

func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewResponseError(ctx context.Context, status int, userMsg string, err error) *ResponseError {
	var event *zerolog.Event
	var strError string

	if status < 500 {
		event = zerolog.Ctx(ctx).Warn().Stack()
	} else {
		event = zerolog.Ctx(ctx).Error().Stack()
	}
	if err != nil {
		event = event.Err(err)
		strError = err.Error()
	}
	if userMsg == "" {
		// take only part up to the first colon to avoid unique ids (UUIDs, database IDs etc)
		userMsg = strings.SplitN(err.Error(), ":", 2)[0]
	}
	event.Msg(userMsg)

	return &ResponseError{
		HTTPStatusCode: status,
		Message:        userMsg,
		TraceId:        logging.TraceId(ctx),
		Error:          strError,
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
	}
}

func NewInvalidRequestError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Invalid request: %s", message)
	return NewResponseError(ctx, http.StatusBadRequest, message, err)
}

func NewWrongArchitectureUserError(ctx context.Context, err error) *ResponseError {
	return NewResponseError(ctx, http.StatusBadRequest, "Image and type architecture mismatch", err)
}

func NewMissingRequestParameterError(ctx context.Context, message string) *ResponseError {
	return NewResponseError(ctx, http.StatusBadRequest, message, nil)
}

func PubkeyDuplicateError(ctx context.Context, message string, err error) *ResponseError {
	return NewResponseError(ctx, http.StatusUnprocessableEntity, message, err)
}

type userPayload struct {
	code    int
	message string
}

func findUserPayload(err error) *userPayload {
	if err == nil {
		return nil
	}

	var ue *usrerr.Error
	if errors.As(err, &ue) {
		msg := ue.UserMessage
		if msg == "" {
			msg = ue.Error()
		}

		return &userPayload{
			code:    ue.StatusCode,
			message: msg,
		}
	}

	return findUserPayload(errors.Unwrap(err))
}

func findUserResponse(ctx context.Context, logMsg string, err error) *ResponseError {
	if payload := findUserPayload(err); payload != nil {
		logger := log.Ctx(ctx).Warn()
		if payload.code >= 500 {
			logger = log.Ctx(ctx).Error()
		}
		logger.Msgf("%s: %s", logMsg, err)
		return NewResponseError(ctx, payload.code, payload.message, err)
	}

	return nil
}

func NewClientError(ctx context.Context, err error) *ResponseError {
	if errors.Is(err, clients.ErrNotFound) {
		return NewNotFoundError(ctx, "unable to get authentication for sources", err)
	}
	if errors.Is(err, clients.ErrBadRequest) {
		return NewResponseError(ctx, http.StatusBadRequest, "unable to get authentication from sources", err)
	}
	if response := findUserResponse(ctx, "Client error", err); response != nil {
		return response
	}
	log.Ctx(ctx).Error().Msgf("Unknown client error: %s", err)
	return NewResponseError(ctx, 500, "backend client error", err)
}

func NewNotFoundError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Not found: %s", message)
	return NewResponseError(ctx, http.StatusNotFound, message, err)
}

func NewMissingIdentityError(ctx context.Context, msg string, err error) *ResponseError {
	message := fmt.Sprintf("Access denied: %s", msg)
	return NewResponseError(ctx, http.StatusForbidden, message, err)
}

func NewMissingPermissionError(ctx context.Context, resource, permission string, err error) *ResponseError {
	message := fmt.Sprintf("Access denied, missing permission %s on %s", permission, resource)
	return NewResponseError(ctx, http.StatusForbidden, message, err)
}

func NewEnqueueTaskError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Task enqueue error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewDAOError(ctx context.Context, message string, err error) *ResponseError {
	if response := findUserResponse(ctx, "DAO error", err); response != nil {
		return response
	}
	message = fmt.Sprintf("DAO error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewRenderError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Rendering error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewURLParsingError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("URL parsing error: %s", message)
	return NewResponseError(ctx, http.StatusBadRequest, message, err)
}

func NewStatusError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Status error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewAWSError(ctx context.Context, message string, err error) *ResponseError {
	var awsAPIErr *smithy.GenericAPIError
	if errors.As(err, &awsAPIErr) && awsAPIErr.Code == "AccessDenied" {
		message = fmt.Sprintf("AWS assume role failed: %s", message)
		return NewResponseError(ctx, http.StatusForbidden, message, err)
	}
	message = fmt.Sprintf("AWS API error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewAzureError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Azure API error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewGCPError(ctx context.Context, message string, err error) *ResponseError {
	message = fmt.Sprintf("Google API error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}
