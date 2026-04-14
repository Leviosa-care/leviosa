#!/usr/bin/env bash
# Refreshes staging AWS resources credentials in Ansible vault.
# Use this to update vault credentials after Terraform apply.
#
# Usage: ./scripts/update-staging-vault.sh
# Or via: make update-staging-vault

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TERRAFORM_DIR="$SCRIPT_DIR/../terraform"
ANSIBLE_DIR="$SCRIPT_DIR/../ansible"
STAGING_VAULT="$ANSIBLE_DIR/group_vars/leviosa_staging.yml"
WORKSPACE="default"

# ── Prerequisites ────────────────────────────────────────────────────
for cmd in terraform python3 jq; do
    command -v "$cmd" >/dev/null 2>&1 || { echo "❌  $cmd is required but not found"; exit 1; }
done

# ── Helpers ───────────────────────────────────────────────────────────
set_yaml_keys() {
    # $1 = file path, remaining args = key=value pairs
    python3 - "$@" <<'PYEOF'
import yaml, sys

path = sys.argv[1]
data = yaml.safe_load(open(path)) or {}
for pair in sys.argv[2:]:
    key, val = pair.split('=', 1)
    data[key] = val
with open(path, 'w') as f:
    yaml.dump(data, f, default_flow_style=False, allow_unicode=True)
PYEOF
}

# ── Navigate to Terraform directory ────────────────────────────────────
cd "$TERRAFORM_DIR"

# ── Show what will happen ─────────────────────────────────────────────
echo ""
echo "=========================================="
echo "Staging Credentials Update"
echo "=========================================="
echo ""
echo "This will:"
echo "  1. Switch to workspace: $WORKSPACE (staging)"
echo "  2. Read Terraform outputs"
echo "  3. Update Ansible vault with credentials"
echo ""
read -p "Continue? [y/N] " confirm && [ "$confirm" = "y" ] || exit 1
echo ""

# ── Switch to staging workspace ───────────────────────────────────────
echo "Switching to staging workspace..."
terraform workspace select "$WORKSPACE" 2>/dev/null || true

# ── Check for existing resources ───────────────────────────────────────
echo "Checking for existing IAM user..."
EXISTING_USER=$(terraform output -raw iam_user_name 2>/dev/null || echo "")
if [ -z "$EXISTING_USER" ]; then
    echo "❌  No existing IAM user found in $WORKSPACE workspace."
    echo ""
    echo "You need to run Terraform first:"
    echo "  terraform apply -var-file=terraform.tfvars.staging"
    exit 1
fi

echo "  ✓ Found: $EXISTING_USER"

# ── Extract outputs ─────────────────────────────────────────────────────
echo ""
echo "Extracting Terraform outputs..."
VAULT_ACCESS_KEY_ID=$(terraform output -raw vault_user_access_key_id)
VAULT_SECRET_ACCESS_KEY=$(terraform output -raw vault_user_secret_access_key)
VAULT_KMS_KEY_ID=$(terraform output -raw vault_kms_key_id)
S3_BUCKET=$(terraform output -raw s3_bucket_name 2>/dev/null || echo "${WORKSPACE}-leviosa-assets")
S3_REGION=$(terraform output -raw s3_bucket_region 2>/dev/null || echo "eu-central-1")
CDN_URL=$(terraform output -raw cdn_url)
BACKUP_BUCKET=$(terraform output -raw backup_bucket_name)
LOKI_ACCESS_KEY_ID=$(terraform output -raw loki_s3_access_key_id)
LOKI_SECRET_ACCESS_KEY=$(terraform output -raw loki_s3_secret_access_key)
LOKI_S3_BUCKET=$(terraform output -raw loki_s3_bucket_name 2>/dev/null || echo "${WORKSPACE}-leviosa-loki-logs")

[ -n "$VAULT_ACCESS_KEY_ID" ] || { echo "❌  Failed to get vault_user_access_key_id"; exit 1; }
[ -n "$VAULT_SECRET_ACCESS_KEY" ] || { echo "❌  Failed to get vault_user_secret_access_key"; exit 1; }

# ── Display credentials ─────────────────────────────────────────────────
echo ""
echo "=========================================="
echo "Staging Credentials"
echo "=========================================="
echo ""
echo "Workspace: $WORKSPACE"
echo "IAM User: $EXISTING_USER"
echo "S3 Bucket: $S3_BUCKET"
echo "CDN URL: $CDN_URL"
echo "Backup Bucket: $BACKUP_BUCKET"
echo "Loki S3 Bucket: $LOKI_S3_BUCKET"
echo ""

# ── Update Ansible vault ───────────────────────────────────────────────
echo "Updating Ansible vault..."
set_yaml_keys "$STAGING_VAULT" \
    "vault_aws_access_key_id=$VAULT_ACCESS_KEY_ID" \
    "vault_aws_secret_access_key=$VAULT_SECRET_ACCESS_KEY" \
    "vault_kms_key_id=$VAULT_KMS_KEY_ID" \
    "vault_s3_bucket=$S3_BUCKET" \
    "vault_s3_region=$S3_REGION" \
    "vault_cdn_url=$CDN_URL" \
    "vault_backup_bucket=$BACKUP_BUCKET" \
    "vault_loki_s3_access_key_id=$LOKI_ACCESS_KEY_ID" \
    "vault_loki_s3_secret_access_key=$LOKI_SECRET_ACCESS_KEY" \
    "vault_loki_s3_bucket=$LOKI_S3_BUCKET"

echo "  ✓ Staging vault updated"
echo ""

# ── Summary ────────────────────────────────────────────────────────────
echo "=========================================="
echo "✅  Staging Credentials Updated"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  - Deploy: cd ../ansible && make deploy-staging"
echo ""
echo "Vault file updated: $STAGING_VAULT"
echo ""
