package migrations

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/migrations/code"
)

type migrationCallback func(ctx context.Context) error

var migrationCallbacks = make(map[int32]migrationCallback)

// Callbacks are executed BEFORE each SQL migration that is matching the migration number. For example,
// a callback with map ID 13 is called before SQL migration 013_xxx.sql.
func init() {
	migrationCallbacks[16] = code.UpdateFingerprints
}

func HasCallback(seq int32) bool {
	_, ok := migrationCallbacks[seq]
	return ok
}

func CallCallback(ctx context.Context, seq int32) error {
	if cb, ok := migrationCallbacks[seq]; ok {
		return cb(ctx)
	}

	return nil
}
