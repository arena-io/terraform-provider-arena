// Copyright (c) ArenaML Labs Pvt Ltd.

package schema

import (
	"context"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
	"github.com/arena-ml/terraform-provider-arenaml/helper"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Pipeline struct {
	ModelCommon
	Disabled types.Bool          `tfsdk:"disabled"`
	Paused   basetypes.BoolValue `tfsdk:"paused"`
}

type PipelineConfig struct {
}

func (c *PipelineConfig) FillFromResp(ctx context.Context, resp client.EntPipeline) (err error) {
	if resp.Config != nil {
		helper.ConvertJSONStructToSimpleTF(ctx, *resp.Config, c)
	}
	return nil
}

func (c *PipelineConfig) ToModelJSON(ctx context.Context) (jsonConf client.SchemaPipelineConfig, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, *c, &jsonConf)
	return
}

func (p *Pipeline) FillFromResp(ctx context.Context, resp client.EntPipeline) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, p)
	mc := &ModelCommon{}
	helper.ConvertJSONStructToSimpleTF(ctx, resp, mc)
	p.ModelCommon = *mc

	return nil
}

func (p *Pipeline) ToModelJSON(ctx context.Context) (jsonPipeline client.EntPipeline, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, p.ModelCommon, &jsonPipeline)
	if err != nil {
		return
	}
	err = helper.ConvertTfModelToApiJSON(ctx, *p, &jsonPipeline)

	var clPipelineConf client.SchemaPipelineConfig
	jsonPipeline.Config = &clPipelineConf

	return
}

func pipelineConfigAttrs() []BaseSchema {
	return []BaseSchema{}
}

const pipelineConfigAttrDesc = "configuration for the pipeline"

func dsPipelineConfigSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Attributes:          DSAttributes(pipelineConfigAttrs()),
		Computed:            true,
		Description:         pipelineConfigAttrDesc,
		MarkdownDescription: pipelineConfigAttrDesc,
	}
}

func resPipelineConfigSchema() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		Attributes:          ResAttributes(pipelineConfigAttrs()),
		Optional:            true,
		Description:         pipelineConfigAttrDesc,
		MarkdownDescription: pipelineConfigAttrDesc,
	}
}

func pipelineAttrs() []BaseSchema {
	attrs := giveCommonAttributes()
	pipelineAttrs := []BaseSchema{
		{
			Name:     "org_id",
			AttrType: TfString,
			Required: true,
			Desc:     "organization ID this pipeline belongs to",
		},
		{
			Name:     "disabled",
			AttrType: TfBoolean,
			Optional: true,
			Desc:     "disabled pipeline don't have any runs queued. Disabling a pipeline also cancels any queue or scheduled runs",
		},
		{
			Name:     "paused",
			AttrType: TfBoolean,
			Optional: true,
			Desc:     "paused pipeline have new run queued but not scheduled until unpaused",
		},
	}

	pipelineAttrs = append(attrs, pipelineAttrs...)

	return pipelineAttrs
}

func PipelineDSchema() dschema.Schema {
	attrs := DSAttributes(pipelineAttrs())

	return dschema.Schema{
		Attributes:  attrs,
		Description: "pipeline resource",
	}
}

func PipelineResourceSchema() rschema.Schema {
	attrs := ResAttributes(pipelineAttrs())

	return rschema.Schema{
		Attributes:  attrs,
		Description: "pipeline resource",
	}
}
