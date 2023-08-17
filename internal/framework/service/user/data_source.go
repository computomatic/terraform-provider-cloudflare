package user

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CloudflareUserDataSource{}

func NewDataSource() datasource.DataSource {
	return &CloudflareUserDataSource{}
}

type CloudflareUserDataSource struct {
	client *cloudflare.API
}

func (r *CloudflareUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *CloudflareUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cloudflare.API)

	if !ok {
		resp.Diagnostics.AddError(
			"unexpected resource configure type",
			fmt.Sprintf("expected *cloudflare.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *CloudflareUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudflareUserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user, err := r.client.UserDetails(ctx)
	if err != nil {
		resp.Diagnostics.AddError("unable to retrieve user details", err.Error())
		return
	}
	data.ID = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.Username = types.StringValue(user.Username)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
