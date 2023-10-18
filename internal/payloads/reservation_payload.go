package payloads

import (
	"net/http"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/go-chi/render"
)

// ReservationRequest is empty, account comes in HTTP header and
// provider type in HTTP URL. All other fields are auto-generated.

type GenericReservationResponse struct {
	ID int64 `json:"id" yaml:"id"`

	// Provider type. Required.
	Provider int `json:"provider" yaml:"provider"`

	// Time when reservation was made.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// Total number of job steps for this reservation.
	Steps int32 `json:"steps" yaml:"steps"`

	// User-facing step descriptions for each step. Length of StepTitles must be equal to Steps.
	StepTitles []string `json:"step_titles" yaml:"step_titles"`

	// Active job step for this reservation. See Status for more details.
	Step int32 `json:"step" yaml:"step"`

	// Textual status of the reservation or error when there was a failure
	Status string `json:"status" yaml:"status"`

	// Error message when reservation was not successful. Only set when Success if false.
	Error string `json:"error" yaml:"error"`

	// Time when reservation was finished or nil when it's still processing.
	FinishedAt *time.Time `json:"finished_at" nullable:"true" yaml:"finished_at"`

	// Flag indicating success, error or unknown state (NULL). See Status for the actual error.
	Success *bool `json:"success" nullable:"true" yaml:"success"`
}

type InstanceResponse struct {
	// Instance ID which has been created on a cloud provider.
	InstanceID string `json:"instance_id" yaml:"instance_id"`

	// Instance's description, ip and dns
	Detail models.ReservationInstanceDetail `json:"detail" yaml:"detail"`
}

type AWSReservationResponse struct {
	ID int64 `json:"reservation_id" yaml:"reservation_id"`

	// Pubkey ID.
	PubkeyID *int64 `json:"pubkey_id,omitempty" yaml:"pubkey_id,omitempty"`

	// Source ID.
	SourceID string `json:"source_id" yaml:"source_id"`

	// AWS region.
	Region string `json:"region" yaml:"region"`

	// AWS Instance type.
	InstanceType string `json:"instance_type" yaml:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 `json:"amount" yaml:"amount"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id" yaml:"image_id"`

	// Optional launch template ID ("lt-9848392734432") or empty for no template.
	LaunchTemplateID string `json:"launch_template_id" yaml:"launch_template_id"`

	// The ID of the aws reservation which was created, or missing if not created yet.
	AWSReservationID string `json:"aws_reservation_id,omitempty" yaml:"aws_reservation_id"`

	// Optional name of the instance(s).
	Name string `json:"name" yaml:"name"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff" yaml:"poweroff"`

	// Instances array, only present for finished reservations
	Instances []InstanceResponse `json:"instances,omitempty" yaml:"instances"`
}

type AzureReservationResponse struct {
	ID int64 `json:"reservation_id" yaml:"reservation_id"`

	PubkeyID *int64 `json:"pubkey_id,omitempty" yaml:"pubkey_id,omitempty"`

	SourceID string `json:"source_id" yaml:"source_id"`

	ResourceGroup string `json:"resource_group" yaml:"resource_group"`

	// Azure Location.
	Location string `json:"location" yaml:"location"`

	// Azure Instance size.
	InstanceSize string `json:"instance_size" yaml:"instance_size"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 `json:"amount" yaml:"amount"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id" yaml:"image_id"`

	// Optional name of the instance(s).
	Name string `json:"name" yaml:"name"`

	// Immediately PowerOff the system after initialization.
	PowerOff bool `json:"poweroff" yaml:"poweroff"`

	// Instances IDs, only present for finished reservations.
	Instances []InstanceResponse `json:"instances,omitempty" yaml:"instances"`
}

