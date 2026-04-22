# Vault Production Setup Guide

This document outlines the steps required to migrate HashiCorp Vault from development mode to production mode.

## Current State (Development Mode)

- Vault runs in dev mode with a single root token
- Auto-unseals on container restart
- Uses insecure token: `dev-only-token-insecure-change-for-production`
- Data is ephemeral (lost on container restart unless using volumes)

## Production Requirements

### 1. Generate Secure Credentials

Generate a secure root token (32+ characters):

```bash
# Generate secure token
openssl rand -base64 32
```

### 2. Update Group Variables

**File:** `infra/ansible/group_vars/leviosa_production.yml`

```yaml
# HashiCorp Vault Configuration
vault_addr: http://vault:8200
vault_token: <YOUR_SECURE_32_CHAR_TOKEN_HERE>
```

**File:** `infra/ansible/group_vars/leviosa_staging.yml`

```yaml
# HashiCorp Vault Configuration (connects to prod's vault container)
vault_addr: http://leviosa_vault:8200
vault_token: <YOUR_SECURE_32_CHAR_TOKEN_HERE>
```

### 3. Update Docker Compose Template

**File:** `infra/ansible/roles/app/templates/docker-compose.yml.j2`

Replace the dev mode Vault service with production configuration:

```yaml
  # HashiCorp Vault (secrets management)
  vault:
    image: hashicorp/vault:1.15
    container_name: {{ app_name }}_vault
    restart: unless-stopped
    cap_add:
      - IPC_LOCK
    environment:
      VAULT_ADDR: http://0.0.0.0:8200
      VAULT_API_ADDR: http://vault:8200
      # Remove VAULT_DEV_ROOT_TOKEN_ID for production
    ports:
      - "8200:8200"
    volumes:
      - {{ app_data_dir }}/vault:/vault/data
      - {{ app_data_dir }}/vault/config:/vault/config
    networks:
      - {{ app_name }}_network
    healthcheck:
      test: ["CMD", "vault", "status"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
    command: server
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

**Note:** Use a specific version (e.g., `1.15`) instead of `:latest` to avoid unexpected breaking changes. The `disable_mlock = true` in the config file (see below) prevents CAP_SETFCAP errors.

### 4. Create Vault Configuration File

Create a config template at `infra/ansible/roles/app/templates/vault.hcl.j2`:

```hcl
storage "file" {
  path = "/vault/data"
}

listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = 1
}

api_addr = "http://vault:8200"
cluster_addr = "https://vault:8201"
disable_mlock = true

# Maximum lease duration
default_lease_ttl = "720h"
max_lease_ttl = "720h"
```

### 5. Deployment Steps

#### Step 1: Deploy Updated Infrastructure

```bash
make ansible-deploy
```

#### Step 2: Initialize Vault

SSH into the VPS and initialize Vault:

```bash
# SSH to VPS
make prod-ssh

# Exec into Vault container
docker exec -it leviosa_vault sh

# Initialize Vault (save the output carefully!)
vault operator init -key-shares=5 -key-threshold=3
```

**IMPORTANT:** Save the unseal keys and root token securely. You cannot recover them if lost!

#### Step 3: Unseal Vault

Unseal Vault with 3 of the 5 unseal keys:

```bash
# Run this 3 times with different unseal keys
vault operator unseal <UNSEAL_KEY_1>
vault operator unseal <UNSEAL_KEY_2>
vault operator unseal <UNSEAL_KEY_3>
```

#### Step 4: Create Secret Paths

Create separate paths for production and staging:

```bash
# Login with root token
vault login <ROOT_TOKEN>

# Enable KV secrets engine
vault secrets enable -path=leviosa_production kv-v2
vault secrets enable -path=leviosa_staging kv-v2

# Create production secrets (example)
vault kv put leviosa_production/config \
  encryption_key="your-production-encryption-key" \
  database_encryption_key="your-db-encryption-key"

# Create staging secrets
vault kv put leviosa_staging/config \
  encryption_key="your-staging-encryption-key" \
  database_encryption_key="your-staging-db-encryption-key"
