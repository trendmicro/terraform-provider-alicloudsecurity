package common

// This file contains the implementation of the VisionOne client.
// It's a http client that interacts with the VisionOne API.
// The client required the API key and region to be set.
// According to the region, the client will set different base URL.
// here is the mapping
// au: api.au.xdr.trendmicro.com (Australia)
// eu: api.eu.xdr.trendmicro.com (European Union)
// in: api.in.xdr.trendmicro.com (India)
// jp: api.xdr.trendmicro.co.jp (Japan)
// sg: api.sg.xdr.trendmicro.com (Singapore)
// uae: api.mea.xdr.trendmicro.com (United Arab Emirates)
// us: api.xdr.trendmicro.com (United States)
import (
	"fmt"
	"net/http"
)

type VisionOneClient struct {
	APIKey     string
	Region     string
	BaseURL    string
	HTTPClient *http.Client
}

var regionBaseURLMap = map[string]string{
	"au":  "https://api.au.xdr.trendmicro.com",
	"eu":  "https://api.eu.xdr.trendmicro.com",
	"in":  "https://api.in.xdr.trendmicro.com",
	"jp":  "https://api.xdr.trendmicro.co.jp",
	"sg":  "https://api.sg.xdr.trendmicro.com",
	"uae": "https://api.mea.xdr.trendmicro.com",
	"us":  "https://api.xdr.trendmicro.com",
}

// NewVisionOneClient creates a new VisionOneClient instance.
func NewVisionOneClient(apiKey, region string) (*VisionOneClient, error) {
	baseURL, exists := regionBaseURLMap[region]
	if !exists {
		return nil, fmt.Errorf("unsupported region: %s", region)
	}

	return &VisionOneClient{
		APIKey:     apiKey,
		Region:     region,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

// DoRequest performs an HTTP request to the VisionOne API.
func (c *VisionOneClient) DoRequest(endpoint string, method string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// A health check to make sure the client is working
func (c *VisionOneClient) HealthCheck() error {
	panic("implement me")
}
