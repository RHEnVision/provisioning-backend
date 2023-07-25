package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSleepCtx(t *testing.T) {
	err := sleepCtx(context.Background(), 1*time.Microsecond)
	require.NoError(t, err)
}

func TestSleepCtxDeadline(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer c()
	err := sleepCtx(ctx, 500*time.Microsecond)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestWaitAndRetryCall(t *testing.T) {
	calls := 0
	err := waitAndRetry(context.Background(), func() error {
		calls += 1
		return nil
	}, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestWaitAndRetryErrFirst(t *testing.T) {
	calls := 0
	err := waitAndRetry(context.Background(), func() error {
		calls += 1
		return ErrTryAgain
	}, 1, 1)
	require.ErrorIs(t, err, ErrTryAgain)
	require.Equal(t, 2, calls)
}

func TestWaitAndRetryErrSecond(t *testing.T) {
	calls := 0
	err := waitAndRetry(context.Background(), func() error {
		calls += 1
		if calls <= 1 {
			return ErrTryAgain
		}
		return nil
	}, 1, 1)
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}
