# Leviosa Terraform Makefile
# Usage: make <target>
#
# Workspace Strategy:
# - default workspace = staging
# - production workspace = production
#
# Each workspace has separate state and AWS resources

export AWS_PROFILE := terraform-leviosa

.PHONY: help init init-migrate init-staging init-production \
        plan apply output credentials destroy destroy-staging destroy-production destroy-backend \
        fmt validate refresh show graph clean \
        workspace-list workspace-new workspace-select \
        plan-staging apply-staging plan-production apply-production \
        update-staging-vault update-production-vault \
        server-info server-ssh server-deploy \
        dns-info dns-verify

# Default target
help:
	@echo "Leviosa Terraform Commands"
	@echo "============================"
	@echo ""
	@echo "Workspaces:"
	@echo "  make workspace-list           - List all workspaces"
	@echo "  make workspace-new NAME=name  - Create new workspace"
	@echo "  make workspace-select         - Select workspace interactively"
	@echo ""
	@echo "Setup:"
	@echo "  make init                     - Initialize Terraform"
	@echo "  make init-migrate             - Re-initialize and migrate state (after backend changes)"
	@echo "  make init-staging             - Initialize + migrate state for staging workspace"
	@echo "  make init-production          - Initialize + migrate state for production workspace"
	@echo "  make validate                 - Validate configuration files"
	@echo ""
	@echo "Staging Environment (default workspace):"
	@echo "  make plan-staging             - Show planned changes for staging"
	@echo "  make apply-staging            - Apply changes for staging"
	@echo "  make update-staging-vault     - Update Ansible vault with staging credentials"
	@echo ""
	@echo "Production Environment:"
	@echo "  make plan-production          - Show planned changes for production"
	@echo "  make apply-production         - Apply changes for production"
	@echo "  make update-production-vault  - Update Ansible vault with production credentials"
	@echo ""
	@echo "Other:"
	@echo "  make refresh                  - Refresh state file"
	@echo "  make show                     - Show current state"
	@echo ""
	@echo "Outputs:"
	@echo "  make output                   - Show all outputs"
	@echo "  make credentials              - Show AWS S3 credentials"
	@echo "  make server-info              - Show VPS connection details"
	@echo "  make dns-info                 - Show DNS records"
	@echo ""
	@echo "Server:"
	@echo "  make server-ssh               - SSH into the VPS"
	@echo "  make server-deploy            - Show deployment guide"
	@echo ""
	@echo "Destruction:"
	@echo "  make destroy                  - Destroy all resources in current workspace"
	@echo "  make destroy-backend          - Destroy backend resources (careful!)"
	@echo ""
	@echo "Utilities:"
	@echo "  make fmt                      - Format Terraform files"
	@echo "  make graph                    - Generate dependency graph"
	@echo ""
	@echo "Current workspace: $$(terraform workspace show)"

# Generic plan/apply (defaults to staging)
plan:
	@$(MAKE) plan-staging

apply:
	@$(MAKE) apply-staging

# ============================================
# Workspace Management
# ============================================

# List all workspaces
workspace-list:
	@echo "Available workspaces:"
	@terraform workspace list
	@echo ""
	@echo "Current workspace: $$(terraform workspace show)"

# Create new workspace
workspace-new:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME parameter is required"; \
		echo "Usage: make workspace-new NAME=<workspace-name>"; \
		exit 1; \
	fi
	@echo "Creating new workspace: $(NAME)"
	terraform workspace new $(NAME)

# Select workspace interactively
workspace-select:
	@echo "Available workspaces:"
	@terraform workspace list
	@echo ""
	@read -p "Enter workspace name: " workspace; \
	terraform workspace select $$workspace

# ============================================
# Terraform Operations
# ============================================

# Initialize Terraform (standard)
init:
	@echo "Initializing Terraform..."
	terraform init

# Initialize and migrate state to new backend (use after backend config changes)
init-migrate:
	@echo "Initializing Terraform with state migration..."
	terraform init -migrate-state

# Initialize for staging (default workspace) with state migration
init-staging:
	@echo "Initializing Terraform for staging (migrating state if needed)..."
	terraform init -migrate-state
	@terraform workspace select default 2>/dev/null || true

# Initialize for production workspace with state migration
init-production:
	@echo "Initializing Terraform for production (migrating state if needed)..."
	terraform init -migrate-state
	@terraform workspace select production 2>/dev/null || terraform workspace new production

# Staging environment targets (uses default workspace)
plan-staging:
	@echo "Switching to staging workspace (default)..."
	@terraform workspace select default 2>/dev/null || true
	@echo "Showing planned changes for staging..."
	terraform plan -var-file=terraform.tfvars.staging $(STAGING_TARGETS)

