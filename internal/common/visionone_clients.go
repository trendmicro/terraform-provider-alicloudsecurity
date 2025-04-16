package common

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type VisionOneClients struct {
	Cam *CamClient
}

func (v *VisionOneClients) Build(apiKey, region string) (*VisionOneClients, error) {
	_, err := v.BuildCamClient(apiKey, region)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *VisionOneClients) BuildCamClient(apiKey, region string) (*CamClient, error) {
	if v.Cam != nil {
		return v.Cam, nil
	}

	if len(apiKey) == 0 {
		apiKey = os.Getenv("VISIONONE_API_KEY")
	}
	if len(region) == 0 {
		region = os.Getenv("VISIONONE_REGION")
	}

	client, err := NewCamClient(apiKey, region)
	if err != nil {
		return nil, err
	}

	tflog.Info(context.Background(), "VisionOne client created successfully", map[string]any{
		"apiKey": apiKey,
		"region": region,
	})

	v.Cam = client
	return v.Cam, nil
}
