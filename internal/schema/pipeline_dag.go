// Copyright (c) ArenaML Labs Pvt Ltd.

package schema

import (
	"context"
	"fmt"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
	"github.com/arena-ml/terraform-provider-arenaml/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type PipelineDag struct {
	PipelineID  types.String `tfsdk:"pipeline_id"`
	InputEdges  types.Set    `tfsdk:"input_edges"`
	OutputEdges types.Set    `tfsdk:"output_edges"`
}

type InputEdges struct {
	NodeID     types.String `tfsdk:"node_id"`
	FromBases  types.Set    `tfsdk:"from_bases"`
	FromInputs types.Set    `tfsdk:"from_inputs"`
	ToSteps    types.Set    `tfsdk:"to_steps"`
	ToInputs   types.Set    `tfsdk:"to_inputs"`
}

type OutputEdge struct {
	NodeID   types.String `tfsdk:"node_id"`
	FromStep types.String `tfsdk:"from_step"`
	ToInputs types.Set    `tfsdk:"to_inputs"`
}

func (ie *InputEdges) FillFromResp(ctx context.Context, resp client.ModelInputEdges) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, ie)

	if resp.FromBases != nil {
		ie.FromBases, err = helper.FromGoStrSliceToTfSet(ctx, *resp.FromBases)
		if err != nil {
			return
		}
	} else {
		ie.FromBases = basetypes.NewSetNull(types.StringType)
	}

	if resp.FromInputs != nil {
		ie.FromInputs, err = helper.FromGoStrSliceToTfSet(ctx, *resp.FromInputs)
		if err != nil {
			return
		}
	} else {
		ie.FromInputs = basetypes.NewSetNull(types.StringType)
	}

	if resp.ToSteps != nil {
		ie.ToSteps, err = helper.FromGoStrSliceToTfSet(ctx, *resp.ToSteps)
		if err != nil {
			return
		}
	} else {
		ie.ToSteps = basetypes.NewSetNull(types.StringType)
	}

	if resp.ToInputs != nil {
		ie.ToInputs, err = helper.FromGoStrSliceToTfSet(ctx, *resp.ToInputs)
		if err != nil {
			return
		}
	} else {
		ie.ToInputs = basetypes.NewSetNull(types.StringType)
	}

	return nil
}

func (ie *InputEdges) ToModelJSON(ctx context.Context) (jsonEdge client.ModelInputEdges, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, *ie, &jsonEdge)
	if err != nil {
		return
	}

	fromBases, ok := helper.TfSetStrToGoSlice(ctx, ie.FromBases)
	if !ok {
		err = fmt.Errorf("from_bases not found in tf from_bases")
		return
	}
	jsonEdge.FromBases = &fromBases

	fromInputs, ok := helper.TfSetStrToGoSlice(ctx, ie.FromInputs)
	if !ok {
		err = fmt.Errorf("from_inputs not found in tf from_inputs")
		return
	}
	jsonEdge.FromInputs = &fromInputs

	toSteps, ok := helper.TfSetStrToGoSlice(ctx, ie.ToSteps)
	if !ok {
		err = fmt.Errorf("to_steps not found in tf to_steps")
		return
	}
	jsonEdge.ToSteps = &toSteps

	toInputs, ok := helper.TfSetStrToGoSlice(ctx, ie.ToInputs)
	if !ok {
		err = fmt.Errorf("to_inputs not found in tf to_inputs")
		return
	}
	jsonEdge.ToInputs = &toInputs

	return
}

func (oe *OutputEdge) FillFromResp(ctx context.Context, resp client.ModelOutputEdges) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, oe)

	if resp.ToInputs != nil {
		oe.ToInputs, err = helper.FromGoStrSliceToTfSet(ctx, *resp.ToInputs)
		if err != nil {
			return
		}
	} else {
		oe.ToInputs = basetypes.NewSetNull(types.StringType)
	}

	return nil
}

