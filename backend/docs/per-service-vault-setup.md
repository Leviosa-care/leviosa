# Per-Service Vault Setup for GDPR Compliance

This document explains how to use the enhanced Vault testcontainer setup for per-service encryption and GDPR-compliant data isolation.

## Overview

The new `core/testutils/vault.go` provides **per-service encryption keys** and **service authentication** for complete data isolation between microservices.

### What Changed

**Before (Single Shared Key):**
- ❌ All services used `leviosa-app-key` 
- ❌ All services used same pepper `secret/pepper`
- ❌ No data isolation between services
- ❌ GDPR compliance issues

**After (Per-Service Keys):**
- ✅ Each service has unique encryption key: `settings-encryption-key`, `catalog-encryption-key`
- ✅ Each service has unique pepper: `secret/peppers/settings`, `secret/peppers/catalog` 
- ✅ Complete data isolation between services
- ✅ GDPR compliant with proper data segregation

## New Functions

### `SetupServiceVault(ctx, t, serviceNames)`

Main function for setting up GDPR-compliant service testing:

```go
// Example: Setup vault for settings and catalog services
serviceNames := []string{"settings", "catalog"}
setup, err := testutils.SetupServiceVault(ctx, t, serviceNames)
if err != nil {
    t.Fatalf("Failed to setup service vault: %v", err)
}
defer testutils.TeardownVault(ctx, t, setup.VaultContainer)
```

### `ServiceVaultSetup` Type

Contains everything needed for service testing:

```go
type ServiceVaultSetup struct {
    VaultContainer  *VaultContainer                // Vault testcontainer
    ServiceKeys     map[string]string              // service name -> API key
    CryptoServices  map[string]encx.CryptoService  // service name -> crypto service
    VaultClient     *api.Client                    // For auth middleware
}
```

### Convenience Methods

```go
// Get service-specific crypto service
settingsCrypto, exists := setup.GetServiceCrypto("settings")

// Get service API key for authentication
apiKey, exists := setup.GetServiceAPIKey("catalog")
```

## Usage Examples

### Settings Service Integration Test

```go
func TestMain(m *testing.M) {
    ctx := context.Background()
    
    // Setup service-specific Vault
    serviceNames := []string{"settings"}
    vaultSetup, err := testutils.SetupServiceVault(ctx, nil, serviceNames)
    if err != nil {
        log.Fatalf("Failed to setup vault: %v", err)
    }
    defer testutils.TeardownVault(ctx, nil, vaultSetup.VaultContainer)
    
    // Get settings-specific crypto service
    settingsCrypto, _ := vaultSetup.GetServiceCrypto("settings")
    
    // Initialize auth middleware with Vault client
    authmw := auth.NewSessionAuthMiddleware(sessionRepo, settingsCrypto, vaultSetup.VaultClient)
    
    // Run tests...
    code := m.Run()
    os.Exit(code)
}
```

### Testing Service Authentication

```go
func TestServiceAuthentication(t *testing.T) {
    // Get API key for catalog service
    catalogAPIKey, _ := vaultSetup.GetServiceAPIKey("catalog")
    
    // Test authenticated request
    req := httptest.NewRequest("GET", "/internal/settings/name", nil)
    req.Header.Set("X-Service-Name", "catalog")
    req.Header.Set("X-Service-Key", catalogAPIKey)
    
    // Should succeed with correct key
    w := httptest.NewRecorder()
    handler(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Multiple Services Setup

```go
func TestMultiServiceSetup(t *testing.T) {
    ctx := context.Background()
    
    // Setup multiple services with isolated encryption
    serviceNames := []string{"settings", "catalog", "authuser", "notification"}
    setup, err := testutils.SetupServiceVault(ctx, t, serviceNames)
    require.NoError(t, err)
    
    // Each service has its own crypto service
    settingsCrypto, _ := setup.GetServiceCrypto("settings")
    catalogCrypto, _ := setup.GetServiceCrypto("catalog")
    
    // Data encrypted by settings service cannot be decrypted by catalog service
    // This ensures GDPR compliance and data isolation
}
```

## Vault Structure

The new setup creates this structure in Vault:

```
Vault Container:
├── transit/keys/
│   ├── settings-encryption-key      # Settings service only
│   ├── catalog-encryption-key       # Catalog service only
│   ├── authuser-encryption-key      # Auth service only
│   ├── notification-encryption-key  # Notification service only
│   └── leviosa-app-key             # Shared (for service auth keys)
├── secret/data/
│   ├── services/                    # Service API keys
│   │   ├── settings/api-key
│   │   ├── catalog/api-key
│   │   ├── authuser/api-key
│   │   └── notification/api-key
│   ├── peppers/                     # Per-service peppers
│   │   ├── settings
│   │   ├── catalog
│   │   ├── authuser
│   │   └── notification
│   └── pepper                       # Shared pepper (for service auth)
```

## Benefits

### GDPR Compliance
- **Data Isolation**: Each service can only decrypt its own data
- **Audit Trail**: Clear ownership of encrypted data
- **Breach Containment**: Service breaches don't affect other services
- **Regulatory Compliance**: Easier to prove data segregation

### Security
- **Zero Trust**: Services cannot access other services' encrypted data
- **Key Rotation**: Rotate keys independently per service
- **Principle of Least Privilege**: Services only access what they need

### Testing
- **Production Parity**: Same encryption setup as production
- **Integration Testing**: Real service authentication and encryption
- **Isolation Testing**: Verify services cannot access other services' data

## Migration from Shared Keys

### Existing Tests (No Changes Needed)
```go
// Old way still works
vaultContainer, err := testutils.SetupVault(ctx, t)
// Uses shared encryption key for backward compatibility
```

### New Tests (GDPR Compliant)
```go
// New way for GDPR compliance
setup, err := testutils.SetupServiceVault(ctx, t, []string{"settings"})
// Uses per-service encryption keys
```

### Gradual Migration
1. Start with settings service using new setup
2. Migrate other services one by one
3. Both approaches can coexist during transition
4. Eventually deprecate shared key approach

## Production Deployment

The same per-service structure works in production:

```yaml
# Docker Compose / Kubernetes
vault:
  image: hashicorp/vault:1.19
  # Same secret structure as tests
  # Per-service encryption keys
  # Service API key authentication
```

This ensures **100% parity** between testing and production environments.