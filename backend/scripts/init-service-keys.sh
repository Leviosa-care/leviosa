#!/bin/bash

# init-service-keys.sh
# Script to initialize service API keys in Vault

set -e

echo "=== Service API Key Initialization ==="
echo "This script will generate and store API keys for all microservices in Vault."
echo

# Check requirements
if [ -z "$VAULT_TOKEN" ]; then
    echo "Error: VAULT_TOKEN environment variable is required"
    echo "Please set your Vault token:"
    echo "  export VAULT_TOKEN=your-vault-token"
    exit 1
fi

if [ -z "$VAULT_ADDR" ]; then
    echo "Warning: VAULT_ADDR not set, using default: http://localhost:8200"
    export VAULT_ADDR="http://localhost:8200"
fi

echo "Vault Address: $VAULT_ADDR"
echo "Vault Token: ${VAULT_TOKEN:0:8}..."
echo

# Check if Vault is accessible
echo "Testing Vault connectivity..."
if ! vault status > /dev/null 2>&1; then
    echo "Error: Cannot connect to Vault at $VAULT_ADDR"
    echo "Please ensure Vault is running and accessible"
    exit 1
fi
echo "✓ Vault is accessible"
echo

# Check if KV v2 secrets engine is enabled
echo "Checking KV secrets engine..."
if ! vault secrets list | grep -q "secret/"; then
    echo "Enabling KV v2 secrets engine at secret/..."
    vault secrets enable -path=secret kv-v2
else
    echo "✓ KV secrets engine already enabled"
fi
echo

# Run the Go program to generate keys
echo "Generating service API keys..."
cd "$(dirname "$0")/.."
go run scripts/init-service-keys.go

echo
echo "=== Service API Key Initialization Complete ==="
echo
echo "SECURITY NOTES:"
echo "- The generated API keys are displayed above"
echo "- Store them securely in your service configurations"
echo "- Consider using environment variables or secure config management"
echo "- These keys provide full access to internal service endpoints"
echo "- Set up key rotation policies for production use"
echo
echo "For production deployment, consider:"
echo "1. Using Vault Agent for automatic key distribution"
echo "2. Setting up AppRole authentication for services"
echo "3. Implementing automatic key rotation"
echo "4. Adding monitoring and alerting for key usage"
