package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
)

// ServiceClient provides HTTP client functionality for service-to-service communication
type ServiceClient struct {
	httpClient  *http.Client
	serviceName string
	apiKey      string
	baseURL     string
}

// ServiceClientConfig holds configuration for service-to-service HTTP communication
type ServiceClientConfig struct {
	ServiceName string
	APIKey      string
	BaseURL     string
	Timeout     time.Duration
}

// NewServiceClient creates a new HTTP client for service-to-service communication
func NewServiceClient(config ServiceClientConfig) (*ServiceClient, error) {
	if config.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	
	if !services.IsValidService(config.ServiceName) {
		return nil, fmt.Errorf("invalid service name: %s", config.ServiceName)
	}
	
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for service %s", config.ServiceName)
	}
	
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second // Default timeout
	}
	
	return &ServiceClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		serviceName: config.ServiceName,
		apiKey:      config.APIKey,
		baseURL:     config.BaseURL,
	}, nil
}

// Do executes an HTTP request with service authentication headers
func (sc *ServiceClient) Do(req *http.Request) (*http.Response, error) {
	// Add service authentication headers
	req.Header.Set(services.ServiceNameHeader, sc.serviceName)
	req.Header.Set(services.ServiceKeyHeader, sc.apiKey)
	
	// Add standard headers for service-to-service communication
	req.Header.Set("User-Agent", fmt.Sprintf("service-%s/1.0", sc.serviceName))
	req.Header.Set("Accept", "application/json")
	
	return sc.httpClient.Do(req)
}

// NewRequest creates a new HTTP request with the service's base URL
func (sc *ServiceClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", sc.baseURL, path)
	return NewJSONRequest(method, url, body)
}

// Get performs a GET request to the specified path
func (sc *ServiceClient) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := sc.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	
	req = req.WithContext(ctx)
	return sc.Do(req)
}

// Post performs a POST request to the specified path with JSON body
func (sc *ServiceClient) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	req, err := sc.NewRequest("POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	
	req = req.WithContext(ctx)
	return sc.Do(req)
}

// Put performs a PUT request to the specified path with JSON body
func (sc *ServiceClient) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	req, err := sc.NewRequest("PUT", path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	
	req = req.WithContext(ctx)
	return sc.Do(req)
}

// Delete performs a DELETE request to the specified path
func (sc *ServiceClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := sc.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}
	
	req = req.WithContext(ctx)
	return sc.Do(req)
}

// GetServiceName returns the service name this client is configured for
func (sc *ServiceClient) GetServiceName() string {
	return sc.serviceName
}

// GetBaseURL returns the base URL this client is configured for
func (sc *ServiceClient) GetBaseURL() string {
	return sc.baseURL
}

// NewJSONRequest creates an HTTP request with JSON body encoding
func NewJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody []byte
	var err error
	
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}
	
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req, nil
}
