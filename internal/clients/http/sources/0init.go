//go:build !test

package sources

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetSourcesClient = newSourcesClient
}
