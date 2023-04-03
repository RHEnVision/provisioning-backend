package logging

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewline(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("x\nx\n"))
	require.Equal(t, "x\nx\n", buf.String())
}

func TestNoNewline(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("x\nx"))
	require.Equal(t, "x\n", buf.String())
}

func TestNoNewlineClose(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("x\nx"))
	repl.Close()
	require.Equal(t, "x\nx", buf.String())
}

func TestBufferFlush(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("a\n"))
	_, _ = repl.Write([]byte("\n"))
	require.Equal(t, "a\n\n", buf.String())
	require.Zero(t, len(repl.buf))
}

func TestUUID(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("ca767444-d1f9-11ed-afa1-0242ac120002\n"))
	require.Equal(t, "?\n", buf.String())
}

func TestUUIDSplit(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("ca767444-d1f9-"))
	_, _ = repl.Write([]byte("11ed-afa1-"))
	_, _ = repl.Write([]byte(""))
	_, _ = repl.Write([]byte("0242ac120002\n"))
	require.Equal(t, "?\n", buf.String())
}

func TestARN(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("arn:aws:iam::4328974392798432:role/my-role-123\n"))
	require.Equal(t, "?\n", buf.String())
}

func TestARNSplit(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("arn:aws:iam::"))
	_, _ = repl.Write([]byte("4328974392798432:role/my-role-123\n"))
	require.Equal(t, "?\n", buf.String())
}
