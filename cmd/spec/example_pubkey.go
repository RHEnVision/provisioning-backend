package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

var PubkeyRequest = payloads.PubkeyRequest{
	Name: "My key",
	Body: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
}

var PubkeyResponse = payloads.PubkeyResponse{
	ID:                1,
	AccountID:         1,
	Name:              "My key",
	Body:              "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
	Type:              "ssh-ed25519",
	Fingerprint:       "gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk=",
	FingerprintLegacy: "ee:f1:d4:62:99:ab:17:d9:3b:00:66:62:32:b2:55:9e",
}

var PubkeyListResponse = payloads.PubkeyListResponse{
	Data: []*payloads.PubkeyResponse{
		{
			ID:                3,
			AccountID:         1,
			Name:              "My key",
			Body:              "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
			Type:              "ssh-ed25519",
			Fingerprint:       "gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk=",
			FingerprintLegacy: "ee:f1:d4:62:99:ab:17:d9:3b:00:66:62:32:b2:55:9e",
		},
	},
	Metadata: page.Metadata{
		Total: 3,
	},
	Links: page.Links{
		Previous: "/api/provisioning/v1/pubkeys?limit=2&offset=0",
		Next:     "",
	},
}
