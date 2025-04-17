package provider

import (
	"context"
	"terraform-provider-alicloudsecurity/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &connectedAccountResource{}
	_ resource.ResourceWithConfigure   = &connectedAccountResource{}
	_ resource.ResourceWithImportState = &connectedAccountResource{}
)

// NewConnectedAccountResource is a helper function to simplify the provider implementation.
func NewConnectedAccountResource() resource.Resource {
	return &connectedAccountResource{}
}

// connectedAccountResource is the resource implementation.
type connectedAccountResource struct {
	cam *common.CamClient
}

// connectedAccountResourceModel maps the resource schema.
type connectedAccountResourceModel struct {
	AccountId         types.String `tfsdk:"account_id"`          // The ID of the AliCloud Account
	RoleArn           types.String `tfsdk:"role_arn"`            // The ARN of the role in AliCloud Account
	OidcProviderId    types.String `tfsdk:"oidc_provider_id"`    // The ID of the OIDC provider in AliCloud Account
	ParentStackRegion types.String `tfsdk:"parent_stack_region"` // The region of the parent stack in VisionOne
	Name              types.String `tfsdk:"name"`                // The name of the connected account in VisionOne
	Description       types.String `tfsdk:"description"`         // The description of the connected account in VisionOne
	State             types.String `tfsdk:"state"`               // The state of the connected account in VisionOne
	CreatedDateTime   types.String `tfsdk:"created_date_time"`   // The creation time of the connected account in VisionOne
	UpdatedDateTime   types.String `tfsdk:"updated_date_time"`   // The last update time of the connected account in VisionOne
}

// Metadata returns the resource type name.
func (r *connectedAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connected_account"
}

// Schema defines the schema for the resource.
func (r *connectedAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a VisionOne connected account.",
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Description: "The ID of the AliCloud Account.",
				Required:    true,
				Computed:    false,
			},
			"role_arn": schema.StringAttribute{
				Description: "The ARN of the role in AliCloud Account.",
				Required:    true,
				Computed:    false,
			},
			"oidc_provider_id": schema.StringAttribute{
				Description: "The ID of the OIDC provider in AliCloud Account.",
				Required:    true,
				Computed:    false,
			},
			"parent_stack_region": schema.StringAttribute{
				Description: "The region of the parent stack in VisionOne.",
				Required:    true,
				Computed:    false,
			},
			"name": schema.StringAttribute{
				Description: "The name of the connected account in VisionOne.",
				Optional:    true,
				Computed:    false,
			},
			"description": schema.StringAttribute{
				Description: "The description of the connected account in VisionOne.",
				Optional:    true,
				Computed:    false,
			},
			"state": schema.StringAttribute{
				Description: "The state of the connected account in VisionOne.",
				Optional:    true,
				Computed:    false,
			},
			"created_date_time": schema.StringAttribute{
				Description: "The creation time of the connected account in VisionOne.",
				Optional:    true,
				Computed:    false,
			},
			"updated_date_time": schema.StringAttribute{
				Description: "The last update time of the connected account in VisionOne.",
				Optional:    true,
				Computed:    false,
			},
		},
	}
}

// Configure prepares the provider for data source operations.
func (r *connectedAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.cam = clients.visiononeClients.Cam
}

// Create creates the resource and sets the initial state.
func (r *connectedAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the values from plan
	var plan connectedAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	panic("not implemented")
}

// Read refreshes the Terraform state with the latest data.
func (r *connectedAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	panic("not implemented")
}

// Update modifies the existing resource and sets the updated state.
func (r *connectedAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("not implemented")
}

// Delete removes the resource from the state.
func (r *connectedAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	panic("not implemented")
}

// ImportState imports the resource state.
func (r *connectedAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("not implemented")
}
