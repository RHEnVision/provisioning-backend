package payloads

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	// root cause
	Err error `json:"-"`
	// HTTP status code
	HTTPStatusCode int `json:"-"`
	// error message
	Message string `json:"msg"`
	// application code
	Code int64 `json:"code,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (e *ErrResponse) Error() error {
	return fmt.Errorf("%s (%d): %w", e.Message, e.Code, e.Err)
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		Message:        "Invalid request",
	}
}

func ErrAWSGeneric(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		Message:        "Error returned from AWS",
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		Message:        "Error rendering response.",
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, Message: "Resource not found"}
var ErrParamParsingError = &ErrResponse{HTTPStatusCode: 500, Message: "Cannot parse parameters"}
var ErrDeleteError = &ErrResponse{HTTPStatusCode: 500, Message: "Cannot delete resource"}
