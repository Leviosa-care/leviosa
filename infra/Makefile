# Leviosa Infrastructure Makefile
# This Makefile provides commands for managing Leviosa infrastructure
# including Terraform and Ansible deployments

.PHONY: help help-terraform help-ansible \
        init-staging init-production \
        plan-staging plan-production \
        apply-staging apply-production \
        update-staging-vault update-production-vault \
        deploy-staging deploy-production \
        destroy-staging destroy-production \
        clean

# Default target
help:
	@echo "Leviosa Infrastructure Commands"
	@echo "=================================="
	@echo ""
	@echo "Quick Start:"
	@echo "  make init-staging         - Initialize staging infrastructure"
	@echo "  make apply-staging        - Apply staging infrastructure changes"
	@echo "  make deploy-staging       - Deploy staging application"
	@echo ""
	@echo "Staging Environment:"
	@echo "  make plan-staging         - Show planned changes for staging"
	@echo "  make apply-staging        - Apply changes for staging"
	@echo "  make update-staging-vault - Update Ansible vault with staging credentials"
	@echo "  make deploy-staging       - Deploy staging application"
	@echo "  make destroy-staging      - Destroy staging infrastructure"
	@echo ""
	@echo "Production Environment:"
	@echo "  make plan-production      - Show planned changes for production"
	@echo "  make apply-production     - Apply changes for production"
	@echo "  make update-production-vault - Update Ansible vault with production credentials"
	@echo "  make deploy-production    - Deploy production application"
	@echo "  make destroy-production   - Destroy production infrastructure"
	@echo ""
	@echo "For more commands, see:"
	@echo "  make help-terraform        - Show Terraform-specific commands"
	@echo "  make help-ansible         - Show Ansible-specific commands"

# ============================================
# Staging Environment (default workspace)
# ============================================

init-staging:
	@echo "Initializing staging infrastructure (migrating state if needed)..."
	@cd terraform && \
	terraform init -migrate-state && \
	terraform workspace select default 2>/dev/null || true

plan-staging:
	@echo "Planning staging infrastructure..."
	@cd terraform && \
	terraform workspace select default 2>/dev/null || true && \
	terraform plan -var-file=terraform.tfvars.staging

apply-staging:
	@echo "Applying staging infrastructure..."
	@cd terraform && \
	terraform workspace select default 2>/dev/null || true && \
	terraform apply -var-file=terraform.tfvars.staging

update-staging-vault:
	@echo "Updating staging vault..."
	@./scripts/update-staging-vault.sh

deploy-staging: update-staging-vault
	@echo "Deploying staging application..."
	@cd ansible && \
	make deploy-leviosa-staging

destroy-staging:
	@echo "⚠️  WARNING: This will destroy staging infrastructure!"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@cd terraform && \
	terraform workspace select default 2>/dev/null || true && \
	terraform destroy -var-file=terraform.tfvars.staging

# ============================================
# Production Environment
# ============================================

init-production:
	@echo "Initializing production infrastructure (migrating state if needed)..."
	@cd terraform && \
	terraform init -migrate-state && \
	terraform workspace select production 2>/dev/null || terraform workspace new production

plan-production:
	@echo "Planning production infrastructure..."
	@cd terraform && \
	terraform workspace select production 2>/dev/null || terraform workspace new production && \
	terraform plan -var-file=terraform.tfvars.production

apply-production:
	@echo "Applying production infrastructure..."
	@cd terraform && \
	terraform workspace select production 2>/dev/null || terraform workspace new production && \
	terraform apply -var-file=terraform.tfvars.production

update-production-vault:
	@echo "Updating production vault..."
	@./scripts/update-production-vault.sh

deploy-production: update-production-vault
	@echo "Deploying production application..."
	@cd ansible && \
	make deploy-leviosa-production

destroy-production:
	@echo "⚠️  WARNING: This will destroy production infrastructure!"
	@read -p "Are you sure? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@cd terraform && \
	terraform workspace select production 2>/dev/null || terraform workspace new production && \
	terraform destroy -var-file=terraform.tfvars.production

# ============================================
# Help Targets
# ============================================

help-terraform:
	@echo "Terraform Commands"
	@echo "=================="
	@echo "See terraform/ directory for Terraform-specific commands:"
	@echo "  cd terraform && make help"

help-ansible:
	@echo "Ansible Commands"
	@echo "================"
	@echo "See ansible/ directory for Ansible-specific commands:"
	@echo "  cd ansible && make help"

# ============================================
# Utilities
# ============================================

clean:
	@echo "Cleaning temporary files..."
	@cd terraform && make clean 2>/dev/null || true
	@echo "Clean complete."
