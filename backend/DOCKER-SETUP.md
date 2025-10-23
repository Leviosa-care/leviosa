# Docker Compose Setup for Leviosa

This document explains how to run the Leviosa application using Docker Compose with full service-to-service authentication.

## Quick Start

```bash
# Start all services
docker-compose up -d

# View logs for service initialization
docker-compose logs vault-init

# Check service status
docker-compose ps

# View specific service logs
docker-compose logs settings
docker-compose logs catalog
```

## Architecture

### Services

1. **vault** - HashiCorp Vault 1.19 for secret management
2. **postgres** - PostgreSQL 17.5 database
3. **redis** - Redis cache
4. **rabbitmq** - RabbitMQ message broker
5. **localstack** - S3-compatible object storage (development)
6. **vault-init** - One-time initialization container for service keys
7. **settings** - Settings microservice (port 8080)
8. **catalog** - Catalog microservice (port 8081)

### Service Discovery

Services communicate using container names:
- `http://vault:8200` - Vault API
- `http://settings:8080` - Settings service
- `http://postgres:5432` - PostgreSQL database
- `http://redis:6379` - Redis cache

## Service Authentication

### How It Works

1. **vault-init** container generates API keys for all services
2. Each service gets a unique API key stored in Vault
3. Services authenticate using `X-Service-Name` and `X-Service-Key` headers
4. Vault validates keys against stored hashes

### API Key Generation

Service API keys are automatically generated during startup:

```bash
# View generated keys
docker-compose logs vault-init

# Output example:
# AUTHUSER_SERVICE_API_KEY=xyz789abc123...
# CATALOG_SERVICE_API_KEY=def456ghi789...
# SETTINGS_SERVICE_API_KEY=jkl012mno345...
# NOTIFICATION_SERVICE_API_KEY=pqr678stu901...
```

### Service-to-Service Communication

Example: Catalog service calling Settings service:

```go
// Catalog service code
client, err := httpx.NewServiceClient(httpx.ServiceClientConfig{
    ServiceName: services.Catalog,
    APIKey:      os.Getenv("SERVICE_API_KEY"),
    BaseURL:     "http://settings:8080",
})

// Make authenticated request
resp, err := client.Get(ctx, "/internal/settings/name")
```

## Development

### Hot Reload

Both settings and catalog services use Air for hot reload:

```bash
# Edit source files and see changes immediately
# No restart needed for Go code changes
```

### Adding New Services

1. Create new service Dockerfile
2. Add service to `docker-compose.yml`
3. Generate API key in `vault-init` script
4. Configure service authentication

## Production Deployment

### Docker Swarm

This setup works with Docker Swarm for production:

```bash
# Deploy to swarm
docker stack deploy -c docker-compose.yml leviosa

# Scale services
docker service scale leviosa_settings=3
docker service scale leviosa_catalog=2
```

### External Vault

For production, replace the dev Vault with an external cluster:

```yaml
# docker-compose.prod.yml
services:
  settings:
    environment:
      VAULT_ADDR: https://vault.company.com:8200
      VAULT_TOKEN: ${VAULT_SERVICE_TOKEN}
```

## Security

### Secret Management

- API keys are generated with 256-bit entropy
- Only key hashes are stored in Vault (never plaintext)
- Each service has unique credentials
- Automatic key rotation supported (extend `vault-init`)

### GDPR Compliance

- Per-service encryption keys
- Data isolation between services
- Audit trail via Vault logs

## Troubleshooting

### Service Won't Start

```bash
# Check dependencies
docker-compose logs vault
docker-compose logs postgres

# Verify Vault initialization
docker-compose logs vault-init

# Check service configuration
docker-compose logs settings
```

### Authentication Failures

```bash
# Verify API keys
docker-compose logs vault-init | grep "SERVICE_API_KEY"

# Check Vault secrets
docker-compose exec vault sh
vault kv list secret/services/
vault kv get secret/services/catalog/api-key
```

### Network Issues

```bash
# Test service connectivity
docker-compose exec settings ping vault
docker-compose exec catalog ping settings

# Check port bindings
docker-compose ps
```

## Monitoring

### Health Checks

All services include health checks:

```bash
# View service health
docker-compose ps

# Check specific service
docker-compose exec settings wget -O- http://localhost:8080/health
```

### Logs

```bash
# Follow all logs
docker-compose logs -f

# Service-specific logs
docker-compose logs -f settings catalog

# Filter by timestamp
docker-compose logs --since="1h" vault-init
```

## Testing

### Integration Tests

Testcontainers use the same setup as Docker Compose:

```bash
# Run integration tests (uses same Vault setup)
cd settings
make test-integration

# Tests use identical:
# - Vault 1.19 container
# - Service authentication flow
# - Secret structure and paths
```

### Manual Testing

```bash
# Test service endpoints directly
curl http://localhost:8080/settings/name

# Test service-to-service auth (requires API key)
curl -H "X-Service-Name: catalog" \
     -H "X-Service-Key: YOUR_API_KEY" \
     http://localhost:8080/internal/settings/name
```

## Cleanup

```bash
# Stop all services
docker-compose down

# Remove volumes (data loss!)
docker-compose down -v

# Remove images
docker-compose down --rmi all
```