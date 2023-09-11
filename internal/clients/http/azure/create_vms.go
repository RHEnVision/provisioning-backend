package azure

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func (c *client) CreateVMs(ctx context.Context, vmParams clients.AzureInstanceParams, amount int64, vmNamePrefix string) ([]clients.InstanceDescription, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "CreateVMs")
	defer span.End()

	logger := logger(ctx)
	logger.Debug().Msgf("Started creating %d Azure VM instances", amount)

	subnet, nsg, err := c.ensureSharedNetworking(ctx, vmParams.Location, vmParams.ResourceGroupName)
	if err != nil {
		return nil, err
	}

	vmDescriptions := make([]clients.InstanceDescription, amount)
	resumeTokens := make([]string, amount)
	var i int64
	for i = 0; i < amount; i++ {
		uid, err := uuid.NewUUID()
		if err != nil {
			return vmDescriptions, fmt.Errorf("could not generate a new UUID: %w", err)
		}
		vmName := fmt.Sprintf("%s-%s", vmNamePrefix, uid.String())

		networkInterface, publicIP, err := c.prepareVMNetworking(ctx, subnet, nsg, vmParams, vmName)
		if err != nil {
			return nil, err
		}

		vmDescriptions[i].PublicIPv4 = *publicIP.Properties.IPAddress

		resumeTokens[i], err = c.BeginCreateVM(ctx, networkInterface, vmParams, vmName)
		if err != nil {
			span.SetStatus(codes.Error, "failed to start creation of Azure instance")
			return vmDescriptions, fmt.Errorf("cannot start a create of Azure instance(s): %w", err)
		}
	}

	for j, token := range resumeTokens {
		instanceId, err := c.WaitForVM(ctx, token)
		if err != nil {
			span.SetStatus(codes.Error, "failed to create Azure instance")
			return vmDescriptions, fmt.Errorf("cannot create Azure instance(s): %w", err)
		}
		vmDescriptions[j].ID = string(instanceId)
		logger.Debug().Msgf("Created new instance (%s) via Azure CreateVM", string(instanceId))
	}

	logger.Debug().Msgf("Created %d new instance", amount)

	return vmDescriptions, nil
}
