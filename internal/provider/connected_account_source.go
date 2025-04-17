package provider

import (
	"context"
	"terraform-provider-alicloudsecurity/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &connectedAccountSource{}
	_ datasource.DataSourceWithConfigure = &connectedAccountSource{}
)

func NewConnectedAccountSource() datasource.DataSource {
	return &connectedAccountSource{}
}

type connectedAccountSource struct {
	cam *common.CamClient
}

func (c *connectedAccountSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connected_account"
}

func (c *connectedAccountSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// Description: "Data source for connected account in VisionOne.",
		// Attributes: map[string]schema.Attribute{
		// 	"id": schema.StringAttribute{
		// 		Description: "The ID of the AliCloud Account.",
		// 		Required:    true,
		// 		Computed:    false,
		// 	},
		// },
	}
}

// Configure prepares the provider for data source operations.
func (c *connectedAccountSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clients := req.ProviderData.(*aliCloudSecurityProviderClients)
	if clients == nil {
		resp.Diagnostics.AddError(
			"Client Error",
			"Client configuration is not set up properly. Please configure the provider.",
		)
		return
	}
	c.cam = clients.visiononeClients.Cam
}

func (c *connectedAccountSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	panic("not implemented")
}
