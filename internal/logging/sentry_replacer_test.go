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

func TestIPv4(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("read tcp 10.128.24.14:42094->10.0.217.126:6379: i/o timeout\n"))
	require.Equal(t, "read tcp ?->?: i/o timeout\n", buf.String())
}

func TestFingerprint(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("pubkey with fingerprint (57:d4:13:ff:c0:74:51:50:41:ec:e1:cd:f1:88:b0:61)\n"))
	require.Equal(t, "pubkey with fingerprint (?)\n", buf.String())
}

func TestAWSResourceID(t *testing.T) {
	buf := bytes.NewBufferString("")
	repl := NewSentryReplacer(buf)
	_, _ = repl.Write([]byte("instance ID 'i-0fe8a8adc1403f5b1' does not exist\n\n\n"))
	require.Equal(t, "instance ID '?' does not exist\n\n\n", buf.String())
}
