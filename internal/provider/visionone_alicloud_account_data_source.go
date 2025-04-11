package provider

import (
	"context"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	resourcemanager "github.com/alibabacloud-go/resourcemanager-20200331/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &visiononeAlicloudAccountDataSource{}
)

// NewVisiononeAlicloudAccountDataSource is a helper function to simplify the provider implementation.
func NewVisiononeAlicloudAccountDataSource() datasource.DataSource {
	return &visiononeAlicloudAccountDataSource{}
}

// visiononeAlicloudAccountDataSource is the data source implementation.
type visiononeAlicloudAccountDataSource struct {
	rmClient  *resourcemanager.Client
	ramClient *ram.Client
}

// visiononeAlicloudAccountDataSourceModel maps the data source schema.
type visiononeAlicloudAccountDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Email                types.String `tfsdk:"email"`
	AlicloudAccessKey    types.String `tfsdk:"alicloud_accesskey"`
	AlicloudAccessSecret types.String `tfsdk:"alicloud_accesskey_secret"`
	AlicloudRegion       types.String `tfsdk:"alicloud_region"`
}

// Metadata returns the data source type name.
func (d *visiononeAlicloudAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_visionone_alicloud_account"
}

// Schema defines the schema for the data source.
func (d *visiononeAlicloudAccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a VisionOne Alicloud account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the account.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the account.",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email of the account.",
				Computed:    true,
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

// Read fetches the account details using the provided ID.
func (d *visiononeAlicloudAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data visiononeAlicloudAccountDataSourceModel
	req.Config.Get(ctx, &data)
	// Fetch account details
	tflog.Info(ctx, "Reading Alicloud account data source", map[string]any{
		"data": data,
	})

	// configure the shared configuration
	config := &openapi.Config{
		AccessKeyId:     tea.String(data.AlicloudAccessKey.ValueString()),
		AccessKeySecret: tea.String(data.AlicloudAccessSecret.ValueString()),
		RegionId:        tea.String(data.AlicloudRegion.ValueString()),
	}

	// Initialize ResourceManager client
	config.Endpoint = tea.String("resourcemanager.aliyuncs.com")
	tflog.Info(ctx, "Creating Alicloud ResourceManager client", map[string]any{
		"config": config,
	})
	rm, err := resourcemanager.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Alicloud ResourceManager client", err.Error())
		return
	}
	d.rmClient = rm

	// Initialize RAM client
	config.Endpoint = tea.String("ram.aliyuncs.com")
	tflog.Info(ctx, "Creating Alicloud RAM client", map[string]any{
		"config": config,
	})
	ramClient, err := ram.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Alicloud RAM client", err.Error())
		return
	}
	d.ramClient = ramClient

	// Create the request to get account details

	getUserRequest := &ram.GetUserRequest{
		UserName: tea.String(data.Name.ValueString()),
	}

	getUserResponse, err := d.ramClient.GetUser(getUserRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching user details", err.Error())
		return
	}
	// Map response to data source model
	data.ID = types.StringValue(tea.StringValue(getUserResponse.Body.User.UserId))
	data.Name = types.StringValue(tea.StringValue(getUserResponse.Body.User.UserName))
	data.Email = types.StringValue(tea.StringValue(getUserResponse.Body.User.Email))

	// Set the state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
