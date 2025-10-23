package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
)

// SettingsClient provides a high-level interface to the settings service
type SettingsClient struct {
	serviceClient *httpx.ServiceClient
}

// SettingsResponse represents the response structure from settings endpoints
type SettingsResponse struct {
	Value string `json:"value"`
}

// NewSettingsClient creates a new settings service client
func NewSettingsClient(serviceName, apiKey, settingsURL string) (*SettingsClient, error) {
	config := httpx.ServiceClientConfig{
		ServiceName: serviceName,
		APIKey:      apiKey,
		BaseURL:     settingsURL,
	}
	
	serviceClient, err := httpx.NewServiceClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create service client: %w", err)
	}
	
	return &SettingsClient{
		serviceClient: serviceClient,
	}, nil
}

// GetCompanyName retrieves the company name from settings service
func (sc *SettingsClient) GetCompanyName(ctx context.Context) (string, error) {
	resp, err := sc.serviceClient.Get(ctx, "/internal/settings/name")
	if err != nil {
		return "", fmt.Errorf("failed to get company name: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("settings service returned status %d", resp.StatusCode)
	}
	
	var settingsResp SettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&settingsResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	return settingsResp.Value, nil
}

// GetCompanyEmail retrieves the company email from settings service
func (sc *SettingsClient) GetCompanyEmail(ctx context.Context) (string, error) {
	resp, err := sc.serviceClient.Get(ctx, "/internal/settings/email")
	if err != nil {
		return "", fmt.Errorf("failed to get company email: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("settings service returned status %d", resp.StatusCode)
	}
	
	var settingsResp SettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&settingsResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	return settingsResp.Value, nil
}

// GetOTPSettings retrieves OTP configuration from settings service
func (sc *SettingsClient) GetOTPSettings(ctx context.Context) (map[string]string, error) {
	otpSettings := make(map[string]string)
	
	// Get OTP duration
	resp, err := sc.serviceClient.Get(ctx, "/internal/settings/otp/duration")
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP duration: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		var settingsResp SettingsResponse
		if err := json.NewDecoder(resp.Body).Decode(&settingsResp); err == nil {
			otpSettings["duration"] = settingsResp.Value
		}
	}
	
	// Get OTP length
	resp, err = sc.serviceClient.Get(ctx, "/internal/settings/otp/length")
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP length: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		var settingsResp SettingsResponse
		if err := json.NewDecoder(resp.Body).Decode(&settingsResp); err == nil {
			otpSettings["length"] = settingsResp.Value
		}
	}
	
	// Get OTP max attempts
	resp, err = sc.serviceClient.Get(ctx, "/internal/settings/otp/max-attempts")
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP max attempts: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		var settingsResp SettingsResponse
		if err := json.NewDecoder(resp.Body).Decode(&settingsResp); err == nil {
			otpSettings["max_attempts"] = settingsResp.Value
		}
	}
	
	return otpSettings, nil
}

// Example usage demonstrating service-to-service communication
func main() {
	// Get service configuration from environment
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = services.Catalog // Default to catalog for this example
	}
	
	apiKey := os.Getenv("SERVICE_API_KEY")
	if apiKey == "" {
		log.Fatal("SERVICE_API_KEY environment variable is required")
	}
	
	settingsURL := os.Getenv("SETTINGS_SERVICE_URL")
	if settingsURL == "" {
		settingsURL = "http://localhost:8080" // Default settings service URL
	}
	
	// Create settings client
	client, err := NewSettingsClient(serviceName, apiKey, settingsURL)
	if err != nil {
		log.Fatalf("Failed to create settings client: %v", err)
	}
	
	ctx := context.Background()
	
	// Example 1: Get company information
	fmt.Println("=== Company Information ===")
	companyName, err := client.GetCompanyName(ctx)
	if err != nil {
		log.Printf("Failed to get company name: %v", err)
	} else {
		fmt.Printf("Company Name: %s\n", companyName)
	}
	
	companyEmail, err := client.GetCompanyEmail(ctx)
	if err != nil {
		log.Printf("Failed to get company email: %v", err)
	} else {
		fmt.Printf("Company Email: %s\n", companyEmail)
	}
	
	// Example 2: Get OTP settings (useful for authuser service)
	if serviceName == services.AuthUser {
		fmt.Println("\n=== OTP Configuration ===")
		otpSettings, err := client.GetOTPSettings(ctx)
		if err != nil {
			log.Printf("Failed to get OTP settings: %v", err)
		} else {
			for key, value := range otpSettings {
				fmt.Printf("OTP %s: %s\n", key, value)
			}
		}
	}
	
	fmt.Println("\n=== Service-to-Service Communication Test Complete ===")
	fmt.Printf("Service '%s' successfully communicated with settings service\n", serviceName)
}
