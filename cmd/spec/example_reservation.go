package main

import (
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
)

var ReservationTime time.Time = MustParseTime("2013-05-13T19:20:25Z")

func MustParseTime(t string) time.Time {
	result, err := time.Parse(time.RFC3339, t)
	if err != nil {
		panic(err)
	}
	return result
}

var GenericReservationResponsePayloadPendingExample = payloads.GenericReservationResponsePayload{
	ID:         1310,
	Provider:   1,
	CreatedAt:  ReservationTime.Add(-10 * time.Second),
	Steps:      3,
	StepTitles: []string{"Ensure public key", "Launch instance(s)", "Fetch instance(s) description"},
	Step:       1,
	Status:     "Started Ensure public key",
	Error:      "",
	FinishedAt: nil,
	Success:    nil,
}

var GenericReservationResponsePayloadSuccessExample = payloads.GenericReservationResponsePayload{
	ID:         1305,
	Provider:   1,
	CreatedAt:  ReservationTime.Add(-10 * time.Second),
	Steps:      3,
	StepTitles: []string{"Ensure public key", "Launch instance(s)", "Fetch instance(s) description"},
	Step:       3,
	Status:     "Finished Fetch instance(s) description",
	Error:      "",
	FinishedAt: ptr.To(ReservationTime),
	Success:    ptr.To(true),
}

var GenericReservationResponsePayloadFailureExample = payloads.GenericReservationResponsePayload{
	ID:         1313,
	Provider:   1,
	CreatedAt:  ReservationTime.Add(-10 * time.Second),
	Steps:      3,
	StepTitles: []string{"Ensure public key", "Launch instance(s)", "Fetch instance(s) description"},
	Step:       2,
	Status:     "Finished Launch instance(s)",
	Error:      "cannot launch ec2 instance: VPCIdNotSpecified: No default VPC for this user. GroupName is only supported for EC2-Classic and default VPC",
	FinishedAt: ptr.To(ReservationTime),
	Success:    ptr.To(false),
}

var GenericReservationResponsePayloadListExample = []payloads.GenericReservationResponsePayload{
	GenericReservationResponsePayloadPendingExample,
	GenericReservationResponsePayloadSuccessExample,
	GenericReservationResponsePayloadFailureExample,
}

var AwsReservationRequestPayloadExample = payloads.AWSReservationRequestPayload{
	PubkeyID:         42,
	SourceID:         "654321",
	Region:           "us-east-1",
	InstanceType:     "t3.small",
	Amount:           1,
	ImageID:          "ami-7846387643232",
	LaunchTemplateID: "",
	Name:             "my-instance",
	PowerOff:         false,
}

var AwsReservationResponsePayloadPendingExample = payloads.AWSReservationResponsePayload{
	PubkeyID:         42,
	SourceID:         "654321",
	Region:           "us-east-1",
	InstanceType:     "t3.small",
	Amount:           1,
	ImageID:          "ami-7846387643232",
	LaunchTemplateID: "",
	Name:             "my-instance",
	PowerOff:         false,
}

var AwsReservationResponsePayloadDoneExample = payloads.AWSReservationResponsePayload{
	ID:               1305,
	PubkeyID:         42,
	SourceID:         "654321",
	Region:           "us-east-1",
	InstanceType:     "t3.small",
	Amount:           1,
	ImageID:          "ami-7846387643232",
	LaunchTemplateID: "",
	AWSReservationID: "r-3743243324231",
	Name:             "my-instance",
	PowerOff:         false,
	Instances: []payloads.InstanceResponse{
		{InstanceID: "i-2324343212", Detail: models.ReservationInstanceDetail{
			PublicDNS:  "",
			PublicIPv4: "10.0.0.88",
		}},
	},
}

var AzureReservationRequestPayloadExample = payloads.AzureReservationRequestPayload{
	PubkeyID:     42,
	SourceID:     "654321",
	Location:     "useast",
	InstanceSize: "Basic_A0",
	Amount:       1,
	ImageID:      "composer-api-081fc867-838f-44a5-af03-8b8def808431",
	Name:         "my-instance",
	PowerOff:     false,
}

var AzureReservationResponsePayloadPendingExample = payloads.AzureReservationResponsePayload{
	ID:           1310,
	PubkeyID:     42,
	SourceID:     "654321",
	Location:     "useast",
	InstanceSize: "Basic_A0",
	Amount:       1,
	ImageID:      "composer-api-081fc867-838f-44a5-af03-8b8def808431",
	Name:         "my-instance",
	PowerOff:     false,
	Instances:    nil,
}

var AzureReservationResponsePayloadDoneExample = payloads.AzureReservationResponsePayload{
	ID:           1310,
	PubkeyID:     42,
	SourceID:     "654321",
	Location:     "useast",
	InstanceSize: "Basic_A0",
	Amount:       1,
	ImageID:      "composer-api-081fc867-838f-44a5-af03-8b8def808431",
	Name:         "my-instance",
	PowerOff:     false,
	Instances: []payloads.InstanceResponse{{
		InstanceID: "/subscriptions/4b9d213f-712f-4d17-a483-8a10bbe9df3a/resourceGroups/redhat-deployed/providers/Microsoft.Compute/images/composer-api-92ea98f8-7697-472e-80b1-7454fa0e7fa7",
		Detail: models.ReservationInstanceDetail{
			PublicDNS:  "",
			PublicIPv4: "10.0.0.88",
		},
	}},
}

var GCPReservationRequestPayloadExample = payloads.GCPReservationRequestPayload{
	PubkeyID:           42,
	SourceID:           "654321",
	Zone:               "us-east-4",
	NamePattern:        "my-instance",
	MachineType:        "e2-micro",
	Amount:             1,
	ImageID:            "08a48fed-de87-40ab-a571-f64e30bd0aa8",
	LaunchTemplateName: "",
}

var GCPReservationResponsePayloadPendingExample = payloads.GCPReservationResponsePayload{
	ID:                 1305,
	PubkeyID:           42,
	SourceID:           "654321",
	Zone:               "us-east-4",
	MachineType:        "e2-micro",
	Amount:             1,
	NamePattern:        "my-instance",
	ImageID:            "08a48fed-de87-40ab-a571-f64e30bd0aa8",
	LaunchTemplateName: "template-1",
	GCPOperationName:   "operation-1686646674436-5fdff07e43209-66146b7e-f3f65ec5",
	PowerOff:           false,
}

var GCPReservationResponsePayloadDoneExample = payloads.GCPReservationResponsePayload{
	ID:                 1305,
	PubkeyID:           42,
	SourceID:           "654321",
	Zone:               "us-east-4",
	MachineType:        "e2-micro",
	Amount:             1,
	ImageID:            "08a48fed-de87-40ab-a571-f64e30bd0aa8",
	LaunchTemplateName: "template-1",
	NamePattern:        "my-instance",
	GCPOperationName:   "operation-1686646674436-5fdff07e43209-66146b7e-f3f65ec5",
	PowerOff:           false,
	Instances: []payloads.InstanceResponse{
		{InstanceID: "3003942005876582747", Detail: models.ReservationInstanceDetail{
			PublicDNS:  "",
			PublicIPv4: "10.0.0.88",
		}},
	},
}

var NoopReservationResponsePayloadExample = payloads.NoopReservationResponsePayload{
	ID: 1310,
}
