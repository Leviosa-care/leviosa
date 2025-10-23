package services

import (
	"testing"
)

func TestAllServices(t *testing.T) {
	services := AllServices()
	
	expected := []string{AuthUser, Catalog, Settings, Notification}
	
	if len(services) != len(expected) {
		t.Errorf("Expected %d services, got %d", len(expected), len(services))
	}
	
	for i, service := range services {
		if service != expected[i] {
			t.Errorf("Expected service %s at index %d, got %s", expected[i], i, service)
		}
	}
}

func TestIsValidService(t *testing.T) {
	tests := []struct {
		name     string
		service  string
		expected bool
	}{
		{"Valid AuthUser", AuthUser, true},
		{"Valid Catalog", Catalog, true},
		{"Valid Settings", Settings, true},
		{"Valid Notification", Notification, true},
		{"Invalid service", "invalid", false},
		{"Empty string", "", false},
		{"Case sensitive", "AuthUser", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidService(tt.service)
			if result != tt.expected {
				t.Errorf("IsValidService(%s) = %v, expected %v", tt.service, result, tt.expected)
			}
		})
	}
}

func TestServiceConstants(t *testing.T) {
	// Test that constants have expected values
	if AuthUser != "authuser" {
		t.Errorf("AuthUser constant should be 'authuser', got '%s'", AuthUser)
	}
	
	if Catalog != "catalog" {
		t.Errorf("Catalog constant should be 'catalog', got '%s'", Catalog)
	}
	
	if Settings != "settings" {
		t.Errorf("Settings constant should be 'settings', got '%s'", Settings)
	}
	
	if Notification != "notification" {
		t.Errorf("Notification constant should be 'notification', got '%s'", Notification)
	}
}