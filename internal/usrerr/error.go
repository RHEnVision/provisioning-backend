package usrerr

import "errors"

// New creates new error with additional information like HTTP code or user-facing message
func New(code int, err, msg string) error {
	//nolint:goerr113
	return &Error{
		StatusCode:  code,
		e:           errors.New(err),
		UserMessage: msg,
	}
}

type Error struct {
	e           error
	StatusCode  int    // HTTP status code
	UserMessage string // User facing optional message
}

func (h Error) Error() string {
	return h.e.Error()
}

func (h Error) Unwrap() error {
	return h.e
}

var (
	ErrBadRequest400   = New(400, "bad request", "")
	ErrNotFound404     = New(404, "not found", "")
	ErrUnauthorized401 = New(401, "unauthorized", "")
	ErrForbidden403    = New(403, "forbidden", "")
)
