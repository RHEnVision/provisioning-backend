package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Test_PatternLogger(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.Background()

	n := chi.NewRouter()
	m := NewPatternMiddleware("patternOnlyTest")
	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	n.Get(`/users/{firstName}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})

	req1, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/users/JoeBob", nil)
	if err != nil {
		t.Error(err)
	}
	req3, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/users/Misty", nil)
	if err != nil {
		t.Error(err)
	}
	req4, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/metrics", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req1)
	n.ServeHTTP(recorder, req2)
	n.ServeHTTP(recorder, req3)
	n.ServeHTTP(recorder, req4)

	body := recorder.Body.String()

	if !strings.Contains(body, metricNameHttpRequestTotal) {
		t.Errorf("body does not contain request total entry '%s'", metricNameHttpRequestTotal)
	}
	if !strings.Contains(body, metricNameHttpRequestDuration) {
		t.Errorf("body does not contain request duration entry '%s'", metricNameHttpRequestDuration)
	}

	req1Count := `http_request_total{code="OK",method="GET",path="/ok",service="patternOnlyTest"} 1`
	joeBobCount := `http_request_total{code="OK",method="GET",path="/users/JoeBob",service="patternOnlyTest"} 1`
	mistyCount := `http_request_total{code="OK",method="GET",path="/users/Misty",service="patternOnlyTest"} 1`
	firstNamePatternCount := `http_request_total{code="OK",method="GET",path="/users/{firstName}",service="patternOnlyTest"} 2`

	if !strings.Contains(body, req1Count) {
		t.Errorf("body does not contain req1 count summary '%s'", req1Count)
	}
	if strings.Contains(body, joeBobCount) {
		t.Errorf("body should not contain Joe Bob count summary '%s'", joeBobCount)
	}
	if strings.Contains(body, mistyCount) {
		t.Errorf("body should not contain Misty count summary '%s'", mistyCount)
	}
	if !strings.Contains(body, firstNamePatternCount) {
		t.Errorf("body does not contain first name pattern count summary '%s'", firstNamePatternCount)
	}
}
