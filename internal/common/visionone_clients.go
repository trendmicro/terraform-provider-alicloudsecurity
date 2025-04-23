package common

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type VisionOneClients struct {
	Cam *CamClient
}

func (v *VisionOneClients) Build(endpoint, endpointType, businessId, apiKey, region string) (*VisionOneClients, error) {
	_, err := v.BuildCamClient(endpoint, endpointType, businessId, apiKey, region)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *VisionOneClients) BuildCamClient(endpoint, endpointType, businessId, apiKey, region string) (*CamClient, error) {
	if v.Cam != nil {
		return v.Cam, nil
	}

	config := &CamClientConfig{
		Endpoint:     &endpoint,
		EndpointType: &endpointType,
		BusinessId:   &businessId,
		ApiKey:       &apiKey,
		Region:       &region,
	}
	client, err := NewCamClient(config)
	if err != nil {
		return nil, err
	}

	tflog.Info(context.Background(), "CAM client created successfully", map[string]any{
		"businessId": businessId,
		"apiKey":     apiKey,
		"region":     region,
	})

	v.Cam = client
	return v.Cam, nil
}
