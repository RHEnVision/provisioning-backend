//go:build !test

package jq

import "github.com/RHEnVision/provisioning-backend/internal/queue"

func init() {
	queue.GetEnqueuer = getEnqueuer
}