```

#### Step 5: Create Policy Files

Create policies at `infra/ansible/roles/app/templates/vault-policy/`:

**leviosa-production-policy.hcl:**
```hcl
path "leviosa_production/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

**leviosa-staging-policy.hcl:**
```hcl
path "leviosa_staging/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

Apply policies:

```bash
vault policy write leviosa-production - <<EOF
path "leviosa_production/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
EOF

vault policy write leviosa-staging - <<EOF
path "leviosa_staging/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
EOF
```

#### Step 6: Deploy Staging

```bash
make ansible-deploy-staging
```

## Post-Deployment Tasks

### 1. Verify Vault Status

```bash
# Check Vault status
docker exec leviosa_vault vault status

# Verify sealed status is false
```

### 2. Create Auto-Unseal Script (Optional)

For better automation, consider using AWS KMS or HSM for auto-unseal. This removes the need to manually unseal after restarts.

### 3. Set Up Vault Backup

Configure periodic backups of Vault data:

```bash
# Add to crontab or systemd timer
0 2 * * * docker exec leviosa_vault vault operator raft snapshot save /backup/vault-snapshot-$(date +\%Y\%m\%d).snap
```

### 4. Monitor Vault Logs

```bash
make prod-logs | grep vault
```

## Security Checklist

- [ ] Root token is 32+ characters and securely stored
- [ ] Unseal keys are distributed to trusted individuals
- [ ] Vault is not accessible from public internet (only internal network)
- [ ] TLS is enabled (remove `tls_disable = 1` and configure certificates)
- [ ] AppRole authentication is configured instead of root tokens
- [ ] Audit logging is enabled
- [ ] Backup strategy is in place
- [ ] Recovery keys are stored offline (e.g., password manager)
- [ ] Vault version is up to date
- [ ] Network policies restrict access to Vault port 8200

## Troubleshooting

### CAP_SETFCAP Error (Vault Restart Loop)

**Symptom:** Vault container stuck in restart loop with logs showing:
```
unable to set CAP_SETFCAP effective capability: Operation not permitted
```

**Solution:** This occurs when Vault tries to use memory locking (mlock) without proper capabilities. Fix by adding `disable_mlock = true` to Vault config:

```hcl
# In vault.hcl
disable_mlock = true
```

**Note:** The development setup uses `privileged: true` as a workaround. For production, use `disable_mlock = true` in the config file instead (already included in the template above).

**Important:** If using Vault 1.16+, always pin to a specific version instead of `:latest` to avoid breaking changes:

```yaml
# Use specific version
image: hashicorp/vault:1.15
# NOT: image: hashicorp/vault:latest
```

### Vault Sealed After Restart

If Vault restarts and becomes sealed:

```bash
# SSH to VPS
make prod-ssh

# Unseal with 3 of 5 keys
docker exec -it leviosa_vault vault operator unseal <KEY_1>
docker exec -it leviosa_vault vault operator unseal <KEY_2>
docker exec -it leviosa_vault vault operator unseal <KEY_3>
```

### Backend Cannot Connect to Vault

1. Check Vault is running: `docker ps | grep vault`
2. Check Vault status: `docker exec leviosa_vault vault status`
3. Check environment variables in backend container: `docker exec leviosa_backend env | grep VAULT`
4. Verify network connectivity: `docker exec leviosa_backend ping leviosa_vault`

### Lost Unseal Keys or Root Token

If unseal keys or root token are lost, the only option is to:
1. Destroy the Vault container and data
2. Re-initialize Vault
3. Re-create all secrets

This is why secure backup is critical!

## Additional Resources

- [Vault Production Guide](https://developer.hashicorp.com/vault/docs/operator/production)
- [Vault Deployment Guide](https://developer.hashicorp.com/vault/docs/install/deployment-guide)
- [Vault Best Practices](https://developer.hashicorp.com/vault/docs/operations/best-practices)

## Timeline Estimate

- Setup and configuration: 1-2 hours
- Initial deployment and testing: 1 hour
- Migration from dev to production: 2-3 hours (including data migration if needed)
- Total: 4-6 hours for a complete migration
