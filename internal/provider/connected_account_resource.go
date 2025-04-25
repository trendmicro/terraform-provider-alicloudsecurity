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
	StackStateRegion types.String `tfsdk:"stack_state_region"` // The region of the AliCloud Account. *required* Alias: alibabaRegion, parentStackRegion
	AccountId        types.String `tfsdk:"account_id"`         // The ID of the AliCloud Account.
	RoleArn          types.String `tfsdk:"role_arn"`           // The ARN of the role in AliCloud Account. *required*
	OidcProviderId   types.String `tfsdk:"oidc_provider_id"`   // The ID of the OIDC provider in AliCloud Account. *required*
	Name             types.String `tfsdk:"name"`               // The name of the connected account in VisionOne. *required*
	Description      types.String `tfsdk:"description"`        // The description of the connected account in VisionOne

	ConnectionState types.String `tfsdk:"connection_state"`  // The state of the connected account in VisionOne
	CreatedDateTime types.String `tfsdk:"created_date_time"` // The creation time of the connected account in VisionOne
	UpdatedDateTime types.String `tfsdk:"updated_date_time"` // The last update time of the connected account in VisionOne
}

type ConnectedSecurityServiceModel struct {
	Name        types.String `tfsdk:"name"`         // The name of the connected security service
	InstanceIds types.List   `tfsdk:"instance_ids"` // The list of instance IDs of the connected security service
}

// Metadata returns the resource type name.
func (r *connectedAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connected_account"
}

