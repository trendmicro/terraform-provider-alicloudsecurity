package provider

import (
	"context"
	"fmt"

	"terraform-provider-alicloudsecurity/internal/common"

	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource = &visiononeAlicloudAccountConnectionResource{}
)

// NewVisiononeAlicloudAccountConnectionResource is a helper function to simplify the provider implementation.
func NewVisiononeAlicloudAccountConnectionResource() resource.Resource {
	return &visiononeAlicloudAccountConnectionResource{}
}

// visiononeAlicloudAccountConnectionResource is the resource implementation.
type visiononeAlicloudAccountConnectionResource struct {
	cam *common.CamClient
	sts *sts.Client
	ram *ram.Client
}

// visiononeAlicloudAccountConnectionResourceModel maps the resource schema.
type visiononeAlicloudAccountConnectionResourceModel struct {
	RoleId          types.String `tfsdk:"role_id"`
	RoleName        types.String `tfsdk:"role_name"`
	RoleArn         types.String `tfsdk:"role_arn"`
	RoleDescription types.String `tfsdk:"role_description"`
}

// Metadata returns the resource type name.
func (r *visiononeAlicloudAccountConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_visionone_alicloud_account_connection"
}

// Schema defines the schema for the resource.
func (r *visiononeAlicloudAccountConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a VisionOne Alicloud account connection.",
		Attributes: map[string]schema.Attribute{
			"role_id": schema.StringAttribute{
				Description: "The ID of the role.",
				Computed:    true,
			},
			"role_name": schema.StringAttribute{
				Description: "The name of the role.",
				Required:    true,
			},
			"role_arn": schema.StringAttribute{
				Description: "The ARN of the role.",
				Computed:    true,
			},
			"role_description": schema.StringAttribute{
				Description: "The description of the role.",
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *visiononeAlicloudAccountConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring resource...")

	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	clients, ok := req.ProviderData.(*aliCloudSecurityProviderClients)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aliCloudSecurityProviderClients, got %T.", req.ProviderData),
		)
		return
	}
	// Set the client
	r.cam = clients.visiononeClients.Cam
	r.sts = clients.alicloudClients.Sts
	r.ram = clients.alicloudClients.Ram
}

// Create creates the resource.
func (r *visiononeAlicloudAccountConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "resource creating...")

	// Retrieve the values from plan
	var plan visiononeAlicloudAccountConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to schema and populate Computed attribute values
	// plan.RoleId = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.RoleId))
	// plan.RoleArn = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.Arn))
	// plan.RoleName = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.RoleName))
	// plan.RoleDescription = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.Description))

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *visiononeAlicloudAccountConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "resource reading...")
	// Get the current state
	var state visiononeAlicloudAccountConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the values from plan parameters, and connect to vision-one-cam backend
	panic("vision-one-cam backend is not implemented yet")

	// Overwrite items with refreshed state
	// state.RoleId = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.RoleId))
	// state.RoleArn = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.Arn))
	// state.RoleName = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.RoleName))
	// state.RoleDescription = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.Description))

	// Set the state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *visiononeAlicloudAccountConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "resource updating...")

	// Retrieve values from plan
	var plan visiononeAlicloudAccountConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the values from plan parameters, and update to vision-one-cam backend
	panic("vision-one-cam backend is not implemented yet")

	// Set the values to plan
	// plan.RoleId = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.RoleId))
	// plan.RoleArn = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.Arn))
	// plan.RoleName = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.RoleName))
	// plan.RoleDescription = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.Description))

	// Set the state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *visiononeAlicloudAccountConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "resource deleting...")

	// Retrieve values from state
	var state visiononeAlicloudAccountConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the values from state parameters, and request vision-one-cam backend to delete the record
	panic("vision-one-cam backend is not implemented yet")
}
