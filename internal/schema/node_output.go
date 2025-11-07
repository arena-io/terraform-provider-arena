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

type NodeOutput struct {
	PipelineID types.String `tfsdk:"pipeline_id"`
	ModelCommon
	StoreID types.String `tfsdk:"store_id"`
}

func (no *NodeOutput) FillFromResp(ctx context.Context, resp client.EntOutput) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, no)
	m := &ModelCommon{}
	helper.ConvertJSONStructToSimpleTF(ctx, resp, m)
	no.ModelCommon = *m

	return nil
}

func (no *NodeOutput) ToModelJSON(ctx context.Context) (client.EntOutput, error) {
	m := client.EntOutput{}
	err := helper.ConvertTfModelToApiJSON(ctx, no.ModelCommon, &m)
	if err != nil {
		return client.EntOutput{}, err
	}

	err = helper.ConvertTfModelToApiJSON(ctx, *no, &m)

	return m, err
}

func nodeOutputAttrs() []BaseSchema {
	attrs := giveCommonAttributes()
	outputAttrs := []BaseSchema{
		{
			Name:     "pipeline_id",
			AttrType: TfString,
			Required: true,
			Desc:     "id of the pipeline this output is part of",
		},
		{
			Name:     "store_id",
			AttrType: TfString,
			Optional: true,
			Desc:     "override for store to be used in place of default store for that env",
		},
	}

	outputAttrs = append(attrs, outputAttrs...)

	return outputAttrs
}

func NodeOutputDSchema() dschema.Schema {
	return dschema.Schema{
		Attributes:  DSAttributes(nodeOutputAttrs()),
		Description: "node of type output",
	}
}

func NodeOutputResourceSchema() rschema.Schema {
	return rschema.Schema{
		Attributes:  ResAttributes(nodeOutputAttrs()),
		Description: "node of type output",
	}
}
