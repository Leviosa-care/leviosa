# Phase 2: Infrastructure Secret Management

This document outlines the complete infrastructure secret management implementation for production deployment with HashiCorp Vault.

## Overview

Phase 2 moves all application secrets (database credentials, API keys, external service credentials) from environment variables to HashiCorp Vault, enabling centralized secret management, automatic rotation, and secure CI/CD pipelines.

## Production Vault Structure

### Per-Service Encryption Secrets
Each service gets its own encryption circle for GDPR compliance and data isolation:

```
# Service-specific peppers for hashing
secret/data/authuser/pepper
secret/data/catalog/pepper  
secret/data/settings/pepper
secret/data/notification/pepper

# Service-specific Key Encryption Keys
transit/keys/authuser-kek
transit/keys/catalog-kek
transit/keys/settings-kek
transit/keys/notification-kek

# Service-to-service API keys
secret/data/services/authuser/api-key
secret/data/services/catalog/api-key
secret/data/services/settings/api-key
secret/data/services/notification/api-key
```

### Infrastructure Secrets
Shared infrastructure credentials accessed by all services:

```
secret/data/infrastructure/postgres/credentials
secret/data/infrastructure/redis/credentials
secret/data/infrastructure/rabbitmq/credentials
secret/data/infrastructure/stripe/credentials
secret/data/infrastructure/s3/credentials
secret/data/infrastructure/gmail/credentials
secret/data/infrastructure/twilio/credentials
```

## Application Integration

### Service Startup Pattern

Each service will initialize with Vault integration:

```go
// main.go for each service
func main() {
    // 1. Initialize Vault client with authentication
    vaultClient := initVaultClient()
    
    // 2. Load service-specific configuration
    config := loadServiceConfig(vaultClient, services.Settings) // Use constants
    
    // 3. Initialize encx with service-specific KEK
    crypto := encx.NewCrypto(vaultClient, services.Settings)
    
    // 4. Initialize infrastructure dependencies
    db := postgres.Connect(config.Database)
    redis := redis.Connect(config.Redis)
    rabbitmq := rabbitmq.Connect(config.RabbitMQ)
    
    // 5. Start service with loaded configuration
    startServer(config, crypto, db, redis, rabbitmq)
}
```

### Configuration Loading

```go
type Config struct {
    // Infrastructure secrets
    Database DatabaseConfig
    Redis    RedisConfig
    RabbitMQ RabbitMQConfig
    Stripe   StripeConfig
    S3       S3Config
    Gmail    GmailConfig
    Twilio   TwilioConfig
    
    // Service-specific secrets
    ServiceAPIKey string
    
    // Service identification
    ServiceName string
    Port        int
}

func loadServiceConfig(vault *vault.Client, serviceName string) *Config {
    return &Config{
        ServiceName:   serviceName,
        Database:     loadDatabaseConfig(vault),
        Redis:        loadRedisConfig(vault),
        RabbitMQ:     loadRabbitMQConfig(vault),
        Stripe:       loadStripeConfig(vault),
        S3:           loadS3Config(vault),
        Gmail:        loadGmailConfig(vault),
        Twilio:       loadTwilioConfig(vault),
        ServiceAPIKey: loadServiceAPIKey(vault, serviceName),
    }
}
```

## Docker Compose Integration

### Vault Agent Pattern

Use Vault Agent as a sidecar container for automatic token renewal and secret injection:

```yaml
# docker-compose.yml
services:
  vault-agent:
    image: hashicorp/vault:1.19
    command: vault agent -config=/vault/config/agent.hcl
    volumes:
      - vault-config:/vault/config:ro
      - vault-secrets:/vault/secrets
    environment:
      VAULT_ADDR: "${VAULT_ADDR}"
    networks:
      - leviosa-network

  settings-service:
    build: ./settings
    environment:
      VAULT_ADDR: "${VAULT_ADDR}"
      SERVICE_NAME: "settings"
    volumes:
      - vault-secrets:/vault/secrets:ro
    depends_on:
      - vault-agent
    networks:
      - leviosa-network

  catalog-service:
    build: ./catalog
    environment:
      VAULT_ADDR: "${VAULT_ADDR}"
      SERVICE_NAME: "catalog"
    volumes:
      - vault-secrets:/vault/secrets:ro
    depends_on:
      - vault-agent
    networks:
      - leviosa-network
```

