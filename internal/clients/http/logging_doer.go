package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/rs/zerolog"
)

type LoggingDoer struct {
	ctx  context.Context
	log  *zerolog.Logger
	doer HttpRequestDoer
}

func NewLoggingDoer(ctx context.Context, doer HttpRequestDoer) *LoggingDoer {
	client := LoggingDoer{
		ctx:  ctx,
		log:  ctxval.Logger(ctx),
		doer: doer,
	}
	return &client
}

func (c *LoggingDoer) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		// read request data into a byte slice
		requestData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read request data: %w", err)
		}

		// rewind the original request reader
		req.Body = ioutil.NopCloser(bytes.NewReader(requestData))

		// perform logging
		c.log.Trace().Str("method", req.Method).
			Str("url", req.URL.RequestURI()).
			Int64("content_length", req.ContentLength).
			Bool("request_trace", true).
			Msg(bytes.NewBuffer(requestData).String())
	}

	// delegate the request
	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, NewDoerErr(err)
	}

	if resp.Body != nil {
		// read response data into a byte slice
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read response data: %w", err)
		}

		// rewind the original response reader
		resp.Body = ioutil.NopCloser(bytes.NewReader(responseData))

		// perform logging
		c.log.Trace().Str("status", resp.Status).
			Int("status_code", resp.StatusCode).
			Int64("content_length", resp.ContentLength).
			Bool("response_trace", true).
			Msg(bytes.NewBuffer(responseData).String())
	}

	return resp, nil
}
