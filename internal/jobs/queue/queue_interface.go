package queue

import (
	"context"

	"github.com/lzap/dejq"
)

type Enqueuer interface {
	Enqueue(ctx context.Context, jobs ...dejq.PendingJob) error
}

var GetEnqueuer func() Enqueuer
