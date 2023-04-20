package payloads

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	httpClients "github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

// ResponseError is used as a payload for all errors
type ResponseError struct {
	// HTTP status code
	HTTPStatusCode int `json:"-"`

	// user facing error message
	Message string `json:"msg"`

	// trace id from context (if provided)
	TraceId string `json:"trace_id,omitempty"`

	// edge id from context (if provided)
	EdgeId string `json:"edge_id,omitempty"`

	// full root cause
	Error string `json:"error"`

	// build commit
	Version string `json:"version"`

	// build time
	BuildTime string `json:"build_time"`

	// environment (prod or stage or ephemeral)
	Environment string `json:"environment,omitempty"`
}

func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewResponseError(ctx context.Context, status int, userMsg string, err error) render.Renderer {
	var event *zerolog.Event
	var strError string

	if status < 500 {
		event = ctxval.Logger(ctx).Warn().Stack()
	} else {
		event = ctxval.Logger(ctx).Error().Stack()
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
		TraceId:        ctxval.TraceId(ctx),
		EdgeId:         ctxval.EdgeRequestId(ctx),
		Error:          strError,
		Version:        version.BuildCommit,
		BuildTime:      version.BuildTime,
		Environment:    config.Environment(),
	}
}

func NewInvalidRequestError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Invalid request: %s", message)
	return NewResponseError(ctx, http.StatusBadRequest, message, err)
}

func NewWrongArchitectureUserError(ctx context.Context, err error) render.Renderer {
	return NewResponseError(ctx, http.StatusBadRequest, "Image and type architecture mismatch", err)
}

func NewMissingRequestParameterError(ctx context.Context, message string) render.Renderer {
	return NewResponseError(ctx, http.StatusBadRequest, message, nil)
}

func PubkeyDuplicateError(ctx context.Context, message string, err error) render.Renderer {
	return NewResponseError(ctx, http.StatusUnprocessableEntity, message, err)
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
	} else if errors.Is(err, httpClients.UploadStatusErr) {
		return 404, "could not fetch upload status from image builder"
	}
	return 0, ""
}

func NewClientError(ctx context.Context, err error) render.Renderer {
	if errors.Is(err, clients.UnknownAuthenticationTypeErr) {
		return NewResponseError(ctx, 500, "Unknown authentication type", err)
	}
	if status, message := ImageBuilderHelper(err); status != 0 {
		return NewResponseError(ctx, status, message, err)
	}
	if status, message := SourcesErrorHelper(err); status != 0 {
		return NewResponseError(ctx, status, message, err)
	}
	if status, message := ClientErrorHelper(err); status != 0 {
		return NewResponseError(ctx, status, message, err)
	}
	return NewResponseError(ctx, 500, "HTTP service returned unknown client error", err)
}

func NewNotFoundError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Not found: %s", message)
	return NewResponseError(ctx, http.StatusNotFound, message, err)
}

func NewEnqueueTaskError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Task enqueue error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewDAOError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("DAO error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewRenderError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Rendering error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewURLParsingError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("URL parsing error: %s", message)
	return NewResponseError(ctx, http.StatusBadRequest, message, err)
}

func NewStatusError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Status error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewAWSError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("AWS API error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewAzureError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Azure API error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}

func NewGCPError(ctx context.Context, message string, err error) render.Renderer {
	message = fmt.Sprintf("Google API error: %s", message)
	return NewResponseError(ctx, http.StatusInternalServerError, message, err)
}
