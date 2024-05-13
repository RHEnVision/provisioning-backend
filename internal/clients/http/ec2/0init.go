//go:build !test

package ec2

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetEC2Client = newAssumedEC2ClientWithRegion
	clients.GetServiceEC2Client = newEC2ClientWithRegion
}
