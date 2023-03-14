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
