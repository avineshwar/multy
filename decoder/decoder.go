package decoder

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/multycloud/multy/hclutil"
	"github.com/multycloud/multy/mhcl"
	"github.com/multycloud/multy/parser"
	"github.com/multycloud/multy/resources"
	"github.com/multycloud/multy/resources/common"
	rg "github.com/multycloud/multy/resources/resource_group"
	"github.com/multycloud/multy/resources/types"
	"github.com/multycloud/multy/validate"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/zclconf/go-cty/cty"
)

type DecodedGlobalConfig struct {
	Location      string
	Clouds        []common.CloudProvider
	DefaultRgName hcl.Expression
}

type DecodedResources struct {
	Resources    map[string]resources.CloudSpecificResource
	Outputs      map[string]DecodedOutput
	Providers    map[string]types.Provider
	GlobalConfig DecodedGlobalConfig
}

type DecodedOutput struct {
	OutputId        string
	Value           cty.Value
	DefinitionRange hcl.Range
}

func Decode(config parser.ParsedConfig) *DecodedResources {
	vars := map[string]cty.Value{}
	for _, v := range config.Variables {
		vars[v.Name] = v.Value("default")
	}
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{"var": cty.ObjectVal(vars)},
		Functions: nil,
	}

	globalConfig := decodeGlobalConfig(config, ctx)

	multyResources := topologicalSort(config.MultyResources)
	r := map[string]resources.CloudSpecificResource{}

	defaultRgId := getDefaultResourceGroupId(
		hclutil.GetOptionalAttributeAsExpr(
			config.GlobalConfig, "default_resource_group_name",
		),
	)
	globalConfig.DefaultRgName = defaultRgId

	resourceDecoder := ResourceDecoder{globalConfig: globalConfig}
	cloudSpecificCtx := InitCloudSpecificContext(ctx, config.Variables)
	mhclProcessor := mhcl.MHCLProcessor{ResourceRefs: r}

	for _, resource := range multyResources {
		decodedResources, outCtx := resourceDecoder.Decode(*resource, cloudSpecificCtx, mhclProcessor)
		cloudSpecificCtx.AddCtx(outCtx)
		for _, decodedResource := range decodedResources {
			uniqueId := decodedResource.GetResourceId()
			if duplicate, ok := r[uniqueId]; ok {
				if !duplicate.ImplicitlyCreated || !decodedResource.ImplicitlyCreated {
					// TODO: allow specifying multiple ranges
					validate.LogFatalWithSourceRange(
						resource.DefinitionRange, "duplicate resources found with id %s ", uniqueId,
					)
				}
			}
			r[uniqueId] = decodedResource
		}
	}

	outputs := map[string]DecodedOutput{}
	for _, output := range config.Outputs {
		val, diags := output.Value.Value(cloudSpecificCtx.GetCloudAgnosticContext())
		if diags != nil {
			validate.LogFatalWithDiags(diags, "unable to parse output %s", output.ID)
		}
		outputs[output.ID] = DecodedOutput{
			OutputId:        output.ID,
			Value:           val,
			DefinitionRange: output.DefinitionRange,
		}
	}

	return &DecodedResources{
		Resources:    r,
		GlobalConfig: globalConfig,
		Outputs:      outputs,
	}
}

func topologicalSort(allResources []*parser.MultyResource) []*parser.MultyResource {
	positions := map[*parser.MultyResource]int{}
	position := 0
	for _, resource := range allResources {
		position = topologicalSortHelper(positions, position, resource, nil)
	}

	result := make([]*parser.MultyResource, len(allResources))
	for resource, i := range positions {
		result[i] = resource
	}
	return result
}

func topologicalSortHelper(positions map[*parser.MultyResource]int, pos int, r *parser.MultyResource, currentStack []parser.MultyResourceDependency) int {
	if _, ok := positions[r]; ok {
		return pos
	}
	for i, otherResource := range currentStack {
		if otherResource.From == r {
			var errorMessages []string
			for j := i; j < len(currentStack); j++ {
				errorMessages = append(
					errorMessages, fmt.Sprintf(
						"[%s]%s=>%s", currentStack[j].SourceRange, currentStack[j].From.ID, currentStack[j].To.ID,
					),
				)
			}
			validate.LogFatalWithSourceRange(
				currentStack[len(currentStack)-1].SourceRange,
				"found cycle while resolving dependencies for resource %s.\n%s", r.ID,
				strings.Join(errorMessages, "\n"),
			)
		}
	}

	for _, dep := range r.Dependencies {
		currentStack = append(currentStack, dep)
		pos = topologicalSortHelper(positions, pos, dep.To, currentStack)
		currentStack = currentStack[:len(currentStack)-1]
	}
	positions[r] = pos
	return pos + 1

}

func decodeGlobalConfig(config parser.ParsedConfig, ctx *hcl.EvalContext) DecodedGlobalConfig {
	if config.GlobalConfig == nil {
		// TODO: here we should probably return nil and handle any attempted accesses gracefully
		return DecodedGlobalConfig{}
	}
	type globalConfigHcl struct {
		Location string   `hcl:"location"`
		Clouds   []string `hcl:"clouds"`
		HclBody  hcl.Body `hcl:",remain"`
	}

	var c globalConfigHcl

	diags := gohcl.DecodeBody(config.GlobalConfig, ctx, &c)
	if diags != nil {
		validate.LogFatalWithDiags(diags, "Unable to decode global configuration.")
	}

	globalConfig := DecodedGlobalConfig{
		Location: c.Location,
		Clouds:   castToCloudProvider(c.Clouds),
	}
	return globalConfig
}

func castToCloudProvider(c []string) []common.CloudProvider {
	var res []common.CloudProvider
	for _, cloud := range c {
		res = append(res, common.CloudProvider(cloud))
	}
	return res
}

func getDefaultResourceGroupId(specifiedDefault *hcl.Expression) hcl.Expression {
	if specifiedDefault != nil {
		return *specifiedDefault
	}
	return rg.GetDefaultResourceGroupId()
}

func resolveAttributes(attrs hcl.Attributes, ctx *hcl.EvalContext) (map[string]cty.Value, hcl.Diagnostics) {
	res := map[string]cty.Value{}
	for key, val := range attrs {
		resolved, diags := val.Expr.Value(ctx)
		if diags != nil {
			return res, diags
		}
		res[key] = resolved
	}
	return res, nil
}
