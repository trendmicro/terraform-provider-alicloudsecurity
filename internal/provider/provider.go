// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-alicloudsecurity/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &aliCloudSecurityProvider{}
	_ provider.ProviderWithFunctions = &aliCloudSecurityProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &aliCloudSecurityProvider{
			version: version,
		}
	}
}

// aliCloudSecurityProvider is the provider implementation.
type aliCloudSecurityProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type aliCloudSecurityProviderClients struct {
	visiononeClients *common.VisionOneClients
	alicloudClients  *common.AliCloudClients
}

// aliCloudSecurityProviderModel maps provider schema data to a Go type.
type aliCloudSecurityProviderModel struct {
	VisiononeEndpoint     types.String `tfsdk:"visionone_endpoint"`
	VisiononeEndpointType types.String `tfsdk:"visionone_endpoint_type"`
	VisiononeBusinessId   types.String `tfsdk:"visionone_business_id"`
	VisiononeAPIKey       types.String `tfsdk:"visionone_api_key"`
	VisiononeRegion       types.String `tfsdk:"visionone_region"`
}

// Metadata returns the provider type name.
func (p *aliCloudSecurityProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "alicloudsecurity"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *aliCloudSecurityProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with VisionOne AliCloud Security.",
		Attributes: map[string]schema.Attribute{
			"visionone_endpoint": schema.StringAttribute{
				Description: "Endpoint for VisionOne AliCloud Security. May also be provided via VISIONONE_ENDPOINT environment variable.",
				Optional:    true,
			},
			"visionone_endpoint_type": schema.StringAttribute{
				Description: "Endpoint type for VisionOne AliCloud Security. May also be provided via VISIONONE_ENDPOINT_TYPE environment variable.",
				Optional:    true,
			},
			"visionone_business_id": schema.StringAttribute{
				Description: "Bussiness Id for VisionOne AliCloud Security. May also be provided via VISIONONE_BUSINESS_ID environment variable.",
				Optional:    true,
			},
			"visionone_api_key": schema.StringAttribute{
				Description: "API key for VisionOne AliCloud Security. May also be provided via VISIONONE_API_KEY environment variable.",
				Optional:    true,
			},
			"visionone_region": schema.StringAttribute{
				Description: "Region for VisionOne AliCloud Security. May also be provided via VISIONONE_REGION environment variable.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a VisionOne AliCloud Security API client for data sources and resources.
func (p *aliCloudSecurityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring VisionOne AliCloud Security client")

	// Retrieve provider data from configuration
	var config aliCloudSecurityProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.VisiononeEndpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_endpoint"),
			"Unknown VisionOne Endpoint",
			"The provider cannot create the VisionOne API client as there is an unknown configuration value for the VisionOne endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VISIONONE_ENDPOINT environment variable.",
		)
	}

	if config.VisiononeEndpointType.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_endpoint_type"),
			"Unknown VisionOne Endpoint Type",
			"The provider cannot create the VisionOne API client as there is an unknown configuration value for the VisionOne endpoint type. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VISIONONE_ENDPOINT_TYPE environment variable.",
		)
	}

	if config.VisiononeBusinessId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_business_id"),
			"Unknown VisionOne Business Id",
			"The provider cannot create the VisionOne API client as there is an unknown configuration value for the VisionOne business id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VISIONONE_BUSINESS_ID environment variable.",
		)
	}

	if config.VisiononeAPIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_api_key"),
			"Unknown VisionOne API Key",
			"The provider cannot create the VisionOne API client as there is an unknown configuration value for the VisionOne API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VISIONONE_API_KEY environment variable.",
		)
	}

	if config.VisiononeRegion.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_region"),
			"Unknown VisionOne Region",
			"The provider cannot create the VisionOne API client as there is an unknown configuration value for the VisionOne region. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VISIONONE_REGION environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	visionone_endpoint := os.Getenv("VISIONONE_ENDPOINT")
	visionone_endpoint_type := os.Getenv("VISIONONE_ENDPOINT_TYPE")
	visionone_business_id := os.Getenv("VISIONONE_BUSINESS_ID")
	visionone_api_key := os.Getenv("VISIONONE_API_KEY")
	visionone_region := os.Getenv("VISIONONE_REGION")

	if !config.VisiononeEndpoint.IsNull() {
		visionone_endpoint = config.VisiononeEndpoint.ValueString()
	}

	if !config.VisiononeEndpointType.IsNull() {
		visionone_endpoint_type = config.VisiononeEndpointType.ValueString()
	}

	if !config.VisiononeBusinessId.IsNull() {
		visionone_business_id = config.VisiononeBusinessId.ValueString()
	}

	if !config.VisiononeAPIKey.IsNull() {
		visionone_api_key = config.VisiononeAPIKey.ValueString()
	}

	if !config.VisiononeRegion.IsNull() {
		visionone_region = config.VisiononeRegion.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if visionone_endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_endpoint"),
			"Missing VisionOne Endpoint",
			"The provider cannot create the VisionOne API client as there is a missing or empty value for the VisionOne endpoint. "+
				"Set the visionone_endpoint value in the configuration or use the VISIONONE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if visionone_endpoint_type == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_endpoint_type"),
			"Missing VisionOne Endpoint Type",
			"The provider cannot create the VisionOne API client as there is a missing or empty value for the VisionOne endpoint type. "+
				"Set the visionone_endpoint_type value in the configuration or use the VISIONONE_ENDPOINT_TYPE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if visionone_business_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_business_id"),
			"Missing VisionOne Business Id",
			"The provider cannot create the VisionOne API client as there is a missing or empty value for the VisionOne business id. "+
				"Set the visionone_business_id value in the configuration or use the VISIONONE_BUSINESS_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if visionone_api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_api_key"),
			"Missing VisionOne API Key",
			"The provider cannot create the VisionOne API client as there is a missing or empty value for the VisionOne API key. "+
				"Set the visionone_api_key value in the configuration or use the VISIONONE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if visionone_region == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("visionone_region"),
			"Missing VisionOne Region",
			"The provider cannot create the VisionOne API client as there is a missing or empty value for the VisionOne region. "+
				"Set the visionone_region value in the configuration or use the VISIONONE_REGION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "visionone_endpoint", visionone_endpoint_type)
	ctx = tflog.SetField(ctx, "visionone_endpoint_type", visionone_endpoint_type)
	ctx = tflog.SetField(ctx, "visionone_business_id", visionone_business_id)
	ctx = tflog.SetField(ctx, "visionone_api_key", visionone_api_key)
	ctx = tflog.SetField(ctx, "visionone_region", visionone_region)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "visionone_api_key")

	tflog.Debug(ctx, "Creating VisionOne API client", map[string]any{
		"visionone_endpoint":      visionone_endpoint,
		"visionone_endpoint_type": visionone_endpoint_type,
		"visionone_business_id":   visionone_business_id,
		"visionone_api_key":       visionone_api_key,
		"visionone_region":        visionone_region,
	})

	visiononeClients := &common.VisionOneClients{}
	_, err := visiononeClients.Build(visionone_endpoint, visionone_endpoint_type, visionone_business_id, visionone_api_key, visionone_region)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create VisionOne API client",
			"Unable to create VisionOne API client: "+err.Error(),
		)
		return
	}

	alicloudClients := &common.AliCloudClients{}
	_, err = alicloudClients.Build()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create AliCloud API client",
			"Unable to create AliCloud API client: "+err.Error(),
		)
		return
	}

	clients := &aliCloudSecurityProviderClients{
		visiononeClients: visiononeClients,
		alicloudClients:  alicloudClients,
	}

	// Set the client on the provider data source and resource
	// configuration so it can be used by the data sources and resources.
	resp.DataSourceData = clients
	resp.ResourceData = clients

	tflog.Info(ctx, "Configured VisionOne API client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *aliCloudSecurityProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConnectedAccountSource, // temporary data source for test
	}
}

// Resources defines the resources implemented in the provider.
func (p *aliCloudSecurityProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConnectedAccountResource,
	}
}

func (p *aliCloudSecurityProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}
