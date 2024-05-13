//go:build test

package stub

import "github.com/RHEnVision/provisioning-backend/internal/queue"

func init() {
	queue.GetEnqueuer = getEnqueuer
}
