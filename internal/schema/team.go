// Copyright (c) ArenaML Labs Pvt Ltd.

package schema

import (
	"context"
	"fmt"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
	"github.com/arena-ml/terraform-provider-arenaml/helper"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kr/pretty"
)

var t client.EntTeam

type Team struct {
	ModelCommon
	OrgId    types.String `tfsdk:"org_id"`
	Inactive types.Bool   `tfsdk:"inactive"`
	Role     types.String `tfsdk:"role"`
	Config   *TeamConfig  `tfsdk:"config"`
}

type TeamConfig struct {
	SyncWith jsontypes.Normalized `tfsdk:"sync_with"`
}

func (c *TeamConfig) FillFromResp(ctx context.Context, resp client.EntTeam) (err error) {
	tflog.Warn(ctx, fmt.Sprintf("\n\n%s\n\n", pretty.Sprint(c)))
	if resp.Config == nil {
		return nil
	}

	conf := *resp.Config
	helper.ConvertJSONStructToSimpleTF(ctx, conf, c)

	if conf.SyncWith != nil {
		c.SyncWith, err = helper.JSONObjToNormalized(*conf.SyncWith)
	} else {
		c.SyncWith = jsontypes.NewNormalizedNull()
	}

	tflog.Warn(ctx, fmt.Sprintf("\n to JSON\n%s\n\n", pretty.Sprint(c, resp.Config)))
	return err
}

func (c *TeamConfig) ToModelJson(ctx context.Context) (jsonConf client.SchemaTeamConfig, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, *c, &jsonConf)

	if !c.SyncWith.IsNull() && !c.SyncWith.IsUnknown() {
		syncWith, err := helper.TfJSONToGoMapInterface(ctx, c.SyncWith)
		if err != nil {
			err = fmt.Errorf("sync_with not found in tf data: \n %s %s ", c.SyncWith, err)
			return jsonConf, err
		}
		jsonConf.SyncWith = &syncWith
	}

	tflog.Warn(ctx, fmt.Sprintf("\nto TF\n%s\n\n", pretty.Sprint(c, jsonConf)))
	return
}

func (t *Team) FillFromResp(ctx context.Context, resp client.EntTeam) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, resp, t)
	mc := &ModelCommon{}
	helper.ConvertJSONStructToSimpleTF(ctx, resp, mc)
	t.ModelCommon = *mc

	if resp.Config != nil {
		teamConf := &TeamConfig{}
		err = teamConf.FillFromResp(ctx, resp)
		if err != nil {
			return
		}
		t.Config = teamConf
	}

	return nil
}

func (t *Team) ToModelJSON(ctx context.Context) (jsonTeam client.EntTeam, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, t.ModelCommon, &jsonTeam)
	if err != nil {
		return
	}
	err = helper.ConvertTfModelToApiJSON(ctx, t, &jsonTeam)

	var clTeamConf client.SchemaTeamConfig
	clTeamConf, err = t.Config.ToModelJson(ctx)
	if err != nil {
		return
	}
	jsonTeam.Config = &clTeamConf

	return
}

func teamConfigAttrs() []BaseSchema {
	return []BaseSchema{
		{
			Name:     "sync_with",
			AttrType: TfJSON,
			Optional: true,
			Desc:     "configuration for autosync with external oauth group",
		},
	}
}

const teamConfigAttrDesc = "configuration for the team"

func dsTeamConfigSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Attributes:          DSAttributes(teamConfigAttrs()),
		Computed:            true,
		Description:         teamConfigAttrDesc,
		MarkdownDescription: teamConfigAttrDesc,
	}
}

func resTeamConfigSchema() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		Attributes:          ResAttributes(teamConfigAttrs()),
		Optional:            true,
		Description:         teamConfigAttrDesc,
		MarkdownDescription: teamConfigAttrDesc,
	}
}

func teamAttrs() []BaseSchema {
	attrs := giveCommonAttributes()
	teamAttrs := []BaseSchema{
		{
			Name:     "org_id",
			AttrType: TfString,
			Required: true,
			Desc:     "organization id this team belongs to",
		},
		{
			Name:     "inactive",
			AttrType: TfBoolean,
			Optional: true,
			Desc:     "whether the team is inactive",
		},
		{
			Name:     "role",
			AttrType: TfString,
			Optional: true,
			Desc:     "role of the team",
		},
	}

	teamAttrs = append(attrs, teamAttrs...)

	return teamAttrs
}

func TeamDSchema() dschema.Schema {
	attrs := DSAttributes(teamAttrs())
	attrs["config"] = dsTeamConfigSchema()

	return dschema.Schema{
		Attributes:  attrs,
		Description: "team resource",
	}
}

func TeamResourceSchema() rschema.Schema {
	attrs := ResAttributes(teamAttrs())
	attrs["config"] = resTeamConfigSchema()

	return rschema.Schema{
		Attributes:  attrs,
		Description: "team resource",
	}
}
