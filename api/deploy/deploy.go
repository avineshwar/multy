package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/multycloud/multy/api/converter"
	"github.com/multycloud/multy/api/errors"
	"github.com/multycloud/multy/api/proto/commonpb"
	"github.com/multycloud/multy/api/proto/configpb"
	"github.com/multycloud/multy/api/proto/credspb"
	"github.com/multycloud/multy/api/util"
	"github.com/multycloud/multy/encoder"
	"github.com/multycloud/multy/flags"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	"github.com/multycloud/multy/resources/output"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/resources/types"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	tfFile = "main.tf"
	tfDir  = ".terraform"
)

var (
	AwsCredsNotSetErr   = status.Error(codes.InvalidArgument, "aws credentials are required but not set")
	AzureCredsNotSetErr = status.Error(codes.InvalidArgument, "azure credentials are required but not set")
)

type EncodedResources struct {
	HclString         string
	affectedResources []string
}

func Encode(credentials *credspb.CloudCredentials, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) (EncodedResources, error) {
	result := EncodedResources{}
	res, err := GetResources(c)
	if err != nil {
		return result, err
	}

	provider, err := getExistingProvider(prev, credentials)
	if err != nil {
		return result, err
	}

	decodedResources := &encoder.DecodedResources{
		Resources: res,
		Providers: provider,
	}

	updateMultyResourceGroups(decodedResources, c)
	addImplicitResourceGroups(c, res)

	encoded, err := encoder.Encode(decodedResources, credentials)
	if err != nil {
		return result, err
	}
	if len(encoded.ValidationErrs) > 0 {
		return result, errors.ValidationErrors(encoded.ValidationErrs)
	}

	result.HclString = encoded.HclString
	for _, r := range decodedResources.Resources.ResourceMap {
		if r.GetCloud() == commonpb.CloudProvider_AWS && !hasValidAwsCreds(credentials) {
			return result, AwsCredsNotSetErr
		}
		if r.GetCloud() == commonpb.CloudProvider_AZURE && !hasValidAzureCreds(credentials) {
			return result, AzureCredsNotSetErr
		}
	}

	updateDeployedResourceList(encoded, c)

	affectedResources := map[string]struct{}{}
	if prev != nil {
		for _, dr := range prev.DeployedResourceGroup.DeployedResource {
			affectedResources[dr] = struct{}{}
		}
	}

	if curr != nil {
		for _, dr := range curr.DeployedResourceGroup.DeployedResource {
			affectedResources[dr] = struct{}{}
		}
		if rgId := getImplicitResourceGroupId(decodedResources.Resources, curr); rgId != nil {
			affectedResources[*rgId] = struct{}{}
		}
	}

	result.affectedResources = maps.Keys(affectedResources)

	return result, nil
}

func getImplicitResourceGroupId(res resources.Resources, curr *configpb.Resource) *string {
	if commonArgs, ok := res.ResourceMap[curr.ResourceId].GetCommonArgs().(*commonpb.ResourceCommonArgs); ok {
		if commonArgs.ResourceGroupId == "" {
			metadata, err := converter.GetConverter(curr.ResourceArgs.ResourceArgs.MessageName())
			if err == nil {
				res := rg.GetDefaultResourceGroupIdString(metadata.AbbreviatedName, curr.DeployedResourceGroup.GroupId)
				return &res
			}
		}
	}

	return nil
}

func addImplicitResourceGroups(c *configpb.Config, res resources.Resources) {
	for _, r := range c.Resources {
		if res.ResourceMap[r.ResourceId].GetCloud() == commonpb.CloudProvider_AZURE {
			// TODO: only create RGs if r.ImplicitResourceGroup
			if args, ok := res.ResourceMap[r.ResourceId].GetCommonArgs().(*commonpb.ResourceCommonArgs); ok {
				resourceGroup := &rg.Type{
					ResourceId: args.ResourceGroupId,
					Name:       args.ResourceGroupId,
					Location:   args.Location,
					Cloud:      commonpb.CloudProvider_AZURE,
				}
				res.ResourceMap[resourceGroup.GetResourceId()] = resourceGroup
			}
		}
	}
}

