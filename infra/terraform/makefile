# Leviosa Terraform Makefile
# Usage: make <target>

.PHONY: help init plan apply output credentials destroy destroy-backend \
        fmt validate refresh show graph clean \
        server-info server-ssh server-deploy \
        dns-info dns-email dns-verify

# Default target
help:
	@echo "Leviosa Terraform Commands"
	@echo "============================"
	@echo ""
	@echo "Setup:"
	@echo "  make init          - Initialize Terraform"
	@echo "  make validate      - Validate configuration files"
	@echo ""
	@echo "Operations:"
	@echo "  make plan          - Show planned changes"
	@echo "  make apply         - Apply changes"
	@echo "  make refresh       - Refresh state file"
	@echo "  make show          - Show current state"
	@echo ""
	@echo "Outputs:"
	@echo "  make output        - Show all outputs"
	@echo "  make credentials   - Show AWS S3 credentials"
	@echo "  make server-info   - Show VPS connection details"
	@echo "  make dns-info      - Show DNS records"
	@echo "  make dns-email     - Show email configuration"
	@echo ""
	@echo "Server:"
	@echo "  make server-ssh    - SSH into the VPS"
	@echo "  make server-deploy - Show deployment guide"
	@echo ""
	@echo "Destruction:"
	@echo "  make destroy       - Destroy all resources"
	@echo "  make destroy-backend - Destroy backend resources (careful!)"
	@echo ""
	@echo "Utilities:"
	@echo "  make fmt           - Format Terraform files"
	@echo "  make graph         - Generate dependency graph"

# Initialize Terraform
init:
	@echo "Initializing Terraform..."
	terraform init

# Show planned changes
plan:
	@echo "Showing planned changes..."
	terraform plan

# Apply changes
apply:
	@echo "Applying Terraform changes..."
	terraform apply

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
	@echo "AWS_ACCESS_KEY_ID=$$(terraform output -raw vault_user_access_key_id)"
	@echo "AWS_SECRET_ACCESS_KEY=$$(terraform output -raw vault_user_secret_access_key)"
	@echo ""
	@echo "==========================================="
	@echo "⚠️  Save these credentials securely!"
	@echo "==========================================="

# Generate dependency graph (requires graphviz)
graph:
	@echo "Generating dependency graph..."
	terraform graph | dot -Tpng > terraform-graph.png
	@echo "Graph saved to terraform-graph.png"

# Destroy all resources
destroy:
	@echo "Destroying all Terraform resources..."
	terraform destroy

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

# Show email configuration (SES)
dns-email:
	@terraform output -raw email_setup_status

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