type GCPReservationResponse struct {
	ID int64 `json:"reservation_id" yaml:"reservation_id"`

	// Pubkey ID.
	PubkeyID *int64 `json:"pubkey_id,omitempty" yaml:"pubkey_id,omitempty"`

	// Source ID.
	SourceID string `json:"source_id" yaml:"source_id"`

	// GCP zone.
	Zone string `json:"zone" yaml:"zone"`

	// GCP Machine type.
	MachineType string `json:"machine_type" yaml:"machine_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 `json:"amount" yaml:"amount"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id" yaml:"image_id"`

	// Optional name pattern of the instance(s).
	NamePattern string `json:"name_pattern" yaml:"name_pattern"`

	// The name of the gcp operation which was created.
	GCPOperationName string `json:"gcp_operation_name,omitempty" yaml:"gcp_operation_name"`

	// Optional launch template id global/instanceTemplates/ID or empty string
	LaunchTemplateID string `json:"launch_template_id,omitempty" yaml:"launch_template_id"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff" yaml:"poweroff"`

	// Instances IDs, only present for finished reservations.
	Instances []InstanceResponse `json:"instances,omitempty" yaml:"instances"`
}

type NoopReservationResponse struct {
	ID int64 `json:"reservation_id" yaml:"reservation_id"`
}

type AWSReservationRequest struct {
	// Pubkey ID. Always required even when launch template provides one.
	PubkeyID int64 `json:"pubkey_id" yaml:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id" yaml:"source_id"`

	// AWS region.
	Region string `json:"region" yaml:"region"`

	// Optional name of the instance(s).
	Name string `json:"name" yaml:"name"`

	// Optional launch template ID ("lt-9848392734432") or empty for no template.
	LaunchTemplateID string `json:"launch_template_id,omitempty" yaml:"launch_template_id"`

	// AWS Instance type.
	InstanceType string `json:"instance_type" yaml:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 ` json:"amount" yaml:"amount"`

	// Image Builder UUID of the image that should be launched. AMI's must be prefixed with 'ami-'.
	ImageID string `json:"image_id" yaml:"image_id"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff" yaml:"poweroff"`
}

type AzureReservationRequest struct {
	PubkeyID int64 `json:"pubkey_id" yaml:"pubkey_id"`

	SourceID string `json:"source_id" yaml:"source_id"`

	// Image Builder UUID of the image that should be launched. This can be directly Azure image ID.
	ImageID string `json:"image_id" yaml:"image_id"`

	// ResourceGroup to use to deploy the resources into
	ResourceGroup string `json:"resource_group" yaml:"resource_group" description:"Azure resource group name to deploy the VM resources into. Optional, defaults to images resource group and when not found to 'redhat-deployed'."`

	// Azure Location also known as region to deploy the VM into.
	// Be aware it needs to be the same as the image location.
	// Defaults to the Resource group location or 'eastus' if new resource group is also created in this request.
	Location string `json:"location" yaml:"location" description:"Location (also known as region) to deploy the VM into, be aware it needs to be the same as the image location. Defaults to the Resource Group location, or 'eastus' when also creating the resource group."`

	// Azure Instance type.
	InstanceSize string `json:"instance_size" yaml:"instance_size"`

	// Amount of instances to provision of size: InstanceSize.
	Amount int64 `json:"amount" yaml:"amount"`

	// Name of the instance(s).
	Name string `json:"name" yaml:"name"`

	// Immediately power off the system after initialization.
	PowerOff bool `json:"poweroff" yaml:"poweroff"`
}

type GCPReservationRequest struct {
	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id" yaml:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id" yaml:"source_id"`

	// Optional launch template id global/instanceTemplates/ID or empty string
	LaunchTemplateID string `json:"launch_template_id,omitempty" yaml:"launch_template_id"`

	// Optional name pattern of the instance(s).
	NamePattern string `json:"name_pattern" yaml:"name_pattern"`

	// GCP zone.
	Zone string `json:"zone" yaml:"zone"`

	// GCP Machine type.
	MachineType string `json:"machine_type" yaml:"machine_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 ` json:"amount" yaml:"amount"`

	// Image Builder UUID of the image that should be launched.
	ImageID string `json:"image_id" yaml:"image_id"`

	// Immediately power off the system after initialization.
	PowerOff bool `json:"poweroff" yaml:"poweroff"`
}

type GenericReservationListResponse struct {
	Data     []*GenericReservationResponse `json:"data" yaml:"data"`
	Metadata page.Metadata                 `json:"metadata" yaml:"metadata"`
}

