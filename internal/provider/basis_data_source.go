// Copyright (c) ArenaML Labs Pvt Ltd.

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arena-ml/terraform-provider-arenaml/generator/client"
	"github.com/arena-ml/terraform-provider-arenaml/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type basisDataSource struct {
	cl *client.ClientWithResponses
}

var _ datasource.DataSource = (*basisDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*basisDataSource)(nil)

type basisDataSourceConfig struct {
	ID types.String `tfsdk:"id"`
}

func NewBasisDatasource() datasource.DataSource {
	return &basisDataSource{}
}

func (b *basisDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_basis"
}

func (b *basisDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.BasisDataSourceSchema(ctx)
}

func (b *basisDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	cl, ok := request.ProviderData.(*client.ClientWithResponses)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *oapi cleint, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}
	b.cl = cl

}

func (b *basisDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data schema.BasisModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := b.cl.GetBasisGetWithResponse(ctx, &client.GetBasisGetParams{Id: data.ID.ValueStringPointer()})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("API client error in GET Basis: id: %s \nerr: %s", data.ID.String(), err))
		return
	}

	if apiResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(fmt.Sprintf("Client Error : %d", apiResp.StatusCode()), fmt.Sprintf("Unable to get basis '%s\n code : %d'",
			data.ID.String(), apiResp.StatusCode()))
		return
	}

	if apiResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error : null data in response", fmt.Sprintf("API Call error \nid : '%s'", data.ID.String()))
		return
	}

	// copy the basic values from resp to data
	err = data.FillFromResp(ctx, *apiResp.JSON200)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("API response parsing error \nid : '%s' , err: %s", data.ID.String(), err.Error()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
