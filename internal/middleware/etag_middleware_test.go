package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateETagFromBuffer(t *testing.T) {
	b := []byte("test")
	etag, _ := GenerateETagFromBuffer("test", time.Minute*1, b)
	assert.Equal(t, "test", etag.Name)
	assert.Equal(t, time.Minute*1, etag.Expiration)
	assert.Equal(t, "fa15fda7c10c75a5", etag.Value)
}

func TestGenerateETagHeader(t *testing.T) {
	b := []byte("test")
	etag, _ := GenerateETagFromBuffer("test", time.Minute*1, b)
	assert.Equal(t, "\"pb-test-fa15fda7c10c75a5\"", etag.Header())
}

func TestGenerateETagCacheControl(t *testing.T) {
	b := []byte("test")
	etag, _ := GenerateETagFromBuffer("test", time.Minute*1, b)
	assert.Equal(t, "max-age=60", etag.CacheControlHeader())
}

func TestGenerateETagFromThreeBuffers(t *testing.T) {
	b1 := []byte("a")
	b2 := []byte("b")
	b3 := []byte("c")
	etag, _ := GenerateETagFromBuffer("test", time.Minute*1, b1, b2, b3)
	assert.Equal(t, "2cd8094a1a277627", etag.Value)
}

func TestGenerateETagFromTwoBuffers(t *testing.T) {
	b1 := []byte("ab")
	b2 := []byte("c")
	etag, _ := GenerateETagFromBuffer("test", time.Minute*1, b1, b2)
	assert.Equal(t, "2cd8094a1a277627", etag.Value)
}