func (p *GenericReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *AWSReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *AWSReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *AzureReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *AzureReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *GCPReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *NoopReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *GenericReservationListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewReservationResponse(reservation *models.Reservation) render.Renderer {
	return reservationResponseMapper(reservation)
}

func (p *GCPReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func NewAWSReservationResponse(reservation *models.AWSReservation, instances []*models.ReservationInstance) render.Renderer {
	instancesResponse := make([]InstanceResponse, len(instances))
	for iter, inst := range instances {
		instancesResponse[iter] = InstanceResponse{InstanceID: inst.InstanceID, Detail: inst.Detail}
	}

	response := AWSReservationResponse{
		PubkeyID:         reservation.PubkeyID,
		ImageID:          reservation.ImageID,
		SourceID:         reservation.SourceID,
		Region:           reservation.Detail.Region,
		Amount:           reservation.Detail.Amount,
		InstanceType:     reservation.Detail.InstanceType,
		ID:               reservation.ID,
		Name:             reservation.Detail.Name,
		PowerOff:         reservation.Detail.PowerOff,
		Instances:        instancesResponse,
		LaunchTemplateID: reservation.Detail.LaunchTemplateID,
	}
	if reservation.AWSReservationID != nil {
		response.AWSReservationID = *reservation.AWSReservationID
	}
	return &response
}

func NewAzureReservationResponse(reservation *models.AzureReservation, instances []*models.ReservationInstance) render.Renderer {
	instanceIds := make([]InstanceResponse, len(instances))
	for iter, inst := range instances {
		instanceIds[iter] = InstanceResponse{InstanceID: inst.InstanceID, Detail: inst.Detail}
	}

	response := AzureReservationResponse{
		PubkeyID:      reservation.PubkeyID,
		ImageID:       reservation.ImageID,
		SourceID:      reservation.SourceID,
		ResourceGroup: reservation.Detail.ResourceGroup,
		Location:      reservation.Detail.Location,
		Amount:        reservation.Detail.Amount,
		InstanceSize:  reservation.Detail.InstanceSize,
		ID:            reservation.ID,
		Name:          reservation.Detail.Name,
		PowerOff:      reservation.Detail.PowerOff,
		Instances:     instanceIds,
	}
	return &response
}

func NewGCPReservationResponse(reservation *models.GCPReservation, instances []*models.ReservationInstance) render.Renderer {
	instanceIds := make([]InstanceResponse, len(instances))
	for iter, inst := range instances {
		instanceIds[iter] = InstanceResponse{
			InstanceID: inst.InstanceID,
			Detail:     inst.Detail,
		}
	}

	response := GCPReservationResponse{
		NamePattern:      *reservation.Detail.NamePattern,
		PubkeyID:         reservation.PubkeyID,
		ImageID:          reservation.ImageID,
		SourceID:         reservation.SourceID,
		Zone:             reservation.Detail.Zone,
		Amount:           reservation.Detail.Amount,
		MachineType:      reservation.Detail.MachineType,
		GCPOperationName: reservation.GCPOperationName,
		ID:               reservation.ID,
		PowerOff:         reservation.Detail.PowerOff,
		Instances:        instanceIds,
		LaunchTemplateID: reservation.Detail.LaunchTemplateID,
	}
	return &response
}

func NewNoopReservationResponse(reservation *models.NoopReservation) render.Renderer {
	return &NoopReservationResponse{
		ID: reservation.ID,
	}
}

func NewReservationListResponse(reservations []*models.Reservation, meta *page.Metadata) render.Renderer {
	list := make([]*GenericReservationResponse, len(reservations))
	for i, reservation := range reservations {
		list[i] = reservationResponseMapper(reservation)
	}
	return &GenericReservationListResponse{Data: list, Metadata: *meta}
}

func reservationResponseMapper(reservation *models.Reservation) *GenericReservationResponse {
	var finishedAt *time.Time
	if reservation.FinishedAt.Valid {
		finishedAt = &reservation.FinishedAt.Time
	}
	var success *bool
	if reservation.Success.Valid {
		success = &reservation.Success.Bool
	}
	return &GenericReservationResponse{
		ID:         reservation.ID,
		Provider:   int(reservation.Provider),
		CreatedAt:  reservation.CreatedAt,
		FinishedAt: finishedAt,
		Status:     reservation.Status,
		Success:    success,
		Steps:      reservation.Steps,
		Step:       reservation.Step,
		StepTitles: reservation.StepTitles,
		Error:      reservation.Error,
	}
}
