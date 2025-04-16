package provider

import (
	"context"
	"fmt"
	"os"

	"terraform-provider-alicloudsecurity/internal/common"

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
	v1Client  *common.VisionOneClient
	ramClient *ram.Client
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

	client, ok := req.ProviderData.(*common.VisionOneClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			"Expected *VisionOneClient, got something else",
		)
		return
	}
	err := client.HealthCheck()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error configuring VisionOne client",
			fmt.Sprintf("Error configuring VisionOne client: %s", err.Error()),
		)
		return
	}

	// Set the client
	r.v1Client = client
	tflog.Info(ctx, "VisionOne client configured successfully")
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

func (r *visiononeAlicloudAccountConnectionResource) GetRamClientInstance(
	ctx context.Context,
) (*ram.Client, error) {
	if r.ramClient != nil {
		return r.ramClient, nil
	}

	// Obtain the key/secret/region from the environment variables
	alicloudAccessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	alicloudAccessSecret := os.Getenv("ALICLOUD_ACCESS_SECRET")
	alicloudRegion := os.Getenv("ALICLOUD_REGION")
	if alicloudAccessKey == "" || alicloudAccessSecret == "" || alicloudRegion == "" {
		return nil, fmt.Errorf("ALICLOUD_ACCESS_KEY, ALICLOUD_ACCESS_SECRET, and ALICLOUD_REGION environment variables must be set")
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
	// Check if the client is not nil and can be used
	if client == nil {
		return nil, fmt.Errorf("failed to create RAM client")
	} else {
		tflog.Info(ctx, "Alicloud Resource Access Management client created successfully")
	}

	r.ramClient = client
	return r.ramClient, nil
}
