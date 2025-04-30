package provider

import (
	"context"
	"terraform-provider-alicloudsecurity/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type connectedAccountSourceModel struct {
	AccountId      types.String `tfsdk:"account_id"`       // The ID of the AliCloud Account.
	RoleArn        types.String `tfsdk:"role_arn"`         // The ARN of the role in AliCloud Account. *required*
	OidcProviderId types.String `tfsdk:"oidc_provider_id"` // The ID of the OIDC provider in AliCloud Account. *required*
	Name           types.String `tfsdk:"name"`             // The name of the connected account in VisionOne. *required*
	Description    types.String `tfsdk:"description"`      // The description of the connected account in VisionOne

	ConnectionState types.String `tfsdk:"connection_state"`  // The state of the connected account in VisionOne
	CreatedDateTime types.String `tfsdk:"created_date_time"` // The creation time of the connected account in VisionOne
	UpdatedDateTime types.String `tfsdk:"updated_date_time"` // The last update time of the connected account in VisionOne
}

func (c *connectedAccountSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connected_account"
}

func (c *connectedAccountSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for connected account in VisionOne.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The ID of the connected AliCloud Account.",
				Required:    true,
				Computed:    false,
			},
			"role_arn": schema.StringAttribute{
				Description: "The ARN of the role in AliCloud Account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"oidc_provider_id": schema.StringAttribute{
				Description: "The ID of the OIDC provider in AliCloud Account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the connected account in VisionOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the connected account in VisionOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_state": schema.StringAttribute{
				Description: "The state of the connected account in VisionOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "The creation time of the connected account in VisionOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"updated_date_time": schema.StringAttribute{
				Description: "The last update time of the connected account in VisionOne.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
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
	// based on the schema, we need to get the account_id from the request
	var data connectedAccountSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// check if the account_id is set
	if data.AccountId.IsNull() || data.AccountId.IsUnknown() {
		resp.Diagnostics.AddError(
			"Account ID is required",
			"Account ID cannot be null or unknown.",
		)
		return
	}

	// Read the connected account from the API
	readConnectionResp, err := c.cam.ReadConnection(ctx, data.AccountId.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error",
			"Unable to read connected account: "+err.Error(),
		)
		return
	}
	// Check if the response is empty
	if readConnectionResp == nil {
		resp.Diagnostics.AddError(
			"Read Error",
			"Cannot find connected account with ID: "+data.AccountId.ValueString(),
		)
		return
	} else {
		// Map the response to the model
		if readConnectionResp.Id != nil {
			data.AccountId = types.StringValue(*readConnectionResp.Id)
		}
		if readConnectionResp.RoleArn != nil {
			data.RoleArn = types.StringValue(*readConnectionResp.RoleArn)
		}
		if readConnectionResp.OidcProviderId != nil {
			data.OidcProviderId = types.StringValue(*readConnectionResp.OidcProviderId)
		}
		if readConnectionResp.Name != nil {
			data.Name = types.StringValue(*readConnectionResp.Name)
		}
		if readConnectionResp.Description != nil {
			data.Description = types.StringValue(*readConnectionResp.Description)
		}
		if readConnectionResp.State != nil {
			data.ConnectionState = types.StringValue(*readConnectionResp.State)
		}
		if readConnectionResp.CreatedDateTime != nil {
			data.CreatedDateTime = types.StringValue(*readConnectionResp.CreatedDateTime)
		}
		if readConnectionResp.UpdatedDateTime != nil {
			data.UpdatedDateTime = types.StringValue(*readConnectionResp.UpdatedDateTime)
		}
	}

	// set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