func (oe *OutputEdge) ToModelJSON(ctx context.Context) (jsonEdge client.ModelOutputEdges, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, *oe, &jsonEdge)
	if err != nil {
		return
	}

	toInputs, ok := helper.TfSetStrToGoSlice(ctx, oe.ToInputs)
	if !ok {
		err = fmt.Errorf("to_inputs not found in tf to_inputs")
		return
	}
	jsonEdge.ToInputs = &toInputs

	return
}

func (pd *PipelineDag) FillFromResp(ctx context.Context, resp client.ModelPipelineDag) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, pd)

	var dg diag.Diagnostics

	// Convert input edges slice
	if resp.InputEdges != nil {
		inputEdges := make([]InputEdges, len(*resp.InputEdges))
		for i, edge := range *resp.InputEdges {
			if err = inputEdges[i].FillFromResp(ctx, edge); err != nil {
				return
			}
		}
		pd.InputEdges, dg = types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"node_id":     types.StringType,
				"from_bases":  types.SetType{ElemType: types.StringType},
				"from_inputs": types.SetType{ElemType: types.StringType},
				"to_steps":    types.SetType{ElemType: types.StringType},
				"to_inputs":   types.SetType{ElemType: types.StringType},
			},
		}, inputEdges)
		if dg.HasError() {
			err = fmt.Errorf("error converting input edges, err : %s", dg.Errors()[0])
			return
		}
	}

	// Convert output edges slice
	if resp.OutputEdges != nil {
		outputEdges := make([]OutputEdge, len(*resp.OutputEdges))
		for i, edge := range *resp.OutputEdges {
			if err = outputEdges[i].FillFromResp(ctx, edge); err != nil {
				return
			}
		}
		pd.OutputEdges, dg = types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"node_id":   types.StringType,
				"from_step": types.StringType,
				"to_inputs": types.SetType{ElemType: types.StringType},
			},
		}, outputEdges)
		if dg.HasError() {
			err = fmt.Errorf("error converting input edges, err : %s", dg.Errors()[0])
			return
		}
	}

	return nil
}

func (pd *PipelineDag) ToModelJSON(ctx context.Context) (jsonDag client.ModelPipelineDag, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, *pd, &jsonDag)
	if err != nil {
		return
	}

	// Convert input edges list
	var inputEdges []InputEdges
	if !pd.InputEdges.IsNull() && !pd.InputEdges.IsUnknown() {
		diags := pd.InputEdges.ElementsAs(ctx, &inputEdges, false)
		if diags.HasError() {
			err = fmt.Errorf("failed to convert input edges: %v", diags)
			return
		}

		apiInputEdges := make([]client.ModelInputEdges, len(inputEdges))
		for i, edge := range inputEdges {
			apiInputEdges[i], err = edge.ToModelJSON(ctx)
			if err != nil {
				return
			}
		}
		jsonDag.InputEdges = &apiInputEdges
	}

	// Convert output edges list
	var outputEdges []OutputEdge
	if !pd.OutputEdges.IsNull() && !pd.OutputEdges.IsUnknown() {
		diags := pd.OutputEdges.ElementsAs(ctx, &outputEdges, false)
		if diags.HasError() {
			err = fmt.Errorf("failed to convert output edges: %v", diags)
			return
		}

		apiOutputEdges := make([]client.ModelOutputEdges, len(outputEdges))
		for i, edge := range outputEdges {
			apiOutputEdges[i], err = edge.ToModelJSON(ctx)
			if err != nil {
				return
			}
		}
		jsonDag.OutputEdges = &apiOutputEdges
	}

	return
}

