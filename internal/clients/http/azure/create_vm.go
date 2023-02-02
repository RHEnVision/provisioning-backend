// strongly inspired by https://github.com/Azure-Samples/azure-sdk-for-go-samples/tree/main/sdk/resourcemanager/compute/create_vm
// from commit 2ea51d9744d5d4680b4983b2e292a7b1d0085ff4

package azure

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var subscriptionId string

const TraceName = "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"

const (
	vmName        = "redhat-vm"
	vnetName      = "redhat-vnet"
	subnetName    = "redhat-subnet"
	nsgName       = "redhat-nsg"
	nicName       = "redhat-nic"
	diskName      = "redhat-disk"
	publicIPName  = "redhat-public-ip"
	adminUsername = "azureuser"
	vpnIPAddress  = "172.22.0.0/16"
)

func (c *client) CreateVM(ctx context.Context, location string, resourceGroupName string, imageID string, pubkey *models.Pubkey, instanceType clients.InstanceTypeName) (*string, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "CreateVM")
	defer span.End()

	logger := logger(ctx)
	logger.Debug().Msg("Creating Azure VM instance")

	virtualNetwork, err := c.createVirtualNetwork(ctx, location, resourceGroupName, vnetName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create virtual network")
		logger.Error().Err(err).Msgf("cannot create virtual network")
		return nil, err
	}
	logger.Trace().Msgf("Using virtual network id=%s", *virtualNetwork.ID)

	subnet, err := c.createSubnets(ctx, resourceGroupName, vnetName, subnetName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create subnet")
		logger.Error().Err(err).Msgf("cannot create subnet")
		return nil, err
	}
	logger.Trace().Msgf("Using subnet id=%s", *subnet.ID)

	publicIP, err := c.createPublicIP(ctx, location, resourceGroupName, publicIPName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create public IP address")
		logger.Error().Err(err).Msgf("cannot create public IP address")
		return nil, err
	}
	logger.Trace().Msgf("Using public IP address id=%s", *publicIP.ID)

	// network security group
	nsg, err := c.createNetworkSecurityGroup(ctx, location, resourceGroupName, nsgName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create network security group")
		logger.Error().Err(err).Msgf("cannot create network security group")
		return nil, err
	}
	logger.Trace().Msgf("Using network security group id=%s", *nsg.ID)

	networkInterface, err := c.createNetworkInterface(ctx, location, resourceGroupName, subnet, publicIP, nsg, nicName)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create network interface")
		logger.Error().Err(err).Msgf("cannot create network interface")
		return nil, err
	}
	logger.Trace().Msgf("Using network interface id=%s", *networkInterface.ID)

	vmParams := c.prepareVirtualMachineParameters(location, armcompute.VirtualMachineSizeTypes(instanceType), networkInterface, imageID, pubkey.Body, diskName)
	virtualMachine, err := c.createVirtualMachine(ctx, resourceGroupName, vmName, vmParams)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create virtual machine")
		logger.Error().Err(err).Msgf("cannot create virtual machine")
		return nil, err
	}
	logger.Debug().Msgf("Created virtual machine id=%s", *virtualMachine.ID)

	return virtualMachine.ID, nil
}

func (c *client) EnsureResourceGroup(ctx context.Context, name string, location string) (*string, error) {
	resourceGroupClient, err := c.newResourceGroupsClient(ctx)
	if err != nil {
		return nil, err
	}

	logger := logger(ctx)
	getResp, err := resourceGroupClient.Get(ctx, name, nil)
	if err == nil {
		return getResp.ResourceGroup.ID, nil
	}
	if err != nil {
		var azErr *azcore.ResponseError
		if errors.As(err, &azErr) && azErr.StatusCode == http.StatusNotFound {
			logger.Debug().Msgf("resource group %s not found, creating", name)
			// 404 is expected, continue
		} else {
			return nil, fmt.Errorf("failed to fetch resource group: %w", err)
		}
	}

	parameters := armresources.ResourceGroup{
		Location: ptr.To(location),
		// TODO tag the resource group by some RH identifier
		// Tags:     map[string]*string{"sample-rs-tag": ptr.To("sample-tag")}, // resource group update tags
	}

	resp, err := resourceGroupClient.CreateOrUpdate(ctx, name, parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create resource group: %w", err)
	}

	return resp.ResourceGroup.ID, nil
}

func (c *client) createVirtualNetwork(ctx context.Context, location string, resourceGroupName string, name string) (*armnetwork.VirtualNetwork, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "createVirtualNetwork")
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

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create virtual network result: %w", err)
	}

	return &resp.VirtualNetwork, nil
}

func (c *client) createSubnets(ctx context.Context, resourceGroupName string, vnetName string, name string) (*armnetwork.Subnet, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "createSubnets")
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

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create subnet result: %w", err)
	}

	return &resp.Subnet, nil
}

func (c *client) createPublicIP(ctx context.Context, location string, resourceGroupName string, name string) (*armnetwork.PublicIPAddress, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "createPublicIP")
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

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create public IP address result: %w", err)
	}
	return &resp.PublicIPAddress, nil
}

func (c *client) createNetworkSecurityGroup(ctx context.Context, location string, resourceGroupName string, name string) (*armnetwork.SecurityGroup, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "createNetworkSecurityGroup")
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

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create network security group result: %w", err)
	}
	return &resp.SecurityGroup, nil
}

func (c *client) createNetworkInterface(ctx context.Context, location string, resourceGroupName string, subnet *armnetwork.Subnet, publicIP *armnetwork.PublicIPAddress, nsg *armnetwork.SecurityGroup, name string) (*armnetwork.Interface, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "createNetworkInterface")
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

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create network interface result: %w", err)
	}

	return &resp.Interface, nil
}

func (c *client) prepareVirtualMachineParameters(location string, instanceType armcompute.VirtualMachineSizeTypes, networkInterface *armnetwork.Interface, imageID string, sshKeyBody string, diskName string) *armcompute.VirtualMachine {
	return &armcompute.VirtualMachine{
		Location: to.Ptr(location),
		Identity: &armcompute.VirtualMachineIdentity{
			Type: to.Ptr(armcompute.ResourceIdentityTypeNone),
		},
		Properties: &armcompute.VirtualMachineProperties{
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					ID: ptr.To(imageID),
				},
				OSDisk: &armcompute.OSDisk{
					Name:         to.Ptr(diskName),
					CreateOption: to.Ptr(armcompute.DiskCreateOptionTypesFromImage),
					Caching:      to.Ptr(armcompute.CachingTypesReadWrite),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: to.Ptr(armcompute.StorageAccountTypesStandardLRS), // OSDisk type Standard/Premium HDD/SSD
					},
					// DiskSizeGB: to.Ptr[int32](100), // default 127G
				},
			},
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: to.Ptr(instanceType), // VM size include vCPUs,RAM,Data Disks,Temp storage.
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
								KeyData: to.Ptr(sshKeyBody),
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
		},
	}
}

func (c *client) createVirtualMachine(ctx context.Context, resourceGroupName string, vmName string, parameters *armcompute.VirtualMachine) (*armcompute.VirtualMachine, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "createVirtualMachine")
	defer span.End()

	vmClient, err := c.newVirtualMachinesClient(ctx)
	if err != nil {
		return nil, err
	}

	pollerResponse, err := vmClient.BeginCreateOrUpdate(ctx, resourceGroupName, vmName, *parameters, nil)
	if err != nil {
		return nil, fmt.Errorf("create of virtual machine failed to start: %w", err)
	}

	resp, err := pollerResponse.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to poll for create virtual machine status: %w", err)
	}

	return &resp.VirtualMachine, nil
}
