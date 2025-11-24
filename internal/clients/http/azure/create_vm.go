// strongly inspired by https://github.com/Azure-Samples/azure-sdk-for-go-samples/tree/main/sdk/resourcemanager/compute/create_vm
// from commit 2ea51d9744d5d4680b4983b2e292a7b1d0085ff4

package azure

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"go.opentelemetry.io/otel/codes"
)

const (
	vnetName              = "redhat-vnet"
	subnetName            = "redhat-subnet"
	nsgName               = "redhat-nsg"
	adminUsername         = "azureuser"
	vpnIPAddress          = "172.22.0.0/16"
	resourcePollFrequency = 5 * time.Second
	vmPollFrequency       = 10 * time.Second
)

func newAzureResourceGroup(group armresources.ResourceGroup) clients.AzureResourceGroup {
	return clients.AzureResourceGroup{ID: *group.ID, Name: *group.Name, Location: *group.Location}
}

func (c *client) BeginCreateVM(ctx context.Context, networkInterface *armnetwork.Interface, vmParams clients.AzureInstanceParams, vmName string) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "BeginCreateVM")
	defer span.End()

	logger := logger(ctx)
	logger.Debug().Msg("Creating Azure VM instance without waiting")

	vmClient, err := c.newVirtualMachinesClient(ctx)
	if err != nil {
		return "", err
	}

	vmAzureParams := c.prepareVirtualMachineParameters(vmParams, networkInterface, vmName)

	poller, err := vmClient.BeginCreateOrUpdate(ctx, vmParams.ResourceGroupName, vmName, *vmAzureParams, nil)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create virtual machine")
		logger.Error().Err(err).Msg("cannot create virtual machine")
		return "", fmt.Errorf("create of virtual machine failed to start: %w", err)
	}

	resumeToken, err := poller.ResumeToken()
	if err != nil {
		span.SetStatus(codes.Error, "cannot generate resume token")
		logger.Error().Err(err).Msg("cannot generate the resume token")
		return "", fmt.Errorf("cannot generate Azure resume token: %w", err)
	}

	return resumeToken, nil
}

func (c *client) WaitForVM(ctx context.Context, resumeToken string) (clients.AzureInstanceID, error) {
	ctx, span := telemetry.StartSpan(ctx, "WaitForVM")
	defer span.End()

	logger := logger(ctx)
	logger.Debug().Msgf("Starting polling for Azure VM instance creation, using token: %s", resumeToken)

	vmClient, err := c.newVirtualMachinesClient(ctx)
	if err != nil {
		return "", err
	}

	pollerResponse, err := vmClient.BeginCreateOrUpdate(ctx, "", "", armcompute.VirtualMachine{}, &armcompute.VirtualMachinesClientBeginCreateOrUpdateOptions{
		ResumeToken: resumeToken,
	})
	if err != nil {
		span.SetStatus(codes.Error, "polling of virtual machine creation status failed to start")
		return "", fmt.Errorf("polling of virtual machine creation status failed to start: %w", err)
	}
	resp, err := pollerResponse.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
		Frequency: vmPollFrequency,
	})
	if err != nil {
		span.SetStatus(codes.Error, "failed to poll for create virtual machine status")
		return "", fmt.Errorf("failed to poll for create virtual machine status: %w", err)
	}

	logger.Debug().Msgf("Done creating virtual machine id=%s", *resp.VirtualMachine.ID)

	return clients.AzureInstanceID(*resp.VirtualMachine.ID), nil
}

func (c *client) ensureSharedNetworking(ctx context.Context, location, resourceGroupName string) (*armnetwork.Subnet, *armnetwork.SecurityGroup, error) {
	ctx, span := telemetry.StartSpan(ctx, "ensureSharedNetworking")
	defer span.End()

	logger := logger(ctx)
	virtualNetwork, err := c.createVirtualNetwork(ctx, location, resourceGroupName, vnetName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create virtual network")
		logger.Error().Err(err).Msg("cannot create virtual network")
		return nil, nil, err
	}
	logger.Trace().Msgf("Using virtual network id=%s", *virtualNetwork.ID)

	subnet, err := c.createSubnets(ctx, resourceGroupName, vnetName, subnetName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create subnet")
		logger.Error().Err(err).Msg("cannot create subnet")
		return nil, nil, err
	}
	logger.Trace().Msgf("Using subnet id=%s", *subnet.ID)

	// network security group
	nsg, err := c.createNetworkSecurityGroup(ctx, location, resourceGroupName, nsgName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create network security group")
		logger.Error().Err(err).Msg("cannot create network security group")
		return nil, nil, err
	}
	logger.Trace().Msgf("Using network security group id=%s", *nsg.ID)

	return subnet, nsg, nil
}