func inputEdgesAttrs() []BaseSchema {
	return []BaseSchema{
		{
			Name:     "node_id",
			AttrType: TfString,
			Required: true,
			Desc:     "id of the input node itself",
		},
		{
			Name:     "from_bases",
			AttrType: TFSet,
			Optional: true,
			SubType:  TfString,
			Desc:     "ids of in edges from basis nodes [the should not be of zero length]",
		},
		{
			Name:     "from_inputs",
			AttrType: TFSet,
			Optional: true,
			SubType:  TfString,
			Desc:     "id's of in edges from other inputs [the should not be of zero length]",
		},
		{
			Name:     "to_steps",
			AttrType: TFSet,
			Optional: true,
			SubType:  TfString,
			Desc:     "id's of out edges to step nodes [the should not be of zero length]",
		},
		{
			Name:     "to_inputs",
			AttrType: TFSet,
			Optional: true,
			SubType:  TfString,
			Desc:     "id's of out edges to other inputs [the should not be of zero length]",
		},
	}
}

func outputEdgeAttrs() []BaseSchema {
	return []BaseSchema{
		{
			Name:     "node_id",
			AttrType: TfString,
			Required: true,
			Desc:     "id of the output node",
		},
		{
			Name:     "from_step",
			AttrType: TfString,
			Required: true,
			Desc:     "id of step node producing the output",
		},
		{
			Name:     "to_inputs",
			AttrType: TFSet,
			Optional: true,
			SubType:  TfString,
			Desc:     "id's of input nodes the output is forwarded to [the should not be of zero length]",
		},
	}
}

const inputEdgesAttrDesc = "input edges in the pipeline DAG"

func dsInputEdgesSchema() dschema.SetNestedAttribute {
	return dschema.SetNestedAttribute{
		NestedObject: dschema.NestedAttributeObject{
			Attributes: DSAttributes(inputEdgesAttrs()),
		},
		Computed:            true,
		Description:         inputEdgesAttrDesc,
		MarkdownDescription: inputEdgesAttrDesc,
	}
}

func resInputEdgesSchema() rschema.SetNestedAttribute {
	return rschema.SetNestedAttribute{
		NestedObject: rschema.NestedAttributeObject{
			Attributes: ResAttributes(inputEdgesAttrs()),
		},
		Optional:            true,
		Description:         inputEdgesAttrDesc,
		MarkdownDescription: inputEdgesAttrDesc,
	}
}

const outputEdgesAttrDesc = "output edges in the pipeline DAG"

func dsOutputEdgesSchema() dschema.SetNestedAttribute {
	return dschema.SetNestedAttribute{
		NestedObject: dschema.NestedAttributeObject{
			Attributes: DSAttributes(outputEdgeAttrs()),
		},
		Computed:            true,
		Description:         outputEdgesAttrDesc,
		MarkdownDescription: outputEdgesAttrDesc,
	}
}

func resOutputEdgesSchema() rschema.SetNestedAttribute {
	return rschema.SetNestedAttribute{
		NestedObject: rschema.NestedAttributeObject{
			Attributes: ResAttributes(outputEdgeAttrs()),
		},
		Optional:            true,
		Description:         outputEdgesAttrDesc,
		MarkdownDescription: outputEdgesAttrDesc,
	}
}

func pipelineDagAttrs() []BaseSchema {
	return []BaseSchema{
		{
			Name:     "pipeline_id",
			AttrType: TfString,
			Required: true,
			Desc:     "id of the pipeline",
		},
	}
}

func PipelineDagDSchema() dschema.Schema {
	attrs := DSAttributes(pipelineDagAttrs())
	attrs["input_edges"] = dsInputEdgesSchema()
	attrs["output_edges"] = dsOutputEdgesSchema()

	return dschema.Schema{
		Attributes:  attrs,
		Description: "pipeline directed acyclic graph (DAG) structure",
	}
}

func PipelineDagResourceSchema() rschema.Schema {
	attrs := ResAttributes(pipelineDagAttrs())
	attrs["input_edges"] = resInputEdgesSchema()
	attrs["output_edges"] = resOutputEdgesSchema()

	return rschema.Schema{
		Attributes:  attrs,
		Description: "pipeline directed acyclic graph (DAG) structure",
	}
}
