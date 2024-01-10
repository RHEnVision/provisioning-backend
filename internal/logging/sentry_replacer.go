package logging

import (
	"regexp"
	"strings"

	"github.com/getsentry/sentry-go"
)

var filters = []string{
	// Example (AWS SDK): RequestID: ca767444-d1f9-11ed-afa1-0242ac120002
	`[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}`,
	// Example (AWS SDK): arn:aws:iam::4328974392798432:role/my-role-123
	`arn:aws:[[:word:]]+::\d+:[[:word:]\*-]+/[[:word:]\*-]+`,
	// Example (AWS SDK): i-1234567890abcdef0
	`[a-z]-[0-9a-f]{17}`,
	// Example: 57:d4:13:ff:c0:74:51:50:41:ec:e1:cd:f1:88:b0:61
	`([0-9a-fA-F]{2}[:-]){15}[0-9a-fA-F]{2}`,
	// Example: 192.168.1.100:32453,
	`\d+\.\d+\.\d+\.\d+:\d+`,
	// Example: The resource 'projects/xxx' was not found
	`'projects/[^']*'`,
	// Example: 2023-06-24T19:34:34.2581206+00:00
	`\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.\d+\+\d\d:\d\d`,
}

var replacement = "?"

type SentryReplacer struct {
	re *regexp.Regexp
}

func NewSentryReplacer() *SentryReplacer {
	sr := SentryReplacer{
		re: regexp.MustCompile("(" + strings.Join(filters, "|") + ")"),
	}
	return &sr
}

func (sr *SentryReplacer) Replace(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	event.Message = sr.re.ReplaceAllString(event.Message, replacement)
	return event
}
