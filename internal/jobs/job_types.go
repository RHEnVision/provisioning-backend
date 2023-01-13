package jobs

import "github.com/RHEnVision/provisioning-backend/pkg/worker"

const (
	TypeNoop              worker.JobType = "no_operation"
	TypeLaunchInstanceAws worker.JobType = "launch_instances_aws"
	TypeLaunchInstanceGcp worker.JobType = "launch_instances_gcp"
)
