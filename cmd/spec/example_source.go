package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

var SourceListResponse = payloads.SourceListResponse{
	Data: []*payloads.SourceResponse{
		{
			ID:       "654321",
			Name:     "My AWS account",
			Provider: models.ProviderTypeAWS.String(),
			Status:   "available",
		}, {
			ID:       "543621",
			Name:     "My other AWS account",
			Provider: models.ProviderTypeAWS.String(),
			Status:   "available",
		},
	},
	Metadata: page.Metadata{
		Total: 4,
		Links: page.Links{
			Previous: "/api/provisioning/v1/sources?limit=2&offset=0",
			Next:     "",
		},
	},
}

var SourceUploadInfoAWSResponse = payloads.SourceUploadInfoResponse{
	Provider: "aws",
	AwsInfo:  &clients.AccountDetailsAWS{AccountID: "78462784632"},
}

var SourceUploadInfoAzureResponse = payloads.SourceUploadInfoResponse{
	Provider: "azure",
	AzureInfo: &clients.AccountDetailsAzure{
		TenantID:       "617807e1-e4e0-481c-983c-be3ce1e49253",
		SubscriptionID: "617807e1-e4e0-4855-983c-1e3ce1e49674",
		ResourceGroups: []string{"MyGroup 1", "MyGroup 42"},
	},
}