func updateDeployedResourceList(encoded encoder.EncodedResources, c *configpb.Config) {
	for _, resource := range c.Resources {
		// compute everything from scratch
		resource.DeployedResourceGroup.DeployedResource = nil
	}

	for r, deployedResources := range encoded.DeployedResources {
		for _, resource := range c.Resources {
			if resource.ResourceId == r.GetResourceId() {
				resource.GetDeployedResourceGroup().DeployedResource = append(resource.GetDeployedResourceGroup().DeployedResource, deployedResources...)
				break
			}
		}
	}
}

func hasValidAzureCreds(credentials *credspb.CloudCredentials) bool {
	return credentials.GetAzureCreds().GetSubscriptionId() != "" &&
		credentials.GetAzureCreds().GetClientId() != "" &&
		credentials.GetAzureCreds().GetTenantId() != "" &&
		credentials.GetAzureCreds().GetClientSecret() != ""
}

func hasValidAwsCreds(credentials *credspb.CloudCredentials) bool {
	return credentials.GetAwsCreds().GetAccessKey() != "" && credentials.GetAwsCreds().GetSecretKey() != ""
}

func updateMultyResourceGroups(decodedResources *encoder.DecodedResources, c *configpb.Config) {
	resourcespbById := map[string]*configpb.Resource{}
	for _, resource := range c.Resources {
		resourcespbById[resource.ResourceId] = resource
	}
	groupIdsByResourceIds := map[string]string{}
	for _, r := range c.Resources {
		if r.DeployedResourceGroup != nil {
			groupIdsByResourceIds[r.ResourceId] = r.DeployedResourceGroup.GroupId
		}
	}
	deployedResourcesByGroupId := map[string]*configpb.DeployedResourceGroup{}
	for groupId, group := range decodedResources.Resources.GetMultyResourceGroups(groupIdsByResourceIds) {
		if _, ok := deployedResourcesByGroupId[groupId]; !ok {
			deployedResourcesByGroupId[groupId] = &configpb.DeployedResourceGroup{GroupId: groupId}
		}
		for _, resource := range group.Resources {
			if _, ok := resourcespbById[resource.GetResourceId()]; ok {
				resourcespbById[resource.GetResourceId()].DeployedResourceGroup = deployedResourcesByGroupId[groupId]
			} else {
				log.Printf("[DEBUG] not adding %s to a group, it doesn't exist in the state\n", resource.GetResourceId())
			}
		}
	}

	for _, r := range c.Resources {
		if commonArgs, ok := decodedResources.Resources.ResourceMap[r.ResourceId].GetCommonArgs().(*commonpb.ResourceCommonArgs); ok {
			if rgId := getImplicitResourceGroupId(decodedResources.Resources, r); rgId != nil {
				r.ImplicitResourceGroup = true
				// this won't be persisted
				commonArgs.ResourceGroupId = *rgId
			}
		}
	}
}

