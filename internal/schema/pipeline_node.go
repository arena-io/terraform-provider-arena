// Copyright (c) ArenaML Labs Pvt Ltd.

package schema

import (
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFPipelineNodes struct {
	PipelineID types.String `tfsdk:"id"`
	Inputs     types.Set    `tfsdk:"inputs"`
	Outputs    types.Set    `tfsdk:"outputs"`
	Steps      types.Set    `tfsdk:"steps"`
}

func InputNodeDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		Attributes: map[string]dschema.Attribute{
			"pipeline_id": dschema.StringAttribute{
				Required:    true,
				Description: "ID of the pipeline this node belongs to",
			},
		},
	}

}
