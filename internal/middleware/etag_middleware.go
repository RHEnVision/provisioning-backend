package middleware

import (
	"errors"
	"fmt"
	"hash/crc64"
	"net/http"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

// An InstanceTypeExpiration represents the default expiration time for instance types
const InstanceTypeExpiration time.Duration = 4 * time.Hour

// An OpenAPIExpiration represents the default expiration time for the OpenAPI JSON
const OpenAPIExpiration time.Duration = 30 * time.Minute

// An InvalidETagErr is returned to prevent empty tag generation
var InvalidETagErr = errors.New("empty etag provided")

// An ETag represents W3C ETag caching header
type ETag struct {
	Name       string
	Expiration time.Duration
	Value      string
	HashTime   time.Duration
}

type ETagValueFunc func() *ETag

var etags = make([]*ETag, 0)

func (etag *ETag) Header() string {
	return fmt.Sprintf("\"pb-%s-%s\"", etag.Name, etag.Value)
}

func (etag *ETag) CacheControlHeader() string {
	return fmt.Sprintf("max-age=%d", int(etag.Expiration.Seconds()))
}

func ETagMiddleware(etagFunc ETagValueFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			etag := etagFunc()
			if etag.Value == "" {
				panic(InvalidETagErr)
			}
			cc := etag.CacheControlHeader()
			w.Header().Set("ETag", etag.Header())
			w.Header().Set("Cache-Control", cc)
			logger := ctxval.Logger(r.Context()).With().Str("etag", etag.Value).Str("etag_name", etag.Name).Logger()
			logger.Trace().Msgf("Returned etag with Cache-Control '%s'", cc)

			if match := r.Header.Get("If-None-Match"); match != "" {
				if strings.Contains(match, etag.Value) {
					logger.Trace().Msgf("ETag cache hit")
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// GenerateETagFromBuffer calculates etag value from one or more buffers (e.g. embedded files)
func GenerateETagFromBuffer(name string, expiration time.Duration, buffers ...[]byte) (*ETag, error) {
	start := time.Now()
	hash := crc64.New(crc64.MakeTable(crc64.ECMA))
	for _, buffer := range buffers {
		_, err := hash.Write(buffer)
		if err != nil {
			return nil, fmt.Errorf("unable to generate etag from buffer: %w", err)
		}
	}
	etag := &ETag{
		Name:       name,
		Expiration: expiration,
		Value:      fmt.Sprintf("%x", hash.Sum64()),
		HashTime:   time.Since(start),
	}
	etags = append(etags, etag)
	return etag, nil
}

// AllETags returns all ETags for diagnostic purposes
func AllETags() []*ETag {
	return etags
}