func (c *client) prepareVMNetworking(ctx context.Context, subnet *armnetwork.Subnet, securityGroup *armnetwork.SecurityGroup, vmParams clients.AzureInstanceParams, vmName string) (*armnetwork.Interface, *armnetwork.PublicIPAddress, error) {
	ctx, span := telemetry.StartSpan(ctx, "prepareVMNetworking")
	defer span.End()

	logger := logger(ctx)

	publicIPName := vmName + "_ip"
	publicIP, err := c.createPublicIP(ctx, vmParams.Location, vmParams.ResourceGroupName, publicIPName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create public IP address")
		logger.Error().Err(err).Msg("cannot create public IP address")
		return nil, nil, err
	}
	logger.Trace().Msgf("Using public IP address id=%s", *publicIP.ID)
	nicName := vmName + "_nic"
	networkInterface, err := c.createNetworkInterface(ctx, vmParams.Location, vmParams.ResourceGroupName, subnet, publicIP, securityGroup, nicName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create network interface")
		logger.Error().Err(err).Msg("cannot create network interface")
		return nil, publicIP, err
	}
	logger.Trace().Msgf("Using network interface id=%s", *networkInterface.ID)
	return networkInterface, publicIP, nil
}

func (c *client) EnsureResourceGroup(ctx context.Context, name string, location string) (clients.AzureResourceGroup, error) {
	resourceGroupClient, err := c.newResourceGroupsClient(ctx)
	if err != nil {
		return clients.AzureResourceGroup{}, err
	}

	logger := logger(ctx)
	getResp, err := resourceGroupClient.Get(ctx, name, nil)
	if err == nil {
		return newAzureResourceGroup(getResp.ResourceGroup), nil
	}
	if err != nil {
		var azErr *azcore.ResponseError
		if errors.As(err, &azErr) && azErr.StatusCode == http.StatusNotFound {
			logger.Debug().Msgf("resource group %s not found, creating", name)
			// 404 is expected, continue
		} else {
			return clients.AzureResourceGroup{}, fmt.Errorf("failed to fetch resource group: %w", err)
		}
	}

	parameters := armresources.ResourceGroup{
		Location: ptr.To(location),
		// TODO tag the resource group by some RH identifier
		// Tags:     map[string]*string{"sample-rs-tag": ptr.To("sample-tag")}, // resource group update tags
	}

	resp, err := resourceGroupClient.CreateOrUpdate(ctx, name, parameters, nil)
	if err != nil {
		return clients.AzureResourceGroup{}, fmt.Errorf("cannot create resource group: %w", err)
	}

	return newAzureResourceGroup(resp.ResourceGroup), nil
}

func (c *client) createVirtualNetwork(ctx context.Context, location string, resourceGroupName string, name string) (*armnetwork.VirtualNetwork, error) {
	ctx, span := telemetry.StartSpan(ctx, "createVirtualNetwork")
	defer span.End()

	vnetClient, err := c.newVirtualNetworksClient(ctx)
	if err != nil {
		return nil, err
	}

	logger := logger(ctx)
	getResp, err := vnetClient.Get(ctx, resourceGroupName, name, nil)
	if err == nil {
		return &getResp.VirtualNetwork, nil
	}
	if err != nil {
		var azErr *azcore.ResponseError
		if errors.As(err, &azErr) && azErr.StatusCode == http.StatusNotFound {
			logger.Debug().Msgf("virtual network %s not found, creating", name)
			// 404 is expected, continue
		} else {
			return nil, fmt.Errorf("failed to fetch virtual network: %w", err)
		}
	}

	parameters := armnetwork.VirtualNetwork{
		Location: to.Ptr(location),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr(vpnIPAddress),
				},
			},
			//Subnets: []*armnetwork.Subnet{
			//	{
			//		Name: to.Ptr(subnetName+"3"),
			//		Properties: &armnetwork.SubnetPropertiesFormat{
			//			AddressPrefix: to.Ptr("10.1.0.0/24"),
			//		},
			//	},
			//},
		},
	}

	pollerResponse, err := vnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("create of virtual network failed to start: %w", err)
	}

	resp, err := pollerResponse.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
		Frequency: resourcePollFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create virtual network result: %w", err)
	}

	return &resp.VirtualNetwork, nil
}

func (c *client) createSubnets(ctx context.Context, resourceGroupName string, vnetName string, name string) (*armnetwork.Subnet, error) {
	ctx, span := telemetry.StartSpan(ctx, "createSubnets")
	defer span.End()

	subnetClient, err := c.newSubnetsClient(ctx)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			// the subnet takes the full vpn address scope
			AddressPrefix: to.Ptr(vpnIPAddress),
		},
	}

	pollerResponse, err := subnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, vnetName, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("create of subnet failed to start: %w", err)
	}

	resp, err := pollerResponse.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
		Frequency: resourcePollFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create subnet result: %w", err)
	}

	return &resp.Subnet, nil
}

