# Service Authentication Implementation Guide

This guide covers the complete implementation of Vault-backed service-to-service authentication for the Leviosa microservices architecture.

## Overview

The service authentication system provides secure communication between microservices using:

- **API Key Authentication**: Each service has a unique API key stored in HashiCorp Vault
- **Header-Based Authentication**: Services authenticate using `X-Service-Name` and `X-Service-Key` headers
- **Vault Integration**: All service keys are securely stored and validated against Vault
- **Per-Service Encryption**: Each service has its own encryption keys and peppers for GDPR compliance

## Architecture Components

### 1. Service Constants (`core/contracts/services/`)

**Files:**
- `names.go` - Service name constants and validation
- `headers.go` - HTTP header constants
- `keys.go` - Vault path generators
- `key_management.go` - Service key management utilities

**Key Functions:**
```go
services.AllServices()                    // List all service names
services.IsValidService("catalog")        // Validate service name
services.ServiceAPIKeyPath("catalog")     // Get Vault path for service key
```

### 2. Authentication Middleware (`core/middleware/auth/`)

**Updated Files:**
- `interface.go` - Added `RequireServiceAuth` method
- `session_auth_middleware.go` - Updated constructor to accept Vault client
- `require_service_auth.go` - New service authentication middleware

**Usage:**
```go
// Create middleware with Vault client
middleware := auth.NewSessionAuthMiddleware(sessionRepo, crypto, vaultClient)

// Protect endpoints
router.HandleFunc("GET /internal/settings/name", middleware.RequireServiceAuth(handler))
```

### 3. Service HTTP Client (`core/httpx/service_client.go`)

**High-level client for service-to-service communication:**
```go
client, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
    ServiceName: services.Catalog,
    APIKey:      os.Getenv("SERVICE_API_KEY"),
    BaseURL:     "http://settings:8080",
})

resp, err := client.Get(ctx, "/internal/settings/name")
```

### 4. Key Management Tools

**Scripts:**
- `scripts/init-service-keys.go` - Generate and store service keys
- `scripts/init-service-keys.sh` - Shell wrapper for key initialization

## Implementation Steps

### Step 1: Initialize Service Keys in Vault

```bash
# Set Vault configuration
export VAULT_ADDR="http://localhost:8200"
export VAULT_TOKEN="your-vault-token"

# Run key initialization script
./scripts/init-service-keys.sh
```

This generates unique API keys for all services and stores them in Vault at:
- `secret/data/services/authuser/api-key`
- `secret/data/services/catalog/api-key`
- `secret/data/services/settings/api-key`
- `secret/data/services/notification/api-key`

### Step 2: Update Service Configurations

Add the generated API keys to each service's configuration:

```bash
# Environment variables for each service
export SERVICE_NAME="catalog"
export SERVICE_API_KEY="generated-key-from-step-1"
export VAULT_ADDR="http://localhost:8200"
export VAULT_TOKEN="service-vault-token"
```

### Step 3: Update Service Initialization

Update each service's main.go to use the new middleware constructor:

```go
// Before
middleware := auth.NewSessionAuthMiddleware(sessionRepo, crypto)

// After  
vaultClient := initVaultClient() // Initialize Vault client
middleware := auth.NewSessionAuthMiddleware(sessionRepo, crypto, vaultClient)
```

### Step 4: Add Internal Routes

Services that need to be called by other services should add internal routes:

```go
// In routes.go
RequireService := h.authmw.RequireServiceAuth

// Add internal endpoints
router.HandleFunc("GET /internal/settings/name", RequireService(mw.EnableCORS(h.GetCompanyName)))
router.HandleFunc("GET /internal/settings/email", RequireService(mw.EnableCORS(h.GetCompanyEmail)))
```

### Step 5: Update Calling Services

Services that need to call other services should use the service client:

```go
// Initialize service client
settingsClient, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
    ServiceName: services.Catalog,
    APIKey:      os.Getenv("SERVICE_API_KEY"),
    BaseURL:     os.Getenv("SETTINGS_SERVICE_URL"),
})

// Make authenticated requests
resp, err := settingsClient.Get(ctx, "/internal/settings/name")
```

## Security Features

### 1. API Key Security
- **Secure Generation**: Cryptographically secure 256-bit keys
- **Hash Storage**: Only key hashes stored in Vault, never plaintext
- **Per-Service Keys**: Each service has unique credentials

### 2. Vault Integration
- **Centralized Storage**: All keys stored in HashiCorp Vault
- **Access Control**: Vault policies control key access
- **Audit Logging**: All key access is logged by Vault

### 3. Request Validation
- **Service Name Validation**: Only valid service names accepted
- **Key Verification**: Keys validated against Vault storage
- **Header Authentication**: Uses standard HTTP headers

### 4. GDPR Compliance
- **Per-Service Encryption**: Each service has separate encryption keys
- **Data Isolation**: Services can only decrypt their own data
- **Key Rotation**: Individual service keys can be rotated independently

## Testing

### Unit Tests
```bash
# Test service constants and key management
go test ./core/contracts/services/ -v

# Test middleware (requires mocks for full test)
go test ./core/middleware/auth/ -run TestServiceConstants -v
```

### Integration Tests
```bash
# Run full Vault integration tests (requires Vault container)
go test ./core/middleware/auth/ -run TestServiceAuthWithVault -v
```

### Manual Testing
```bash
# Test service-to-service communication
go run examples/service-integration/settings_client.go
```

## Production Deployment

### 1. Vault Configuration
```hcl
# Enable KV v2 secrets engine
path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Service-specific policies
path "secret/data/services/catalog/*" {
  capabilities = ["read"]
}
```

### 2. Service Configuration
```yaml
# docker-compose.yml
services:
  catalog:
    environment:
      - SERVICE_NAME=catalog
      - SERVICE_API_KEY=${CATALOG_API_KEY}
      - VAULT_ADDR=${VAULT_ADDR}
      - VAULT_TOKEN=${CATALOG_VAULT_TOKEN}
```

### 3. Key Rotation
```bash
# Rotate a service key
vault kv put secret/services/catalog/api-key key_hash="new-hash"

# Restart the service with new key
docker-compose restart catalog
```

## Error Handling

The middleware provides detailed error responses:

- **401 Unauthorized**: Missing or invalid service key
- **400 Bad Request**: Invalid service name
- **500 Internal Server Error**: Vault connection issues

All errors are logged with structured logging for debugging.

## Monitoring and Alerts

### Key Metrics to Monitor
- Service authentication success/failure rates
- Vault connection health
- API key usage patterns
- Service-to-service communication latency

### Recommended Alerts
- High authentication failure rate
- Vault service unavailable
- Unusual service key usage patterns
- Service communication timeouts

## Next Steps

1. **Production Deployment**: Deploy with proper Vault cluster
2. **Key Rotation**: Implement automatic key rotation
3. **Monitoring**: Add metrics and alerting
4. **Service Mesh**: Consider integration with service mesh for enhanced security
5. **Rate Limiting**: Add rate limiting for service-to-service calls

This implementation provides a robust foundation for secure microservice communication while maintaining the flexibility to evolve with your architecture needs.