package queue

import (
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

var GetEnqueuer func() worker.JobEnqueuer
