#!/usr/bin/env bash
# Refreshes production AWS credentials in Ansible group_vars after a Terraform apply.
#
# Usage: ./scripts/update-production-vault.sh
# Or via: make update-production-vault

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TERRAFORM_DIR="$SCRIPT_DIR/../terraform"
ANSIBLE_DIR="$SCRIPT_DIR/../ansible"
PRODUCTION_VAULT="$ANSIBLE_DIR/group_vars/leviosa_production.yml"
WORKSPACE="production"

# ── Prerequisites ─────────────────────────────────────────────────────
for cmd in terraform python3; do
    command -v "$cmd" >/dev/null 2>&1 || { echo "ERROR: $cmd is required but not found"; exit 1; }
done

# ── Helpers ───────────────────────────────────────────────────────────
set_yaml_keys() {
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

cd "$TERRAFORM_DIR"

echo ""
echo "=========================================="
echo "Production Credentials Update"
echo "=========================================="
echo ""
echo "WARNING: PRODUCTION ENVIRONMENT"
echo ""
read -p "Continue? [y/N] " confirm && [ "$confirm" = "y" ] || exit 1
echo ""

echo "Switching to production workspace..."
terraform workspace select "$WORKSPACE" 2>/dev/null || {
    echo "ERROR: Production workspace not found. Run: make init-production"
    exit 1
}

echo "Reading Terraform outputs..."
ACCESS_KEY_ID=$(terraform output -raw app_user_access_key_id)
SECRET_ACCESS_KEY=$(terraform output -raw app_user_secret_access_key)
S3_BUCKET=$(terraform output -raw media_bucket_name)
S3_REGION=$(terraform output -raw media_bucket_region)

[ -n "$ACCESS_KEY_ID" ]    || { echo "ERROR: failed to read app_user_access_key_id";    exit 1; }
[ -n "$SECRET_ACCESS_KEY" ] || { echo "ERROR: failed to read app_user_secret_access_key"; exit 1; }

echo ""
echo "S3 Bucket : $S3_BUCKET"
echo "S3 Region : $S3_REGION"
echo "IAM User  : $(terraform output -raw app_user_access_key_id | cut -c1-8)..."
echo ""

echo "Updating $PRODUCTION_VAULT ..."
set_yaml_keys "$PRODUCTION_VAULT" \
    "aws_access_key_id=$ACCESS_KEY_ID" \
    "aws_secret_access_key=$SECRET_ACCESS_KEY" \
    "s3_bucket=$S3_BUCKET" \
    "s3_region=$S3_REGION"

echo ""
echo "Production vault updated."
echo ""
echo "Next: cd ansible && make deploy-production"
echo ""