// Schema defines the schema for the resource.
func (r *connectedAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The resource schema for connected account.",
		Attributes: map[string]schema.Attribute{
			"stack_state_region": schema.StringAttribute{
				Description: "The region of the AliCloud Account where the terraform state is located. *required*", // example: us-west-1
				Required:    true,
			},
			"account_id": schema.StringAttribute{
				Description: "The ID of the AliCloud Account.",
				Required:    true,
			},
			"role_arn": schema.StringAttribute{
				Description: "The ARN of the role in AliCloud Account. *required*",
				Required:    true,
			},
			"oidc_provider_id": schema.StringAttribute{
				Description: "The ID of the OIDC provider in AliCloud Account. *required*",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the connected account in VisionOne. *required*",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the connected account in VisionOne",
				Optional:    true,
				Computed:    true,
			},
			"connection_state": schema.StringAttribute{
				Description: "The state of the connected account in VisionOne",
				Optional:    true,
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "The creation time of the connected account in VisionOne",
				Optional:    true,
				Computed:    true,
			},
			"updated_date_time": schema.StringAttribute{
				Description: "The last update time of the connected account in VisionOne",
				Optional:    true,
				Computed:    true,
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

	createConnectionReq := &common.CreateConnectionRequest{
		AlibabaRegion:  plan.StackStateRegion.ValueStringPointer(),
		RoleArn:        plan.RoleArn.ValueStringPointer(),
		OidcProviderId: plan.OidcProviderId.ValueStringPointer(),
		Name:           plan.Name.ValueStringPointer(),
		Description:    plan.Description.ValueStringPointer(),
	}
	err := r.cam.CreateConnection(createConnectionReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create Connection Error",
			"Failed to create connection: "+err.Error(),
		)
		return
	}

	readConnectionResp := &common.ReadConnectionResponse{}
	readConnectionResp, err = r.cam.ReadConnection(plan.AccountId.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Connection Error",
			"Failed to read connection: "+err.Error(),
		)
		return
	}
	if readConnectionResp == nil {
		resp.Diagnostics.AddError(
			"Read Connection Error",
			"Failed to read connection: response is nil",
		)
		return
	} else {
		// Overwrite the plan with the read response
		plan.AccountId = types.StringValue(*readConnectionResp.Id)
		plan.StackStateRegion = types.StringValue(*readConnectionResp.ParentStackRegion)
		plan.RoleArn = types.StringValue(*readConnectionResp.RoleArn)
		plan.OidcProviderId = types.StringValue(*readConnectionResp.OidcProviderId)
		plan.Name = types.StringValue(*readConnectionResp.Name)
		plan.Description = types.StringValue(*readConnectionResp.Description)
		plan.ConnectionState = types.StringValue(*readConnectionResp.State)
		plan.CreatedDateTime = types.StringValue(*readConnectionResp.CreatedDateTime)
		plan.UpdatedDateTime = types.StringValue(*readConnectionResp.UpdatedDateTime)
	}

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *connectedAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get the current state
	var state connectedAccountResourceModel
	diags := req.State.Get(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Get refreshed data from the API
	readConnectionResp, err := r.cam.ReadConnection(state.AccountId.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Connection Error",
			"Failed to read connection: "+err.Error(),
		)
		return
	}
	if readConnectionResp == nil {
		resp.Diagnostics.AddError(
			"Read Connection Error",
			"Failed to read connection: response is nil",
		)
		return
	} else {
		// Overwrite the state with the read response
		state.AccountId = types.StringValue(*readConnectionResp.Id)
		state.StackStateRegion = types.StringValue(*readConnectionResp.ParentStackRegion)
		state.RoleArn = types.StringValue(*readConnectionResp.RoleArn)
		state.OidcProviderId = types.StringValue(*readConnectionResp.OidcProviderId)
		state.Name = types.StringValue(*readConnectionResp.Name)
		state.Description = types.StringValue(*readConnectionResp.Description)
		state.ConnectionState = types.StringValue(*readConnectionResp.State)
		state.CreatedDateTime = types.StringValue(*readConnectionResp.CreatedDateTime)
		state.UpdatedDateTime = types.StringValue(*readConnectionResp.UpdatedDateTime)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update modifies the existing resource and sets the updated state.
func (r *connectedAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the values from plan
	var plan connectedAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Update the connection
	updateConnectionReq := &common.UpdateConnectionRequest{
		Name:        plan.Name.ValueStringPointer(),
		Description: plan.Description.ValueStringPointer(),
	}
	err := r.cam.UpdateConnection(plan.AccountId.ValueStringPointer(), updateConnectionReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update Connection Error",
			"Failed to update connection: "+err.Error(),
		)
		return
	}

	// Read the updated connection
	readConnectionResp, err := r.cam.ReadConnection(plan.AccountId.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Connection Error",
			"Failed to read connection: "+err.Error(),
		)
		return
	}
	if readConnectionResp == nil {
		resp.Diagnostics.AddError(
			"Read Connection Error",
			"Failed to read connection: response is nil",
		)
		return
	} else {
		// Overwrite the plan with the read response
		plan.AccountId = types.StringValue(*readConnectionResp.Id)
		plan.StackStateRegion = types.StringValue(*readConnectionResp.ParentStackRegion)
		plan.RoleArn = types.StringValue(*readConnectionResp.RoleArn)
		plan.OidcProviderId = types.StringValue(*readConnectionResp.OidcProviderId)
		plan.Name = types.StringValue(*readConnectionResp.Name)
		plan.Description = types.StringValue(*readConnectionResp.Description)
		plan.ConnectionState = types.StringValue(*readConnectionResp.State)
		plan.CreatedDateTime = types.StringValue(*readConnectionResp.CreatedDateTime)
		plan.UpdatedDateTime = types.StringValue(*readConnectionResp.UpdatedDateTime)
	}

	// Set state to fully populated plan
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete removes the resource from the state.
func (r *connectedAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve the values from state
	var state connectedAccountResourceModel
	diags := req.State.Get(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Delete the connection
	err := r.cam.DeleteConnection(state.AccountId.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Connection Error",
			"Failed to delete connection: "+err.Error(),
		)
		return
	}

	// Optionally, you can set the state to nil or empty
	resp.State.RemoveResource(ctx)
	if resp.Diagnostics.HasError() {
		return
	}
}

// ImportState imports the resource state.
func (r *connectedAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("not implemented")
}
