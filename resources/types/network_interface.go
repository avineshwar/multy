package types

import (
	"fmt"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/resourcespb"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	"github.com/multycloud/multy/resources/output/network_interface"
	"github.com/multycloud/multy/validate"
)

var networkInterfaceMetadata = resources.ResourceMetadata[*resourcespb.NetworkInterfaceArgs, *NetworkInterface, *resourcespb.NetworkInterfaceResource]{
	DepsGetter: func(_ string, args *resourcespb.NetworkInterfaceArgs, _ resources.Resources) (res []string, err error) {
		res = append(res, args.SubnetId)
		return res, err
	},
	CreateFunc:        CreateNetworkInterface,
	UpdateFunc:        UpdateNetworkInterface,
	ReadFromStateFunc: NetworkInterfaceFromState,
	ExportFunc:        func(r *NetworkInterface) (*resourcespb.NetworkInterfaceArgs, error) { return r.Args, nil },
	ImportFunc:        NewNetworkInterface,
	AbbreviatedName:   "nic",
}

type NetworkInterface struct {
	resources.ResourceWithId[*resourcespb.NetworkInterfaceArgs]

	Subnet *Subnet
}

func (r *NetworkInterface) GetMetadata() resources.ResourceMetadataInterface {
	return &networkInterfaceMetadata
}

func CreateNetworkInterface(resourceId string, args *resourcespb.NetworkInterfaceArgs, others resources.Resources) (*NetworkInterface, error) {
	if args.CommonParameters.ResourceGroupId == "" {
		subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
		if err != nil {
			return nil, err
		}
		rgId := subnet.VirtualNetwork.Args.CommonParameters.ResourceGroupId
		args.CommonParameters.ResourceGroupId = rgId
	}

	return NewNetworkInterface(resourceId, args, others)
}

func UpdateNetworkInterface(resource *NetworkInterface, vn *resourcespb.NetworkInterfaceArgs, others resources.Resources) error {
	_, err := NewNetworkInterface(resource.ResourceId, vn, others)
	return err
}

func NetworkInterfaceFromState(resource *NetworkInterface, _ *output.TfState) (*resourcespb.NetworkInterfaceResource, error) {
	return &resourcespb.NetworkInterfaceResource{
		CommonParameters: &commonpb.CommonResourceParameters{
			ResourceId:      resource.ResourceId,
			ResourceGroupId: resource.Args.CommonParameters.ResourceGroupId,
			Location:        resource.Args.CommonParameters.Location,
			CloudProvider:   resource.Args.CommonParameters.CloudProvider,
			NeedsUpdate:     false,
		},
		Name:     resource.Args.Name,
		SubnetId: resource.Args.SubnetId,
	}, nil
}

func NewNetworkInterface(resourceId string, args *resourcespb.NetworkInterfaceArgs, others resources.Resources) (*NetworkInterface, error) {
	subnet, err := resources.Get[*Subnet](resourceId, others, args.SubnetId)
	if err != nil {
		return nil, err
	}
	return &NetworkInterface{
		ResourceWithId: resources.ResourceWithId[*resourcespb.NetworkInterfaceArgs]{
			ResourceId: resourceId,
			Args:       args,
		},
		Subnet: subnet,
	}, nil
}

func (r *NetworkInterface) Translate(ctx resources.MultyContext) ([]output.TfBlock, error) {
	subnetId, err := resources.GetMainOutputId(r.Subnet)
	if err != nil {
		return nil, err
	}

	if r.GetCloud() == commonpb.CloudProvider_AWS {
		return []output.TfBlock{
			network_interface.AwsNetworkInterface{
				AwsResource: common.NewAwsResource(r.ResourceId, r.Args.Name),
				SubnetId:    subnetId,
			},
		}, nil
	} else if r.GetCloud() == commonpb.CloudProvider_AZURE {
		rgName := GetResourceGroupName(r.Args.CommonParameters.ResourceGroupId)
		nic := network_interface.AzureNetworkInterface{
			AzResource: common.NewAzResource(
				r.ResourceId, r.Args.Name, rgName,
				r.GetCloudSpecificLocation(),
			),
			// by default, virtual_machine will have a private ip
			IpConfigurations: []network_interface.AzureIpConfiguration{{
				Name:                       "internal", // this name shouldn't be vm.name
				PrivateIpAddressAllocation: "Dynamic",
				SubnetId:                   subnetId,
				Primary:                    true,
			}},
		}
		// associate a public ip configuration in case a public_ip resource references this network_interface
		if publicIpReference := r.getPublicIpReferences(ctx, subnetId); len(publicIpReference) != 0 {
			nic.IpConfigurations = publicIpReference
		}
		return []output.TfBlock{nic}, nil
	}

	return nil, fmt.Errorf("cloud %s is not supported for this resource type ", r.GetCloud().String())
}

func (r *NetworkInterface) GetId(cloud commonpb.CloudProvider) string {
	types := map[commonpb.CloudProvider]string{common.AWS: network_interface.AwsResourceName, common.AZURE: network_interface.AzureResourceName}
	return fmt.Sprintf("%s.%s.id", types[cloud], r.ResourceId)
}

func (r *NetworkInterface) getPublicIpReferences(ctx resources.MultyContext, subnetId string) []network_interface.AzureIpConfiguration {
	var ipConfigurations []network_interface.AzureIpConfiguration
	for _, resource := range resources.GetAllResourcesWithRef(ctx, func(i *PublicIp) *NetworkInterface { return i.NetworkInterface }, r) {
		ipConfigurations = append(
			ipConfigurations, network_interface.AzureIpConfiguration{
				Name:                       fmt.Sprintf("external-%s", resource.Args.Name),
				PrivateIpAddressAllocation: "Dynamic",
				PublicIpAddressId:          resource.GetId(common.AZURE),
				SubnetId:                   subnetId,
				Primary:                    true,
			},
		)
	}
	return ipConfigurations
}

func (r *NetworkInterface) Validate(ctx resources.MultyContext) (errs []validate.ValidationError) {
	errs = append(errs, r.ResourceWithId.Validate()...)
	return errs
}

func (r *NetworkInterface) GetMainResourceName() (string, error) {
	switch r.GetCloud() {
	case commonpb.CloudProvider_AWS:
		return network_interface.AwsResourceName, nil
	case commonpb.CloudProvider_AZURE:
		return network_interface.AzureResourceName, nil
	default:
		return "", fmt.Errorf("unknown cloud %s", r.GetCloud().String())
	}
}