func (c *client) createNetworkSecurityGroup(ctx context.Context, location string, resourceGroupName string, name string) (*armnetwork.SecurityGroup, error) {
	ctx, span := telemetry.StartSpan(ctx, "createNetworkSecurityGroup")
	defer span.End()

	nsgClient, err := c.newSecurityGroupsClient(ctx)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.SecurityGroup{
		Location: to.Ptr(location),
		Properties: &armnetwork.SecurityGroupPropertiesFormat{
			SecurityRules: []*armnetwork.SecurityRule{
				// inbound
				{
					Name: to.Ptr("inbound_22"), //
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
						SourcePortRange:          to.Ptr("*"),
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
						DestinationPortRange:     to.Ptr("22"),
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
						Priority:                 to.Ptr[int32](100),
						Description:              to.Ptr("network security group inbound port 22"),
						Direction:                to.Ptr(armnetwork.SecurityRuleDirectionInbound),
					},
				},
				// outbound
				{
					Name: to.Ptr("outbound_22"), //
					Properties: &armnetwork.SecurityRulePropertiesFormat{
						SourceAddressPrefix:      to.Ptr("0.0.0.0/0"),
						SourcePortRange:          to.Ptr("*"),
						DestinationAddressPrefix: to.Ptr("0.0.0.0/0"),
						DestinationPortRange:     to.Ptr("22"),
						Protocol:                 to.Ptr(armnetwork.SecurityRuleProtocolTCP),
						Access:                   to.Ptr(armnetwork.SecurityRuleAccessAllow),
						Priority:                 to.Ptr[int32](100),
						Description:              to.Ptr("network security group outbound port 22"),
						Direction:                to.Ptr(armnetwork.SecurityRuleDirectionOutbound),
					},
				},
			},
		},
	}

	pollerResponse, err := nsgClient.BeginCreateOrUpdate(ctx, resourceGroupName, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("create of network security group failed to start: %w", err)
	}

	resp, err := pollerResponse.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
		Frequency: resourcePollFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create network security group result: %w", err)
	}
	return &resp.SecurityGroup, nil
}

func (c *client) createPublicIP(ctx context.Context, location string, resourceGroupName string, name string) (*armnetwork.PublicIPAddress, error) {
	ctx, span := telemetry.StartSpan(ctx, "createPublicIP")
	defer span.End()

	publicIPAddressClient, err := c.newPublicIPAddressesClient(ctx)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.PublicIPAddress{
		Location: to.Ptr(location),
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic), // Static or Dynamic
		},
	}

	pollerResponse, err := publicIPAddressClient.BeginCreateOrUpdate(ctx, resourceGroupName, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("create of public IP address failed to start: %w", err)
	}

	resp, err := pollerResponse.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
		Frequency: resourcePollFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create public IP address result: %w", err)
	}
	return &resp.PublicIPAddress, nil
}

func (c *client) createNetworkInterface(ctx context.Context, location string, resourceGroupName string, subnet *armnetwork.Subnet, publicIP *armnetwork.PublicIPAddress, nsg *armnetwork.SecurityGroup, name string) (*armnetwork.Interface, error) {
	ctx, span := telemetry.StartSpan(ctx, "createNetworkInterface")
	defer span.End()

	nicClient, err := c.newInterfacesClient(ctx)
	if err != nil {
		return nil, err
	}

	parameters := armnetwork.Interface{
		Location: to.Ptr(location),
		Properties: &armnetwork.InterfacePropertiesFormat{
			IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
				{
					Name: to.Ptr("ipConfig"),
					Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
						Subnet: &armnetwork.Subnet{
							ID: subnet.ID,
						},
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID: publicIP.ID,
						},
					},
				},
			},
			NetworkSecurityGroup: &armnetwork.SecurityGroup{
				ID: nsg.ID,
			},
		},
	}

	pollerResponse, err := nicClient.BeginCreateOrUpdate(ctx, resourceGroupName, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("create of network interface failed to start: %w", err)
	}

	resp, err := pollerResponse.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
		Frequency: resourcePollFrequency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create network interface result: %w", err)
	}

	return &resp.Interface, nil
}

func (c *client) prepareVirtualMachineParameters(vmParams clients.AzureInstanceParams, networkInterface *armnetwork.Interface, vmName string) *armcompute.VirtualMachine {
	userDataEncoded := make([]byte, base64.StdEncoding.EncodedLen(len(vmParams.UserData)))
	base64.StdEncoding.Encode(userDataEncoded, vmParams.UserData)

	return &armcompute.VirtualMachine{
		Location: to.Ptr(vmParams.Location),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Tags: vmParams.Tags,
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					ID: ptr.To(vmParams.ImageID),
				},
				OSDisk: &armcompute.OSDisk{
					// Name:         ptr.To(vmName + "_disk1"),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
					},
					// DiskSizeGB: to.Ptr[int32](100), // default 127G
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(armcompute.VirtualMachineSizeTypes(vmParams.InstanceType)), // VM size include vCPUs,RAM,Data Disks,Temp storage.
			},
			OSProfile: &armcompute.OSProfile{ //
				ComputerName:  to.Ptr(vmName),
				AdminUsername: to.Ptr(adminUsername),
				// require ssh key for authentication
				LinuxConfiguration: &armcompute.LinuxConfiguration{
					DisablePasswordAuthentication: to.Ptr(true),
					SSH: &armcompute.SSHConfiguration{
						PublicKeys: []*armcompute.SSHPublicKey{
							{
								Path:    to.Ptr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", adminUsername)),
								KeyData: to.Ptr(vmParams.Pubkey.Body),
							},
						},
					},
				},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						ID: networkInterface.ID,
					},
				},
			},
			UserData: to.Ptr(string(userDataEncoded)),
		},
	}
}
