package resources

import (
	"fmt"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/resources/output"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

type ResourceArgs map[string]proto.Message

type ImportFunc[T proto.Message, O Resource] func(string, T, Resources) (O, error)

// TODO: have resources create RGs as well, and store them
type ResourceCreateFunc[T proto.Message, O Resource] func(string, T, Resources) (O, error)
type ResourceUpdateFunc[T proto.Message, O Resource] func(O, T, Resources) error
type ResourceReadFromStateFunc[O Resource, OutT proto.Message] func(O, *output.TfState) (OutT, error)

type ExportFunc[T proto.Message, R Resource] func(R) (T, error)
type DepsGetterFunc[T proto.Message, R Resource] func(string, T, Resources) ([]string, error)

type ResourceMetadata[ArgsT proto.Message, R Resource, OutT proto.Message] struct {
	DepsGetter DepsGetterFunc[ArgsT, R]

	CreateFunc        ResourceCreateFunc[ArgsT, R]
	UpdateFunc        ResourceUpdateFunc[ArgsT, R]
	ReadFromStateFunc ResourceReadFromStateFunc[R, OutT]

	ImportFunc ImportFunc[ArgsT, R]
	ExportFunc ExportFunc[ArgsT, R]

	AbbreviatedName string
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Create(resourceId string, args proto.Message, resources Resources) (Resource, error) {
	return m.CreateFunc(resourceId, args.(ArgsT), resources)
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Update(resource Resource, args proto.Message, resources Resources) error {
	return m.UpdateFunc(resource.(R), args.(ArgsT), resources)
}

func (m *ResourceMetadata[ArgsT, R, OutT]) ReadFromState(resource Resource, state *output.TfState) (proto.Message, error) {
	return m.ReadFromStateFunc(resource.(R), state)
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Export(resource Resource) (proto.Message, error) {
	return m.ExportFunc(resource.(R))
}

func (m *ResourceMetadata[ArgsT, R, OutT]) Import(resourceId string, args proto.Message, resources Resources) (Resource, error) {
	return m.ImportFunc(resourceId, args.(ArgsT), resources)
}

func (m *ResourceMetadata[ArgsT, R, OutT]) GetAbbreviatedName() string {
	return m.AbbreviatedName
}

type ResourceMetadataInterface interface {
	Create(string, proto.Message, Resources) (Resource, error)
	Update(Resource, proto.Message, Resources) error
	ReadFromState(Resource, *output.TfState) (proto.Message, error)

	Export(Resource) (proto.Message, error)
	Import(string, proto.Message, Resources) (Resource, error)
	GetAbbreviatedName() string
}

func FindArgsOfType[T proto.Message](args ResourceArgs) map[string]T {
	res := map[string]T{}
	for resourceId, arg := range args {
		if casted, ok := arg.(T); ok {
			res[resourceId] = casted
		}
	}
	return res
}

type MultyConfig struct {
	Resources                  Resources
	c                          *configpb.Config
	groupsByResourceId         map[string]*MultyResourceGroup
	affectedResourcesByGroupId map[string][]string
	metadatas                  map[proto.Message]ResourceMetadataInterface
}

func (c *MultyConfig) GetUserId() string {
	return c.c.UserId
}

func LoadConfig(c *configpb.Config, metadatas map[proto.Message]ResourceMetadataInterface) (*MultyConfig, error) {
	multyc := &MultyConfig{
		c:         c,
		metadatas: metadatas,
	}
	res := NewResources()

	for _, r := range c.Resources {
		conv, err := multyc.getConverter(r.ResourceArgs.ResourceArgs.MessageName())
		if err != nil {
			return multyc, err
		}
		err = addMultyResource(r, res, conv)
		if err != nil {
			return multyc, err
		}
	}

	multyc.Resources = res

	affectedResources := map[string][]string{}
	for _, r := range c.Resources {
		affectedResources[r.ResourceId] = r.GetDeployedResourceGroup().GetDeployedResource()
	}

	return multyc, nil
}

func addMultyResource(r *configpb.Resource, res Resources, metadata ResourceMetadataInterface) error {
	m, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
	if err != nil {
		return err
	}

	translatedResource, err := metadata.Import(r.ResourceId, m, res)
	if err != nil {
		return err
	}

	res.Add(translatedResource)
	return nil
}

func (c *MultyConfig) CreateResource(args proto.Message) (Resource, error) {
	conv, err := c.getConverter(proto.MessageName(args))
	if err != nil {
		return nil, err
	}
	c.c.ResourceCounter += 1
	resourceId := fmt.Sprintf("multy_%s_r%d", conv.GetAbbreviatedName(), c.c.ResourceCounter)
	r, err := conv.Create(resourceId, args, c.Resources)
	if err != nil {
		return nil, err
	}
	c.Resources.Add(r)
	return r, nil
}

func (c *MultyConfig) UpdateResource(resourceId string, args proto.Message) (Resource, error) {
	conv, err := c.getConverter(proto.MessageName(args))
	if err != nil {
		return nil, err
	}
	r, exists := c.Resources.ResourceMap[resourceId]
	if !exists {
		return nil, errors.ResourceNotFound(resourceId)
	}
	err = conv.Update(r, args, c.Resources)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *MultyConfig) DeleteResource(resourceId string) (Resource, error) {
	if r, exists := c.Resources.ResourceMap[resourceId]; exists {
		c.Resources.Delete(resourceId)
		return r, nil
	} else {
		return nil, errors.ResourceNotFound(resourceId)
	}

}

func (c *MultyConfig) UpdateMultyResourceGroups() {
	groupIdsByResourceIds := map[string]string{}
	for _, r := range c.c.Resources {
		if r.DeployedResourceGroup != nil {
			groupIdsByResourceIds[r.ResourceId] = r.DeployedResourceGroup.GroupId
		}
	}
	c.groupsByResourceId = c.Resources.GetMultyResourceGroups(groupIdsByResourceIds)
}

func (c *MultyConfig) UpdateDeployedResourceList(deployedResources map[Resource][]string) {
	affectedResources := map[string][]string{}

	for r, deployedResource := range deployedResources {
		groupId := c.groupsByResourceId[r.GetResourceId()]
		for _, d := range deployedResource {
			if !slices.Contains(affectedResources[groupId.GroupId], d) {
				affectedResources[groupId.GroupId] = append(affectedResources[r.GetResourceId()], d)
			}
		}
	}

	c.affectedResourcesByGroupId = affectedResources
}

func (c *MultyConfig) GetAffectedResources(resourceId string) []string {
	return c.affectedResourcesByGroupId[c.groupsByResourceId[resourceId].GroupId]
}

func (c *MultyConfig) ExportConfig() (*configpb.Config, error) {
	result := &configpb.Config{
		UserId:          c.GetUserId(),
		ResourceCounter: c.c.ResourceCounter,
	}

	for _, r := range c.Resources.GetAll() {
		msg, err := r.GetMetadata().Export(r)
		if err != nil {
			return nil, err
		}
		a, err := anypb.New(msg)
		if err != nil {
			return nil, err
		}
		resource := configpb.Resource{
			ResourceId:   r.GetResourceId(),
			ResourceArgs: &configpb.ResourceArgs{ResourceArgs: a},
		}
		result.Resources = append(result.Resources, &resource)
	}

	return result, nil
}

func (c *MultyConfig) getConverter(name protoreflect.FullName) (ResourceMetadataInterface, error) {
	for messageType, conv := range c.metadatas {
		if name == proto.MessageName(messageType) {
			return conv, nil
		}
	}
	return nil, fmt.Errorf("unknown resource type %s", name)
}
