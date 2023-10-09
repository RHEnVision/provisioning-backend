package main

import (
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
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

var GenericReservationResponsePayloadPendingExample = payloads.GenericReservationResponse{
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

var GenericReservationResponsePayloadSuccessExample = payloads.GenericReservationResponse{
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

var GenericReservationResponsePayloadFailureExample = payloads.GenericReservationResponse{
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

var GenericReservationResponsePayloadListExample = payloads.GenericReservationListResponse{
	Data: []*payloads.GenericReservationResponse{
		&GenericReservationResponsePayloadPendingExample,
		&GenericReservationResponsePayloadSuccessExample,
		&GenericReservationResponsePayloadFailureExample,
	},
	Metadata: page.Metadata{
		Total: 3,
		Links: page.Links{
			Previous: "",
			Next:     "",
		},
	},
}

var AwsReservationRequestPayloadExample = payloads.AWSReservationRequest{
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

var AwsReservationResponsePayloadPendingExample = payloads.AWSReservationResponse{
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

var AwsReservationResponsePayloadDoneExample = payloads.AWSReservationResponse{
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
			PublicDNS:   "ec2-184-73-141-211.compute-1.amazonaws.com",
			PublicIPv4:  "184.73.141.211",
			PrivateIPv4: "172.31.36.10",
			PrivateIPv6: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		}},
	},
}

var AzureReservationRequestPayloadExample = payloads.AzureReservationRequest{
	PubkeyID:      42,
	SourceID:      "654321",
	Location:      "useast_1",
	ResourceGroup: "redhat-hcc",
	InstanceSize:  "Basic_A0",
	Amount:        1,
	ImageID:       "composer-api-081fc867-838f-44a5-af03-8b8def808431",
	Name:          "my-instance",
	PowerOff:      false,
}

var AzureReservationResponsePayloadPendingExample = payloads.AzureReservationResponse{
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

var AzureReservationResponsePayloadDoneExample = payloads.AzureReservationResponse{
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
			PublicDNS:   "",
			PublicIPv4:  "10.0.0.88",
			PrivateIPv4: "172.22.0.1",
		},
	}},
}

var GCPReservationRequestPayloadExample = payloads.GCPReservationRequest{
	PubkeyID:         42,
	SourceID:         "654321",
	Zone:             "us-east-4",
	NamePattern:      "my-instance",
	MachineType:      "e2-micro",
	Amount:           1,
	ImageID:          "08a48fed-de87-40ab-a571-f64e30bd0aa8",
	LaunchTemplateID: "",
}

var GCPReservationResponsePayloadPendingExample = payloads.GCPReservationResponse{
	ID:               1305,
	PubkeyID:         42,
	SourceID:         "654321",
	Zone:             "us-east-4",
	MachineType:      "e2-micro",
	Amount:           1,
	NamePattern:      "my-instance",
	ImageID:          "08a48fed-de87-40ab-a571-f64e30bd0aa8",
	LaunchTemplateID: "4883371230199373111",
	GCPOperationName: "operation-1686646674436-5fdff07e43209-66146b7e-f3f65ec5",
	PowerOff:         false,
}

var GCPReservationResponsePayloadDoneExample = payloads.GCPReservationResponse{
	ID:               1305,
	PubkeyID:         42,
	SourceID:         "654321",
	Zone:             "us-east-4",
	MachineType:      "e2-micro",
	Amount:           1,
	ImageID:          "08a48fed-de87-40ab-a571-f64e30bd0aa8",
	LaunchTemplateID: "4883371230199373111",
	NamePattern:      "my-instance",
	GCPOperationName: "operation-1686646674436-5fdff07e43209-66146b7e-f3f65ec5",
	PowerOff:         false,
	Instances: []payloads.InstanceResponse{
		{InstanceID: "3003942005876582747", Detail: models.ReservationInstanceDetail{
			PublicDNS:  "",
			PublicIPv4: "10.0.0.88",
		}},
	},
}

var NoopReservationResponsePayloadExample = payloads.NoopReservationResponse{
	ID: 1310,
}
