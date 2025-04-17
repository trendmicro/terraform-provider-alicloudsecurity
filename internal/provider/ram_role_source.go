package provider

import (
	"context"

	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ramRoleSource{}
	_ datasource.DataSourceWithConfigure = &ramRoleSource{}
)

// NewRamRoleSource is a helper function to simplify the provider implementation.
func NewRamRoleSource() datasource.DataSource {
	return &ramRoleSource{}
}

// ramRoleSource is the resource implementation.
type ramRoleSource struct {
	ram *ram.Client
}

// ramRoleSourceModel maps the resource schema.
type ramRoleSourceModel struct {
	RoleName types.String `tfsdk:"role_name"` // The name of the role.
}

// Metadata returns the resource type name.
func (r *ramRoleSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ram_role"
}

// Schema defines the schema for the resource.
func (r *ramRoleSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a RAM role.",
		Attributes: map[string]schema.Attribute{
			"role_name": schema.StringAttribute{
				Description: "The name of the role.",
				Required:    true,
				Computed:    false,
			},
		},
	}
}

// Configure prepares the provider for data source operations.
func (r *ramRoleSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	r.ram = clients.alicloudClients.Ram
}

// Read refreshes the Terraform state with the latest data.
func (r *ramRoleSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	panic("not implemented")
}
