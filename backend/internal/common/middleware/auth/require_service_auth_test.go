package auth

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/stretchr/testify/assert"
)

func TestServiceConstants(t *testing.T) {
	// Test that service header constants are correctly defined
	assert.Equal(t, "X-Service-Key", services.ServiceKeyHeader)
	assert.Equal(t, "X-Service-Name", services.ServiceNameHeader)
}

func TestGetServiceInfoFromContext(t *testing.T) {
	tests := []struct {
		name        string
		contextValue interface{}
		expectError  bool
		expectedName string
	}{
		{
			name:         "Valid service info",
			contextValue: &ServiceInfo{Name: services.Catalog},
			expectError:  false,
			expectedName: services.Catalog,
		},
		{
			name:         "Missing service info",
			contextValue: nil,
			expectError:  true,
			expectedName: "",
		},
		{
			name:         "Wrong type in context",
			contextValue: "not-service-info",
			expectError:  true,
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.contextValue != nil {
				ctx = context.WithValue(ctx, ServiceContextKey, tt.contextValue)
			}

			serviceInfo, err := GetServiceInfoFromContext(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, serviceInfo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, serviceInfo)
				assert.Equal(t, tt.expectedName, serviceInfo.Name)
			}
		})
	}
}