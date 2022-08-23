package clients

import "net/http"

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
