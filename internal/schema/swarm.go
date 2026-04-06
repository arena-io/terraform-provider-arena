package schema

import (
	"context"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
	"github.com/arena-ml/terraform-provider-arenaml/helper"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SwarmModel struct {
	ModelCommon
	Kind     types.String `tfsdk:"profile_id"`
	Inactive types.Bool   `tfsdk:"inactive"`
	Spec     *SwarmSpec   `tfsdk:"swarm_id"`
	// TODO add drone members field
}

type SwarmSpec struct {
	MaxUnits          types.Int32   `tfsdk:"max_units"`
	ActiveQuorumRatio types.Float32 `json:"active_quorum"` // min percentage of units to be connected for an active fleet e.g., 0.5 for 50%
}

func (s *SwarmModel) FillFromResp(ctx context.Context, r client.EntSwarm) (err error) {
	helper.ConvertJSONStructToSimpleTF(ctx, r, s)

	mc := &ModelCommon{}
	helper.ConvertJSONStructToSimpleTF(ctx, r, mc)
	s.ModelCommon = *mc

	if r.Spec != nil {
		spec := &SwarmSpec{}
		helper.ConvertJSONStructToSimpleTF(ctx, *r.Spec, spec)
		s.Spec = spec
	}

	return
}

func (s *SwarmModel) ToModelJSON(ctx context.Context) (sw client.EntSwarm, err error) {
	err = helper.ConvertTfModelToApiJSON(ctx, *s, sw)
	if err != nil {
		return
	}

	spec := &client.SchemaSwarmSpec{}
	if s.Spec != nil {
		err = helper.ConvertTfModelToApiJSON(ctx, s.Spec, spec)
		if err != nil {
			return
		}
	}

	return
}

func SwarmAttrs() []BaseSchema {
	return append(giveCommonAttributes(),
		BaseSchema{Name: "kind", AttrType: TfString, Optional: true, Desc: "drone kind"},
		BaseSchema{Name: "inactive", AttrType: TfBoolean, Optional: true, Desc: "whether the drone is inactive"},
	)
}
