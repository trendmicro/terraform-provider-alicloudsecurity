package common

import (
	"context"
	"fmt"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	sts "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AliCloudClients struct {
	Sts *sts.Client
	Ram *ram.Client
}

type AliCloudClientConfig struct {
	AccessKey       string // Access Key ID
	AccessKeySecret string // Access Key Secret
	Region          string // Region ID
}

func NewAliCloudClients() *AliCloudClients {
	return &AliCloudClients{}
}

func (a *AliCloudClients) Build() (*AliCloudClients, error) {
	if _, err := a.BuildStsClient(context.Background(), ""); err != nil {
		return nil, err
	}
	if _, err := a.BuildRamClient(context.Background(), ""); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *AliCloudClients) BuildStsClient(
	ctx context.Context, region string) (*sts.Client, error) {
	if a.Sts != nil {
		return a.Sts, nil
	}

	// Configure the shared configuration
	config, err := obtainDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain default config: %v", err)
	}
	if region != "" {
		config.RegionId = tea.String(region)
	}
	config.Endpoint = tea.String(fmt.Sprintf("sts.%s.aliyuncs.com", *config.RegionId))

	// Initialize STS client
	client, err := sts.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Check if the client is not nil and can be used
	if client == nil {
		return nil, fmt.Errorf("failed to create STS client")
	} else {
		tflog.Info(ctx, "Alicloud STS client created successfully", map[string]any{
			"config": config,
		})
		a.Sts = client
	}

	return a.Sts, nil
}

func (a *AliCloudClients) BuildRamClient(
	ctx context.Context, region string) (*ram.Client, error) {
	if a.Ram != nil {
		return a.Ram, nil
	}

	// Configure the shared configuration
	config, err := obtainDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to obtain default config: %v", err)
	}
	if region != "" {
		config.RegionId = tea.String(region)
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

	a.Ram = client
	return a.Ram, nil
}

// Obtain the default configuration for the AliCloud client from the environment variables.
func obtainDefaultConfig() (*openapi.Config, error) {
	alicloudAccessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	alicloudAccessSecret := os.Getenv("ALICLOUD_ACCESS_SECRET")
	alicloudRegion := os.Getenv("ALICLOUD_REGION")

	if alicloudAccessKey == "" || alicloudAccessSecret == "" || alicloudRegion == "" {
		return nil, fmt.Errorf("ALICLOUD_ACCESS_KEY, ALICLOUD_ACCESS_SECRET, and ALICLOUD_REGION environment variables must be set")
	}

	return &openapi.Config{
		AccessKeyId:     tea.String(alicloudAccessKey),
		AccessKeySecret: tea.String(alicloudAccessSecret),
		RegionId:        tea.String(alicloudRegion),
	}, nil
}
