// Import this package in all Go tests which need to enqueue a job. The implementation
// is to silently handle all incoming jobs without any effects.
package stub

import "github.com/RHEnVision/provisioning-backend/internal/jobs"

func init() {
	jobs.InitializeStub()
}
