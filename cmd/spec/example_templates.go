package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

var LaunchTemplateListResponse = payloads.LaunchTemplateListResponse{
	Data: []*payloads.LaunchTemplateResponse{
		{
			ID:   "lt-9843797432897342",
			Name: "XXL large backend API",
		},
	},
	Metadata: page.Metadata{
		Links: page.Links{
			Next: "",
		},
	},
}
