//go:build !test

package image_builder

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetImageBuilderClient = newImageBuilderClient
}
