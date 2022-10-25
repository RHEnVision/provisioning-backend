package http

import "net/http"

// DoerErr is a simple wrapped error without any message. Additional message would
// stack for each request as multiple doers are called leading to:
//
// "error in doer1: error in doer2: error in doer3: something happened"
type DoerErr struct {
	Err error
}

func NewDoerErr(err error) *DoerErr {
	return &DoerErr{Err: err}
}

func (e *DoerErr) Error() string {
	return e.Err.Error()
}

func (e *DoerErr) Unwrap() error {
	return e.Err
}

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