func GetResources(c *configpb.Config) (resources.Resources, error) {
	res := resources.NewResources()

	for _, r := range c.Resources {
		conv, err := converter.GetConverter(r.ResourceArgs.ResourceArgs.MessageName())
		if err != nil {
			return res, err
		}
		err = addMultyResource(r, res, conv)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

type tfOutput struct {
	Level      string `json:"@level"`
	Message    string `json:"@message"`
	Diagnostic struct {
		Summary string `json:"summary"`
		Detail  string `json:"detail"`
		Address string `json:"address"`
	} `json:"diagnostic"`
}

func Deploy(ctx context.Context, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource) (*output.TfState, error) {
	tmpDir := filepath.Join(os.TempDir(), "multy", c.UserId)
	encoded, err := EncodeAndStoreTfFile(ctx, c, prev, curr, false)
	if err != nil {
		return nil, err
	}

	err = MaybeInit(ctx, c.UserId, false)
	if err != nil {
		return nil, err
	}

	var targetArgs []string

	log.Println("[INFO] Running apply for targets:")
	for _, id := range encoded.affectedResources {
		log.Printf("[INFO] %s", id)
		targetArgs = append(targetArgs, "-target="+id)
	}

	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] apply finished in %s", time.Since(start))
	}()

	if len(encoded.affectedResources) != 0 {
		cmd := exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + tmpDir, "apply", "-refresh=false", "-auto-approve", "--json"}, targetArgs...)...)
		if flags.DryRun {
			cmd = exec.CommandContext(ctx, "terraform", append([]string{"-chdir=" + tmpDir, "plan", "-refresh=false", "--json"}, targetArgs...)...)
		}
		outputJson := new(bytes.Buffer)
		cmd.Stdout = outputJson
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			outputs, parseErr := parseTfOutputs(outputJson)
			if parseErr != nil {
				return nil, errors.InternalServerErrorWithMessage("error deploying resources", parseErr)
			}
			if parseErr := getFirstError(outputs); parseErr != nil {
				return nil, errors.DeployError(parseErr)
			}
			return nil, errors.InternalServerErrorWithMessage("error deploying resources", err)
		}
	}

	state, err := GetState(ctx, c.UserId, false)
	if err != nil {
		return state, errors.InternalServerErrorWithMessage("error parsing state", err)
	}

	return state, nil
}

func EncodeAndStoreTfFile(ctx context.Context, c *configpb.Config, prev *configpb.Resource, curr *configpb.Resource, readonly bool) (EncodedResources, error) {
	credentials, err := util.ExtractCloudCredentials(ctx)
	if err != nil {
		return EncodedResources{}, err
	}
	encoded, err := Encode(credentials, c, prev, curr)
	if err != nil {
		return encoded, err
	}

	tfBlock, err := GetTerraformBlock(c.UserId)
	if err != nil {
		return encoded, err
	}

	// TODO: move this to a proper place
	hclOutput := tfBlock + encoded.HclString

	dir := "multy"
	if readonly {
		dir = "multy/readonly"
	}

	tmpDir := filepath.Join(os.TempDir(), dir, c.UserId)
	err = os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
	if err != nil {
		return EncodedResources{}, err
	}
	err = os.WriteFile(filepath.Join(tmpDir, tfFile), []byte(hclOutput), os.ModePerm&0664)
	return encoded, err
}

func getFirstError(outputs []tfOutput) error {
	var err error
	for _, o := range outputs {
		if o.Level == "error" {
			log.Printf("[ERROR] %s\n%s\n", o.Diagnostic.Summary, o.Diagnostic.Detail)
			if err == nil {
				err = fmt.Errorf(o.Diagnostic.Summary)
			}
		}
	}
	return err
}

func parseTfOutputs(outputJson *bytes.Buffer) ([]tfOutput, error) {
	var out []tfOutput
	line, err := outputJson.ReadString('\n')
	for ; err == nil; line, err = outputJson.ReadString('\n') {
		elem := tfOutput{}
		err = json.Unmarshal([]byte(line), &elem)
		if err != nil {
			return nil, err
		}
		out = append(out, elem)
	}

	if err == io.EOF {
		return out, nil
	}

	return nil, err
}

func MaybeInit(ctx context.Context, userId string, readonly bool) error {
	dir := "multy"
	if readonly {
		dir = "multy/readonly"
	}

	tmpDir := filepath.Join(os.TempDir(), dir, userId)
	_, err := os.Stat(filepath.Join(tmpDir, tfDir))
	if os.IsNotExist(err) {
		start := time.Now()
		defer func() {
			log.Printf("[DEBUG] init finished in %s", time.Since(start))
		}()

		err := os.MkdirAll(tmpDir, os.ModeDir|(os.ModePerm&0775))
		if err != nil {
			return errors.InternalServerErrorWithMessage("error creating output file", err)
		}

		cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "init")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func GetState(ctx context.Context, userId string, readonly bool) (*output.TfState, error) {
	dir := "multy"
	if readonly {
		dir = "multy/readonly"
	}

	tmpDir := filepath.Join(os.TempDir(), dir, userId)
	state := output.TfState{}
	stateJson := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "show", "-json")
	cmd.Stdout = stateJson
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return &state, err
	}

	err = json.Unmarshal(stateJson.Bytes(), &state)
	if err != nil {
		return nil, err
	}
	return &state, err
}

