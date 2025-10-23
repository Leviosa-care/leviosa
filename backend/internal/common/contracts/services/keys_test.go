package services

import (
	"testing"
)

func TestServicePepperPath(t *testing.T) {
	tests := []struct {
		service  string
		expected string
	}{
		{AuthUser, "secret/data/authuser/pepper"},
		{Catalog, "secret/data/catalog/pepper"},
		{Settings, "secret/data/settings/pepper"},
		{Notification, "secret/data/notification/pepper"},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			result := ServicePepperPath(tt.service)
			if result != tt.expected {
				t.Errorf("ServicePepperPath(%s) = %s, expected %s", tt.service, result, tt.expected)
			}
		})
	}
}

func TestServiceKEKPath(t *testing.T) {
	tests := []struct {
		service  string
		expected string
	}{
		{AuthUser, "transit/keys/authuser-kek"},
		{Catalog, "transit/keys/catalog-kek"},
		{Settings, "transit/keys/settings-kek"},
		{Notification, "transit/keys/notification-kek"},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			result := ServiceKEKPath(tt.service)
			if result != tt.expected {
				t.Errorf("ServiceKEKPath(%s) = %s, expected %s", tt.service, result, tt.expected)
			}
		})
	}
}

func TestServiceAPIKeyPath(t *testing.T) {
	tests := []struct {
		service  string
		expected string
	}{
		{AuthUser, "secret/data/services/authuser/api-key"},
		{Catalog, "secret/data/services/catalog/api-key"},
		{Settings, "secret/data/services/settings/api-key"},
		{Notification, "secret/data/services/notification/api-key"},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			result := ServiceAPIKeyPath(tt.service)
			if result != tt.expected {
				t.Errorf("ServiceAPIKeyPath(%s) = %s, expected %s", tt.service, result, tt.expected)
			}
		})
	}
}

func TestServiceVaultPaths(t *testing.T) {
	paths := ServiceVaultPaths(Settings)

	expected := map[string]string{
		"pepper":  "secret/data/settings/pepper",
		"kek":     "transit/keys/settings-kek",
		"api_key": "secret/data/services/settings/api-key",
	}

	for key, expectedPath := range expected {
		if paths[key] != expectedPath {
			t.Errorf("ServiceVaultPaths()[%s] = %s, expected %s", key, paths[key], expectedPath)
		}
	}
}