apply-staging:
	@echo "Switching to staging workspace (default)..."
	@terraform workspace select default 2>/dev/null || true
	@echo "Applying Terraform changes for staging..."
	terraform apply -var-file=terraform.tfvars.staging $(STAGING_TARGETS)

# Production environment targets
# Only manages AWS resources (S3 + IAM). The VPS and DNS are shared and already
# live in the staging/default workspace state — do not recreate them here.
# Staging targets: VPS + Cloudflare DNS + staging S3/IAM
# Excludes backend.tf resources (terraform state bucket) which must persist
STAGING_TARGETS := \
	-target=hcloud_server.manager \
	-target=cloudflare_dns_record.leviosa_care_a \
	-target=cloudflare_dns_record.leviosa_care_aaaa \
	-target=cloudflare_dns_record.www_leviosa_care \
	-target=cloudflare_dns_record.staging_leviosa_care \
	-target=cloudflare_dns_record.admin_leviosa_care \
	-target=cloudflare_dns_record.staff_leviosa_care \
	-target=cloudflare_dns_record.admin_staging_leviosa_care \
	-target=cloudflare_dns_record.staff_staging_leviosa_care \
	-target=aws_s3_bucket.bucket \
	-target=aws_s3_bucket_versioning.assets_versioning \
	-target=aws_s3_bucket_server_side_encryption_configuration.assets_encryption \
	-target=aws_s3_bucket_public_access_block.assets_block \
	-target=aws_s3_bucket_cors_configuration.assets_cors \
	-target=aws_iam_user.app_user \
	-target=aws_iam_access_key.app_user_key \
	-target=aws_iam_user_policy.app_s3_policy

# Production targets: production S3 + IAM only (VPS lives in staging workspace)
PRODUCTION_TARGETS := \
	-target=aws_s3_bucket.bucket \
	-target=aws_s3_bucket_versioning.assets_versioning \
	-target=aws_s3_bucket_server_side_encryption_configuration.assets_encryption \
	-target=aws_s3_bucket_public_access_block.assets_block \
	-target=aws_s3_bucket_cors_configuration.assets_cors \
	-target=aws_iam_user.app_user \
	-target=aws_iam_access_key.app_user_key \
	-target=aws_iam_user_policy.app_s3_policy \
	-target=aws_iam_policy.app_leviosa \
	-target=aws_iam_user.app_leviosa \
	-target=aws_iam_user_policy_attachment.app_leviosa \
	-target=aws_iam_access_key.app_leviosa

plan-production:
	@echo "Switching to production workspace..."
	@terraform workspace select production 2>/dev/null || terraform workspace new production
	@echo "Showing planned changes for production (AWS resources only)..."
	terraform plan -var-file=terraform.tfvars.production $(PRODUCTION_TARGETS)

apply-production:
	@echo "Switching to production workspace..."
	@terraform workspace select production 2>/dev/null || terraform workspace new production
	@echo "Applying Terraform changes for production (AWS resources only)..."
	terraform apply -var-file=terraform.tfvars.production $(PRODUCTION_TARGETS)

# Update Ansible vault with credentials
update-staging-vault:
	@../scripts/update-staging-vault.sh

update-production-vault:
	@../scripts/update-production-vault.sh

# Validate configuration
validate:
	@echo "Validating Terraform configuration..."
	terraform fmt -check
	terraform validate

# Format configuration files
fmt:
	@echo "Formatting Terraform files..."
	terraform fmt -recursive

# Refresh state file
refresh:
	@echo "Refreshing Terraform state..."
	terraform refresh

# Show current state
show:
	@echo "Showing Terraform state..."
	terraform show

# Show all outputs
output:
	@echo "Terraform Outputs:"
	@echo ""
	terraform output -json | jq -r 'to_entries | .[] | "\(.key): \(.value.value // .value)"' 2>/dev/null || terraform output

# Show IAM credentials for .env file
credentials:
	@echo "==========================================="
	@echo "AWS S3 Credentials for .env file"
	@echo "==========================================="
	@echo ""
	@echo "S3_BUCKET=$$(terraform output -raw media_bucket_name 2>/dev/null || echo 'staging-leviosa-assets')"
	@echo "S3_REGION=$$(terraform output -raw media_bucket_region 2>/dev/null || echo 'eu-central-1')"
	@echo "AWS_ACCESS_KEY_ID=$$(terraform output -raw app_user_access_key_id)"
	@echo "AWS_SECRET_ACCESS_KEY=$$(terraform output -raw app_user_secret_access_key)"
	@echo ""
	@echo "==========================================="
	@echo "Save these credentials securely!"
	@echo "==========================================="

