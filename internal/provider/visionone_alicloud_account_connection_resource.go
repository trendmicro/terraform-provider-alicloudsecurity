package provider

import (
	"context"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
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
	ramClient *ram.Client
}

// visiononeAlicloudAccountConnectionResourceModel maps the resource schema.
type visiononeAlicloudAccountConnectionResourceModel struct {
	RoleId               types.String `tfsdk:"role_id"`
	RoleName             types.String `tfsdk:"role_name"`
	RoleArn              types.String `tfsdk:"role_arn"`
	RoleDescription      types.String `tfsdk:"role_description"`
	AlicloudAccessKey    types.String `tfsdk:"alicloud_accesskey"`
	AlicloudAccessSecret types.String `tfsdk:"alicloud_accesskey_secret"`
	AlicloudRegion       types.String `tfsdk:"alicloud_region"`
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
			"alicloud_accesskey": schema.StringAttribute{
				Description: "The Alicloud Access Key.",
				Required:    true,
				Sensitive:   true,
			},
			"alicloud_accesskey_secret": schema.StringAttribute{
				Description: "The Alicloud Access Key Secret.",
				Required:    true,
				Sensitive:   true,
			},
			"alicloud_region": schema.StringAttribute{
				Description: "The Alicloud region.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
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

	ramClient, err := r.GetRamClientInstance(ctx, plan.AlicloudAccessKey.ValueString(), plan.AlicloudAccessSecret.ValueString(), plan.AlicloudRegion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting RAM client", err.Error())
		return
	}

	// Create the resource in Alicloud
	createRoleRequest := &ram.CreateRoleRequest{
		RoleName: tea.String(plan.RoleName.ValueString()),
		AssumeRolePolicyDocument: tea.String(`{
	"Statement": [{
		"Action": "sts:AssumeRole",
		"Effect": "Allow",
		"Principal": {
			"Service": [
				"ecs.aliyuncs.com"
			]
		}
	}],
	"Version": "1"
}`),
		Description: tea.String(plan.RoleDescription.ValueString()),
	}
	createRoleResponse, err := ramClient.CreateRole(createRoleRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.RoleId = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.RoleId))
	plan.RoleArn = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.Arn))
	plan.RoleName = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.RoleName))
	plan.RoleDescription = types.StringValue(tea.StringValue(createRoleResponse.Body.Role.Description))
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

	ramClient, err := r.GetRamClientInstance(ctx, state.AlicloudAccessKey.ValueString(), state.AlicloudAccessSecret.ValueString(), state.AlicloudRegion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting RAM client", err.Error())
		return
	}
	// Fetch the role details
	roleDetailsRequest := &ram.GetRoleRequest{
		RoleName: tea.String(state.RoleName.ValueString()),
	}
	roleDetailsResponse, err := ramClient.GetRole(roleDetailsRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching role details", err.Error())
		return
	}
	// Overwrite items with refreshed state
	state.RoleId = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.RoleId))
	state.RoleArn = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.Arn))
	state.RoleName = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.RoleName))
	state.RoleDescription = types.StringValue(tea.StringValue(roleDetailsResponse.Body.Role.Description))
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

	ramClient, err := r.GetRamClientInstance(ctx, plan.AlicloudAccessKey.ValueString(), plan.AlicloudAccessSecret.ValueString(), plan.AlicloudRegion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting RAM client", err.Error())
		return
	}
	// Update the role name
	updateRoleRequest := &ram.UpdateRoleRequest{
		RoleName:       tea.String(plan.RoleName.ValueString()),
		NewDescription: tea.String(plan.RoleDescription.ValueString()),
	}
	updateRoleResponse, err := ramClient.UpdateRole(updateRoleRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating role", err.Error())
		return
	}

	plan.RoleId = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.RoleId))
	plan.RoleArn = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.Arn))
	plan.RoleName = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.RoleName))
	plan.RoleDescription = types.StringValue(tea.StringValue(updateRoleResponse.Body.Role.Description))

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

	ramClient, err := r.GetRamClientInstance(ctx, state.AlicloudAccessKey.ValueString(), state.AlicloudAccessSecret.ValueString(), state.AlicloudRegion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting RAM client", err.Error())
		return
	}
	// Delete the role
	deleteRoleRequest := &ram.DeleteRoleRequest{
		RoleName: tea.String(state.RoleName.ValueString()),
	}
	deleteRoleResponse, err := ramClient.DeleteRole(deleteRoleRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting role", err.Error())
		return
	}
	tflog.Info(ctx, "Role deleted successfully", map[string]any{
		"deletion_request_id": deleteRoleResponse.Body.RequestId,
	})
}

func (r *visiononeAlicloudAccountConnectionResource) GetRamClientInstance(
	ctx context.Context,
	alicloudAccessKey string,
	alicloudAccessSecret string,
	alicloudRegion string,
) (*ram.Client, error) {
	if r.ramClient != nil {
		return r.ramClient, nil
	}
	// configure the shared configuration
	config := &openapi.Config{
		AccessKeyId:     tea.String(alicloudAccessKey),
		AccessKeySecret: tea.String(alicloudAccessSecret),
		RegionId:        tea.String(alicloudRegion),
	}
	// Initialize RAM client
	config.Endpoint = tea.String("ram.aliyuncs.com")
	tflog.Info(ctx, "Creating Alicloud Resource Access Management client", map[string]any{
		"config": config,
	})
	client, err := ram.NewClient(config)
	if err != nil {
		return nil, err
	}
	r.ramClient = client
	return r.ramClient, nil
}
