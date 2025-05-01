package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CamClient struct {
	Config *CamClientConfig
	Client *http.Client
}

type CamClientConfig struct {
	Endpoint     *string
	EndpointType *string
	Region       *string
	ApiKey       *string
	BusinessId   *string
}

type CreateConnectionRequest struct {
	AccountId      *string `json:"accountId" validate:"max=16"`
	Region         *string `json:"region" validate:"max=254"`
	RoleArn        *string `json:"roleArn" validate:"max=254"`
	OidcProviderId *string `json:"oidcProviderId" validate:"max=254"`
	Name           *string `json:"name" validate:"max=254"`
	Description    *string `json:"description" validate:"omitempty,max=254"`
}

type UpdateConnectionRequest struct {
	Name        *string `json:"name"`        // The name of the Alibaba Cloud account to be used in Cloud Account Management.
	Description *string `json:"description"` // The description of the Alibaba Cloud account. The default value is an empty string if the field is omitted.
}

type ReadConnectionResponse struct {
	Id                 *string `json:"id"`                 // The ID of the Alibaba Cloud account.
	ParentStackRegion  *string `json:"parentStackRegion"`  // The region of Terraform backend where the state files are stored.
	RoleArn            *string `json:"roleArn"`            // The Alibaba Cloud resource name (ARN) of the user role for Trend Vision One.
	OidcProviderId     *string `json:"oidcProviderId"`     // The ID of the Alibaba Cloud OpenID Connect (OIDC) provider.
	Name               *string `json:"name"`               // The name of the Alibaba Cloud account used in Cloud Account Management.
	Description        *string `json:"description"`        // The description of the Alibaba Cloud account.
	CreatedDateTime    *string `json:"createdDateTime"`    // The timestamp indicating when the Alibaba Cloud account was added to Trend Vision One.
	UpdatedDateTime    *string `json:"updatedDateTime"`    // The timestamp indicating the last time the Alibaba Cloud account was modified.
	State              *string `json:"state"`              // The status of the Alibaba Cloud account.
	LastSyncedDateTime *string `json:"lastSyncedDateTime"` // The timestamp indicating the most recent synchronization of the Alibaba Cloud account with the cloud provider.
}

var apiPathMap = map[string]map[string]string{
	"create": {
		"automation": "/v3.0/cam/alibabaAccounts",          // The endpoint would be like: https://api-int.visionone.trendmicro.com
		"express":    "/public/cam/api/ui/alibabaAccounts", // The endpoint would be like: https://cloudaccounts-dev-us.visionone.trendmicro.com
	},
	"update": {
		"automation": "/v3.0/cam/alibabaAccounts/%s",
		"express":    "/public/cam/api/ui/alibabaAccounts/%s",
	},
	"delete": {
		"automation": "/v3.0/cam/alibabaAccounts/%s",
		"express":    "/public/cam/api/ui/alibabaAccounts/%s",
	},
	"read": {
		"automation": "/v3.0/cam/alibabaAccounts/%s",
		"express":    "/public/cam/api/ui/alibabaAccounts/%s",
	},
}

// NewCamClient creates a new CamClient instance.
func NewCamClient(config *CamClientConfig) (*CamClient, error) {
	// Ensure the config is not nil
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Validate that all required fields in the config are not nil or empty
	if config.Endpoint == nil || *config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be nil or empty")
	}
	if config.ApiKey == nil || *config.ApiKey == "" {
		return nil, fmt.Errorf("API key cannot be nil or empty")
	}
	if config.Region == nil || *config.Region == "" {
		return nil, fmt.Errorf("region cannot be nil or empty")
	}

	// Create an HTTP client with the Authorization header
	client := &http.Client{}

	// Create and return the CamClient instance
	return &CamClient{
		Config: config,
		Client: client,
	}, nil
}

// DoRequest performs an HTTP request to the VisionOne API.
func (c *CamClient) DoRequest(ctx context.Context, method, url string, body []byte) (*http.Response, error) {
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *c.Config.ApiKey))
	req.Header.Set("x-customer-id", *c.Config.BusinessId)
	req.Header.Set("x-task-id", GenerateUUID())
	req.Header.Set("x-trace-id", GenerateUUID())
	req.Header.Set("Content-Type", "application/json")

	return c.Client.Do(req)
}

func (c *CamClient) CreateConnection(ctx context.Context, req *CreateConnectionRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	urlPattern, err := BuildUrlPattern(c.Config, "create")
	if err != nil {
		return err
	}
	url := urlPattern

	tflog.Debug(ctx, fmt.Sprintf("CreateConnection URL: %s", url))
	tflog.Debug(ctx, fmt.Sprintf("CreateConnection Request: %s", string(body)))

	resp, err := c.DoRequest(ctx, "POST", url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		// Read the response body to get more details about the error
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}

		var respBody map[string]interface{}
		if err := json.Unmarshal(respBodyBytes, &respBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %v", err)
		} else {
			return fmt.Errorf("failed to create connection: status code %d, response: %s", resp.StatusCode, string(respBodyBytes))
		}
	}

	return nil
}

