package jobs

import "github.com/RHEnVision/provisioning-backend/pkg/worker"

const (
	TypeNoop                worker.JobType = "no_operation"
	TypeLaunchInstanceAws   worker.JobType = "launch_instances_aws"
	TypeLaunchInstanceAzure worker.JobType = "launch_instances_azure"
	TypeLaunchInstanceGcp   worker.JobType = "launch_instances_gcp"
)