### Vault Agent Configuration

```hcl
# vault-config/agent.hcl
vault {
  address = "https://vault.example.com:8200"
}

auto_auth {
  method "approle" {
    mount_path = "auth/approle"
    config = {
      role_id_file_path = "/vault/config/role-id"
      secret_id_file_path = "/vault/config/secret-id"
    }
  }

  sink "file" {
    config = {
      path = "/vault/secrets/token"
    }
  }
}

template {
  source = "/vault/config/templates/database.tpl"
  destination = "/vault/secrets/database.json"
}

template {
  source = "/vault/config/templates/redis.tpl" 
  destination = "/vault/secrets/redis.json"
}
```

## CI/CD Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v4
      
      - name: Authenticate with Vault
        run: |
          # Use GitHub OIDC to get Vault token
          VAULT_TOKEN=$(vault write -field=token auth/jwt/login \
            role=github-actions \
            jwt=${{ github.token }})
          echo "VAULT_TOKEN=$VAULT_TOKEN" >> $GITHUB_ENV
      
      - name: Initialize Vault Secrets
        run: |
          # Create any missing secrets
          ./scripts/init-vault-secrets.sh
      
      - name: Deploy Services
        run: |
          export VAULT_ADDR=${{ secrets.VAULT_ADDR }}
          export VAULT_TOKEN=${{ env.VAULT_TOKEN }}
          docker-compose up -d --build
      
      - name: Health Checks
        run: |
          ./scripts/health-check.sh
```

### Self-Hosted Runner Setup

The self-hosted runner needs:

1. **Vault Agent**: Running as a service for automatic token renewal
2. **Docker Compose**: With access to Vault secrets volume
3. **AppRole Credentials**: Stored securely on the runner
4. **Network Access**: To your Vault instance

```bash
# /etc/systemd/system/vault-agent.service
[Unit]
Description=Vault Agent
After=network.target

[Service]
Type=simple
User=vault-agent
ExecStart=/usr/bin/vault agent -config=/etc/vault-agent/config.hcl
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Security Considerations

### Access Control

1. **Service Policies**: Each service gets minimal required permissions
2. **Infrastructure Policies**: Separate read-only policies for infrastructure secrets
3. **Time-Limited Tokens**: All tokens have TTL and auto-renewal
4. **Audit Logging**: All secret access is logged and monitored

### Secret Rotation

1. **Automatic Rotation**: Database passwords, API keys rotate automatically
2. **Zero-Downtime**: Services reload configuration on secret rotation
3. **Version Management**: Keep multiple versions during rotation period

## Migration Strategy

### Phase 2.1: Infrastructure Secrets
- Move database, Redis, RabbitMQ credentials to Vault
- Update service initialization to load from Vault
- Test with staging environment

### Phase 2.2: External Service Secrets
- Move Stripe, S3, Gmail, Twilio credentials to Vault
- Update configuration loading
- Test payment and notification flows

### Phase 2.3: CI/CD Integration
- Set up Vault authentication for GitHub Actions
- Configure self-hosted runner with Vault Agent
- Implement health checks and monitoring

### Phase 2.4: Production Deployment
- Deploy Vault cluster for high availability
- Configure backup and disaster recovery
- Set up monitoring and alerting

## Benefits

1. **Centralized Management**: All secrets in one secure location
2. **Automatic Rotation**: Reduces manual secret management overhead
3. **Audit Trail**: Complete visibility into secret access
4. **GDPR Compliance**: Per-service encryption keys enable data isolation
5. **CI/CD Security**: No secrets in environment variables or config files
6. **Scalability**: Easy to add new services with proper secret isolation

## Implementation Timeline

- **Week 1**: Phase 2.1 - Infrastructure secrets migration
- **Week 2**: Phase 2.2 - External service secrets migration  
- **Week 3**: Phase 2.3 - CI/CD integration and testing
- **Week 4**: Phase 2.4 - Production deployment and monitoring setup

This comprehensive approach ensures secure, scalable secret management while maintaining the per-service architecture and GDPR compliance requirements.
