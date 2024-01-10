package logging

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

func TestNewline(t *testing.T) {
	evt := sentry.Event{Message: "x\nx\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "x\nx\n", result.Message)
}

func TestNoNewlineClose(t *testing.T) {
	evt := sentry.Event{Message: "x\nx"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "x\nx", result.Message)
}

func TestUUID(t *testing.T) {
	evt := sentry.Event{Message: "ca767444-d1f9-11ed-afa1-0242ac120002\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "?\n", result.Message)
}

func TestARN(t *testing.T) {
	evt := sentry.Event{Message: "arn:aws:iam::4328974392798432:role/my-role-123\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "?\n", result.Message)
}

func TestIPv4(t *testing.T) {
	evt := sentry.Event{Message: "read tcp 10.128.24.14:42094->10.0.217.126:6379: i/o timeout\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "read tcp ?->?: i/o timeout\n", result.Message)
}

func TestFingerprint(t *testing.T) {
	evt := sentry.Event{Message: "pubkey with fingerprint (57:d4:13:ff:c0:74:51:50:41:ec:e1:cd:f1:88:b0:61)\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "pubkey with fingerprint (?)\n", result.Message)
}

func TestAWSResourceID(t *testing.T) {
	evt := sentry.Event{Message: "instance ID 'i-0fe8a8adc1403f5b1' does not exist\n\n\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "instance ID '?' does not exist\n\n\n", result.Message)
}

func TestGoogleProject(t *testing.T) {
	evt := sentry.Event{Message: "The resource 'projects/xxx' was not found\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "The resource ? was not found\n", result.Message)
}

func TestAzureTime(t *testing.T) {
	evt := sentry.Event{Message: "'start time': '2023-06-24T19:34:34.2581206+00:00'\n"}
	repl := NewSentryReplacer()
	result := repl.Replace(&evt, nil)
	require.Equal(t, "'start time': '?'\n", result.Message)
}
