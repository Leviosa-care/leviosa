#!/usr/bin/env bash
# Refreshes staging AWS credentials in Ansible group_vars after a Terraform apply.
#
# Usage: ./scripts/update-staging-vault.sh
# Or via: make update-staging-vault

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TERRAFORM_DIR="$SCRIPT_DIR/../terraform"
ANSIBLE_DIR="$SCRIPT_DIR/../ansible"
STAGING_VAULT="$ANSIBLE_DIR/group_vars/leviosa_staging.yml"
WORKSPACE="default"

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
echo "Staging Credentials Update"
echo "=========================================="
echo ""
read -p "Continue? [y/N] " confirm && [ "$confirm" = "y" ] || exit 1
echo ""

echo "Switching to staging workspace ($WORKSPACE)..."
terraform workspace select "$WORKSPACE" 2>/dev/null || true

echo "Reading Terraform outputs..."
ACCESS_KEY_ID=$(terraform output -no-color -raw app_user_access_key_id 2>/dev/null) \
  || { echo "ERROR: failed to read app_user_access_key_id — run 'make apply-staging' first"; exit 1; }
SECRET_ACCESS_KEY=$(terraform output -no-color -raw app_user_secret_access_key 2>/dev/null) \
  || { echo "ERROR: failed to read app_user_secret_access_key — run 'make apply-staging' first"; exit 1; }
S3_BUCKET=$(terraform output -no-color -raw media_bucket_name 2>/dev/null) \
  || { echo "ERROR: failed to read media_bucket_name — run 'make apply-staging' first"; exit 1; }
S3_REGION=$(terraform output -no-color -raw media_bucket_region 2>/dev/null) \
  || { echo "ERROR: failed to read media_bucket_region — run 'make apply-staging' first"; exit 1; }

[ -n "$ACCESS_KEY_ID" ]     || { echo "ERROR: app_user_access_key_id is empty";     exit 1; }
[ -n "$SECRET_ACCESS_KEY" ] || { echo "ERROR: app_user_secret_access_key is empty"; exit 1; }
[ -n "$S3_BUCKET" ]         || { echo "ERROR: media_bucket_name is empty";          exit 1; }
[ -n "$S3_REGION" ]         || { echo "ERROR: media_bucket_region is empty";        exit 1; }

echo ""
echo "S3 Bucket : $S3_BUCKET"
echo "S3 Region : $S3_REGION"
echo "IAM User  : $(terraform output -raw app_user_access_key_id | cut -c1-8)..."
echo ""

echo "Updating $STAGING_VAULT ..."
set_yaml_keys "$STAGING_VAULT" \
    "aws_access_key_id=$ACCESS_KEY_ID" \
    "aws_secret_access_key=$SECRET_ACCESS_KEY" \
    "s3_bucket=$S3_BUCKET" \
    "s3_region=$S3_REGION"

echo ""
echo "Staging vault updated."
echo ""
echo "Next: cd ansible && make deploy-staging"
echo ""