func RefreshState(ctx context.Context, userId string, readonly bool) error {
	start := time.Now()
	defer func() {
		log.Printf("[DEBUG] refresh finished in %s", time.Since(start))
	}()
	dir := "multy"
	if readonly {
		dir = "multy/readonly"
	}

	tmpDir := filepath.Join(os.TempDir(), dir, userId)
	outputJson := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "terraform", "-chdir="+tmpDir, "refresh", "-json")
	cmd.Stdout = outputJson
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		outputs, parseErr := parseTfOutputs(outputJson)
		if parseErr != nil {
			return errors.InternalServerErrorWithMessage("error querying resources", parseErr)
		}
		if parseErr := getFirstError(outputs); parseErr != nil {
			return errors.InternalServerErrorWithMessage("error querying resources", parseErr)
		}
		return errors.InternalServerErrorWithMessage("error querying resources", err)
	}
	return err
}

func addMultyResource(r *configpb.Resource, res resources.Resources, metadata *converter.ResourceMetadata) error {
	m, err := r.ResourceArgs.ResourceArgs.UnmarshalNew()
	if err != nil {
		return err
	}

	translatedResource, err := metadata.InitFunc(r.ResourceId, m, res)
	if err != nil {
		return err
	}
	res.ResourceMap[translatedResource.GetResourceId()] = translatedResource
	return nil
}

type WithCommonParams interface {
	GetCommonParameters() *commonpb.ResourceCommonArgs
}

func getExistingProvider(r *configpb.Resource, creds *credspb.CloudCredentials) (map[commonpb.CloudProvider]map[string]*types.Provider, error) {
	providers := map[commonpb.CloudProvider]map[string]*types.Provider{}
	if r != nil {
		args := r.GetResourceArgs().GetResourceArgs()
		m, err := args.UnmarshalNew()
		if err != nil {
			return nil, err
		}
		if wcp, ok := m.(WithCommonParams); ok {
			location, err := common.GetCloudLocation(wcp.GetCommonParameters().Location, wcp.GetCommonParameters().CloudProvider)
			if err != nil {
				return nil, err
			}
			providers[wcp.GetCommonParameters().CloudProvider] = map[string]*types.Provider{
				location: {
					Cloud:        wcp.GetCommonParameters().CloudProvider,
					Location:     location,
					NumResources: 1,
				},
			}

		}
	}

	// Here we use a default location so that if there are lingering resources in the state we don't throw an error.
	// It doesn't work perfectly tho -- AWS resources will be removed by terraform from the state if they don't exist
	// in our config and will no longer be tracked.
	defaultAzureLocation := common.LOCATION[commonpb.Location_EU_WEST_1][commonpb.CloudProvider_AZURE]
	defaultAwsLocation := common.LOCATION[commonpb.Location_EU_WEST_1][commonpb.CloudProvider_AWS]

	if hasValidAwsCreds(creds) && providers[commonpb.CloudProvider_AZURE] == nil {
		providers[commonpb.CloudProvider_AZURE] = map[string]*types.Provider{
			defaultAzureLocation: {
				Cloud:        commonpb.CloudProvider_AZURE,
				Location:     defaultAzureLocation,
				NumResources: 1,
			},
		}
	}
	if hasValidAzureCreds(creds) && providers[commonpb.CloudProvider_AWS] == nil {
		providers[commonpb.CloudProvider_AWS] = map[string]*types.Provider{
			defaultAwsLocation: {
				Cloud:        commonpb.CloudProvider_AWS,
				Location:     defaultAwsLocation,
				NumResources: 1,
			},
		}
	}

	return providers, nil
}