func (c *CamClient) UpdateConnection(ctx context.Context, accountId *string, req *UpdateConnectionRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	urlPattern, err := BuildUrlPattern(c.Config, "update")
	if err != nil {
		return err
	}
	url := fmt.Sprintf(urlPattern, *accountId)

	tflog.Debug(ctx, fmt.Sprintf("UpdateConnection URL: %s", url))
	tflog.Debug(ctx, fmt.Sprintf("UpdateConnection Request: %s", string(body)))

	resp, err := c.DoRequest(ctx, http.MethodPatch, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		// Read the response body to get more details about the error
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		var respBody map[string]interface{}
		if err := json.Unmarshal(respBodyBytes, &respBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %v", err)
		} else {
			return fmt.Errorf("failed to update connection: status code %d, response: %s", resp.StatusCode, string(respBodyBytes))
		}
	}
	return nil
}

func (c *CamClient) DeleteConnection(ctx context.Context, accountId *string) error {
	if len(*accountId) == 0 {
		return fmt.Errorf("account id cannot be empty")
	}

	urlPattern, err := BuildUrlPattern(c.Config, "delete")
	if err != nil {
		return err
	}
	url := fmt.Sprintf(urlPattern, *accountId)

	tflog.Debug(ctx, fmt.Sprintf("DeleteConnection URL: %s", url))
	tflog.Debug(ctx, fmt.Sprintf("DeleteConnection Account ID: %s", *accountId))

	resp, err := c.DoRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		// Read the response body to get more details about the error
		respBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		var respBody map[string]interface{}
		if err := json.Unmarshal(respBodyBytes, &respBody); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %v", err)
		} else {
			return fmt.Errorf("failed to delete connection: status code %d, response: %s", resp.StatusCode, string(respBodyBytes))
		}
	}
	return nil
}

func (c *CamClient) ReadConnection(ctx context.Context, accountId *string) (*ReadConnectionResponse, error) {
	if len(*accountId) == 0 {
		return nil, fmt.Errorf("account id cannot be empty")
	}

	urlPattern, err := BuildUrlPattern(c.Config, "read")
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(urlPattern, *accountId)

	tflog.Debug(ctx, fmt.Sprintf("ReadConnection URL: %s", url))
	tflog.Debug(ctx, fmt.Sprintf("ReadConnection Account ID: %s", *accountId))

	resp, err := c.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusNotFound {
			// Handle the case where the account is not found, return empty response or nil
			tflog.Debug(ctx, fmt.Sprintf("ReadConnection: account %s not found", *accountId))
			return nil, nil
		} else {
			// Read the response body to get more details about the error
			respBodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}
			var respBody map[string]interface{}
			if err := json.Unmarshal(respBodyBytes, &respBody); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
			} else {
				return nil, fmt.Errorf("failed to read connection: status code %d, response: %s", resp.StatusCode, string(respBodyBytes))
			}
		}
	}

	response := &ReadConnectionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("response is nil")
	} else {
		if response.Id == nil {
			response.Id = new(string)
			*response.Id = ""
		}
		if response.ParentStackRegion == nil {
			response.ParentStackRegion = new(string)
			*response.ParentStackRegion = ""
		}
		if response.RoleArn == nil {
			response.RoleArn = new(string)
			*response.RoleArn = ""
		}
		if response.OidcProviderId == nil {
			response.OidcProviderId = new(string)
			*response.OidcProviderId = ""
		}
		if response.Name == nil {
			response.Name = new(string)
			*response.Name = ""
		}
		if response.Description == nil {
			response.Description = new(string)
			*response.Description = ""
		}
		if response.CreatedDateTime == nil {
			response.CreatedDateTime = new(string)
			*response.CreatedDateTime = ""
		}
		if response.UpdatedDateTime == nil {
			response.UpdatedDateTime = new(string)
			*response.UpdatedDateTime = ""
		}
		if response.State == nil {
			response.State = new(string)
			*response.State = ""
		}
		if response.LastSyncedDateTime == nil {
			response.LastSyncedDateTime = new(string)
			*response.LastSyncedDateTime = ""
		}
	}
	return response, nil
}

func BuildUrlPattern(c *CamClientConfig, method string) (string, error) {
	pattern := ""

	if len(*c.EndpointType) == 0 {
		return "", fmt.Errorf("endpoint type cannot be missing or empty")
	}

	if m, ok := apiPathMap[method]; ok {
		if path, ok := m[*c.EndpointType]; ok && path != "" {
			pattern = fmt.Sprintf("%s%s", *c.Endpoint, path)
		} else {
			return "", fmt.Errorf("path is missing or empty. method: %s, endpoint type: %s", method, *c.EndpointType)
		}
	} else {
		return "", fmt.Errorf("unrecognized endpoint type")
	}

	return pattern, nil
}

// GenerateUUID generates a new UUID string.
func GenerateUUID() string {
	return uuid.New().String()
}
