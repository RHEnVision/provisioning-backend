package logging

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
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
}

var replacement = []byte{'?'}

type SentryReplacer struct {
	buf    []byte
	w      io.Writer
	re     *regexp.Regexp
	mu     sync.Mutex
	closed bool
}

func NewSentryReplacer(w io.Writer) *SentryReplacer {
	sr := SentryReplacer{
		w:  w,
		re: regexp.MustCompile("(" + strings.Join(filters, "|") + ")"),
	}
	return &sr
}

func (sr *SentryReplacer) Write(p []byte) (n int, err error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	if sr.closed {
		return 0, io.EOF
	}

	// Append p to our own buffer, see if there's anything we need to censor.
	sr.buf = append(sr.buf, p...)

	// If we've appended at least a line, censor it and write it out.
	for {
		if len(sr.buf) == 0 {
			// Buffer flushed out completely.
			return len(p), nil
		}

		idx := bytes.IndexRune(sr.buf, '\n')
		if idx < 0 {
			// No line yet, just lie to the caller and tell them we wrote p.
			return len(p), nil
		}

		var line []byte
		line, sr.buf = sr.buf[:idx+1], sr.buf[idx+1:]
		line = sr.re.ReplaceAll(line, replacement)

		_, err := sr.w.Write(line)
		if err != nil {
			// This is not strictly the error related to the incoming `p`, but the best we can do.
			return 0, fmt.Errorf("cannot filter: %w", err)
		}
	}
}

func (sr *SentryReplacer) Close() error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	replaced := sr.re.ReplaceAll(sr.buf, replacement)
	_, err := sr.w.Write(replaced)
	if err != nil {
		return fmt.Errorf("cannot close filter: %w", err)
	}

	sr.closed = true
	return nil
}