# Generate dependency graph (requires graphviz)
graph:
	@echo "Generating dependency graph..."
	terraform graph | dot -Tpng > terraform-graph.png
	@echo "Graph saved to terraform-graph.png"

# Destroy all resources in current workspace
destroy:
	@echo "Destroying all Terraform resources in current workspace..."
	terraform destroy

# Destroy staging resources (VPS + Cloudflare DNS + staging S3/IAM)
# Does NOT destroy the Terraform state backend (backend.tf)
destroy-staging:
	@echo "Switching to staging workspace (default)..."
	@terraform workspace select default 2>/dev/null || true
	@echo "Destroying staging resources (excluding backend)..."
	terraform destroy -var-file=terraform.tfvars.staging $(STAGING_TARGETS)

# Destroy production resources (S3 + IAM only — VPS lives in staging workspace)
destroy-production:
	@echo "Switching to production workspace..."
	@terraform workspace select production 2>/dev/null || terraform workspace new production
	@echo "Destroying production resources (AWS only)..."
	terraform destroy -var-file=terraform.tfvars.production $(PRODUCTION_TARGETS)

# Destroy backend resources (use with caution!)
destroy-backend:
	@echo "⚠️  WARNING: This will destroy the Terraform state backend!"
	@echo "You will need to recreate it afterwards."
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	terraform destroy -target=aws_dynamodb_table.terraform_locks
	terraform destroy -target=aws_s3_bucket.terraform_state
	terraform destroy -target=aws_s3_bucket_versioning.terraform_state_versioning
	terraform destroy -target=aws_s3_bucket_server_side_encryption_configuration.terraform_state_sse

# Clean temporary files
clean:
	@echo "Cleaning temporary files..."
	-rm -f terraform-graph.png
	-rm -rf .terraform/
	-rm -f .terraform.lock.hcl
	-rm -f terraform.tfstate.backup
	@echo "Clean complete. Run 'make init' to reinitialize."

# ============================================
# Server Commands
# ============================================

# Show server connection details
server-info:
	@echo "==========================================="
	@echo "VPS Connection Details"
	@echo "==========================================="
	@echo ""
	@echo "IPv4: $$(terraform output -raw server_ipv4_address)"
	@echo "IPv6: $$(terraform output -raw server_ipv6_address)"
	@echo ""
	@echo "SSH connection:"
	@echo "  $$(terraform output -raw ssh_connection_string)"
	@echo ""
	@echo "==========================================="

# SSH into the server
server-ssh:
	@ssh root@$$(terraform output -raw server_ipv4_address)

# Show deployment guide
server-deploy:
	@echo "==========================================="
	@echo "Leviosa VPS Deployment Guide"
	@echo "==========================================="
	@echo ""
	@echo "Server: leviosa-staging"
	@echo "IPv4: $$(terraform output -raw server_ipv4_address)"
	@echo "IPv6: $$(terraform output -raw server_ipv6_address)"
	@echo ""
	@echo "Connect to server:"
	@echo "  $$(terraform output -raw ssh_connection_string)"
	@echo ""
	@echo "Run Ansible setup:"
	@echo "  cd ../ansible"
	@echo "  make setup"
	@echo ""
	@echo "==========================================="

# ============================================
# DNS Commands
# ============================================

# Show application URL and DNS records
dns-info:
	@echo "==========================================="
	@echo "Application DNS Information"
	@echo "==========================================="
	@echo ""
	@echo "Application URL: https://$$(terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care')"
	@echo ""
	@echo "DNS Records configured:"
	@echo "  A (main):     $$(terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care') → $$(terraform output -raw server_ipv4_address)"
	@echo "  CNAME (www):  www.$$(terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care') → $$(terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care')"
	@echo "  CNAME (cdn):  cdn.$$(terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care') → $$(terraform output -raw cloudfront_domain_name)"
	@echo ""
	@echo "==========================================="

# Show Cloudflare dashboard link for verification
dns-verify:
	@echo "==========================================="
	@echo "Cloudflare Dashboard"
	@echo "==========================================="
	@echo ""
	@echo "Visit the Cloudflare dashboard to verify DNS:"
	@echo "  https://dash.cloudflare.com/$(cat terraform.tfvars 2>/dev/null | grep zone_id | cut -d'"' -f2)/$$(terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care')/dns"
	@echo ""
	@echo "==========================================="
