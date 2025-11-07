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

var orgJson client.EntOrg

type Org struct {
	ModelCommon
	Inactive types.Bool `tfsdk:"inactive"`
}

func (o *Org) FillFromResp(ctx context.Context, resp client.EntOrg) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, o)
	mc := &ModelCommon{}
	helper.ConvertJSONStructToSimpleTF(ctx, resp, mc)
	o.ModelCommon = *mc

	return nil
}

func (o *Org) ToModelJSON(ctx context.Context) (jsonOrg client.EntOrg, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, o.ModelCommon, &jsonOrg)
	if err != nil {
		return
	}
	err = helper.ConvertTfModelToApiJSON(ctx, o, &jsonOrg)

	return
}

func orgAttrs() []BaseSchema {
	attrs := giveCommonAttributes()
	orgAttrs := []BaseSchema{
		{
			Name:     "inactive",
			AttrType: TfBoolean,
			Optional: true,
			Desc:     "whether the organization is inactive",
		},
	}

	orgAttrs = append(attrs, orgAttrs...)

	return orgAttrs
}

func OrgDSchema() dschema.Schema {
	attrs := DSAttributes(orgAttrs())

	return dschema.Schema{
		Attributes:  attrs,
		Description: "organization resource",
	}
}

func OrgResourceSchema() rschema.Schema {
	attrs := ResAttributes(orgAttrs())

	return rschema.Schema{
		Attributes:  attrs,
		Description: "organization resource",
	}
}
