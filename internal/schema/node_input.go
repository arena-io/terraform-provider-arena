// Copyright (c) ArenaML Labs Pvt Ltd.

package schema

import (
	"context"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
	"github.com/arena-ml/terraform-provider-arenaml/helper"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputNode struct {
	PipelineID types.String `tfsdk:"pipeline_id"`
	Kind       types.String `tfsdk:"kind"`
	Frozen     types.Bool   `tfsdk:"frozen"`
	StoreID    types.String `tfsdk:"store_id"`
	ModelCommon
}

func (ni *InputNode) FillFromResp(ctx context.Context, resp client.EntInput) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, ni)
	m := &ModelCommon{}
	helper.ConvertJSONStructToSimpleTF(ctx, resp, m)
	ni.ModelCommon = *m

	return nil
}

func (ni *InputNode) ToModelJSON(ctx context.Context) (client.EntInput, error) {
	m := client.EntInput{}
	// any embedded structs need to be converted separately
	err := helper.ConvertTfModelToApiJSON(ctx, ni.ModelCommon, &m)
	if err != nil {
		return m, err
	}
	err = helper.ConvertTfModelToApiJSON(ctx, *ni, &m)

	return m, err
}

func inputNodeAttrs() []BaseSchema {
	attrs := giveCommonAttributes()
	inputAttrs := []BaseSchema{
		{
			Name:     "pipeline_id",
			AttrType: TfString,
			Required: true,
			Desc:     "id of the pipeline this input is part of",
		},
		{
			Name:     "kind",
			AttrType: TfString,
			Optional: true,
			Desc:     "The type of input node.",
		},
		{
			Name:     "frozen",
			AttrType: TfBoolean,
			Optional: true,
			Desc:     "frozen input it will pin to its last artifact for runs and ignore new artifacts",
		},
		{
			Name:     "store_id",
			AttrType: TfString,
			Optional: true,
			Desc:     "override for store to be used in place of default store for that env",
		},
	}

	inputAttrs = append(attrs, inputAttrs...)

	return inputAttrs
}

func InputNodeDSchema() dschema.Schema {
	return dschema.Schema{
		Attributes:  DSAttributes(inputNodeAttrs()),
		Description: "node of type input",
	}
}

func InputNodeResourceSchema() rschema.Schema {
	return rschema.Schema{
		Attributes:  ResAttributes(inputNodeAttrs()),
		Description: "node of type input",
	}
}
