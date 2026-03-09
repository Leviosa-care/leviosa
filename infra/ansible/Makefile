# Leviosa Ansible Makefile
# Usage: make <target>

.PHONY: help setup deploy deploy-staging restart restart-staging backup check-clean lint \
        ssh-info status staging-status logs staging-logs health-check \
        install-deps vars-test terraform-outputs

# Default target
help:
	@echo "Leviosa Ansible Commands"
	@echo "============================"
	@echo ""
	@echo "Setup:"
	@echo "  make setup           - Complete VPS setup (first time only)"
	@echo "  make deploy          - Deploy/update production"
	@echo "  make deploy-staging  - Deploy/update staging"
	@echo "  make restart         - Pull latest image and restart production (quick)"
	@echo "  make restart-staging - Pull latest image and restart staging (quick)"
	@echo "  make backup          - Configure database backups"
	@echo ""
	@echo "Verification:"
	@echo "  make check-clean     - Check for syntax errors"
	@echo "  make lint            - Run ansible-lint"
	@echo "  make vars-test       - Display effective variables"
	@echo ""
	@echo "Server Operations:"
	@echo "  make ssh-info        - Display SSH connection details"
	@echo "  make status          - Check production service status"
	@echo "  make staging-status  - Check staging service status"
	@echo "  make logs            - View production logs"
	@echo "  make staging-logs    - View staging logs"
	@echo "  make health-check    - Run health check on server"
	@echo ""
	@echo "Dependencies:"
	@echo "  make install-deps    - Install Ansible dependencies"
	@echo "  make terraform-outputs - Show Terraform outputs"

# Deploy user (the non-root user created by the common role)
DEPLOY_USER ?= leviosa
APP_NAME ?= leviosa

# Get VPS IP: try Terraform output first, fall back to inventory parsing
VPS_IP=$(shell cd ../terraform && terraform output -raw server_ipv4_address 2>/dev/null || grep -A1 '\[webserver\]' inventory/hosts | tail -1 | awk '{print $$1}')
DOMAIN=$(shell cd ../terraform && terraform output -raw domain_name 2>/dev/null || echo 'leviosa.care')

# ============================================
# Setup Commands
# ============================================

# Complete VPS setup (run once on fresh server)
setup:
	@echo "==========================================="
	@echo "Leviosa VPS Setup"
	@echo "==========================================="
	@echo ""
	@if [ -z "$(VPS_IP)" ]; then \
		echo "Error: Could not get VPS IP from inventory."; \
		exit 1; \
	fi
	@echo "VPS IP: $(VPS_IP)"
	@echo "Domain: $(DOMAIN)"
	@echo ""
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo ""
	@echo "Running setup playbook..."
	ansible-playbook -i inventory/hosts playbooks/site.yml --ask-vault-pass

# Deploy production application
deploy:
	@echo "==========================================="
	@echo "Leviosa Production Deployment"
	@echo "==========================================="
	@echo ""
	@echo "Deploying to: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts playbooks/deploy.yml --ask-vault-pass

# Deploy staging application
deploy-staging:
	@echo "==========================================="
	@echo "Leviosa Staging Deployment"
	@echo "==========================================="
	@echo ""
	@if [ ! -f "group_vars/leviosa_staging.vault.yml" ]; then \
		echo "Error: group_vars/leviosa_staging.vault.yml not found."; \
		echo "Copy the example and encrypt it:"; \
		echo "  cp group_vars/leviosa_staging.vault.example.yml group_vars/leviosa_staging.vault.yml"; \
		echo "  ansible-vault encrypt group_vars/leviosa_staging.vault.yml"; \
		exit 1; \
	fi
	@echo "Deploying staging to: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts playbooks/deploy-staging.yml --ask-vault-pass

# Quick restart: pull latest app image and restart only the app container
restart:
	@echo "Pulling latest prod image and restarting..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "cd /opt/$(APP_NAME) && docker compose pull && docker compose up -d"

# Quick restart for staging
restart-staging:
	@echo "Pulling latest staging image and restarting..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "cd /opt/$(APP_NAME)-staging && docker compose pull && docker compose up -d"

# Configure database backups
backup:
	@echo "==========================================="
	@echo "Database Backup Configuration"
	@echo "==========================================="
	@echo ""
	@echo "VPS IP: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts playbooks/backup.yml --ask-vault-pass \
		-e "backup_s3_bucket=$${BACKUP_S3_BUCKET:-staging-leviosa-backups}" \
		-e "backup_s3_region=$${BACKUP_S3_REGION:-eu-central-1}"

# ============================================
# Verification Commands
# ============================================

# Check playbooks for syntax errors
check-clean:
	@echo "Checking playbooks for syntax errors..."
	@ansible-playbook playbooks/site.yml --syntax-check
	@ansible-playbook playbooks/deploy.yml --syntax-check
	@ansible-playbook playbooks/deploy-staging.yml --syntax-check
	@ansible-playbook playbooks/backup.yml --syntax-check
	@echo "All playbooks are clean!"

# Run ansible-lint (requires installation)
lint:
	@echo "Running ansible-lint..."
	@which ansible-lint > /dev/null || (echo "Install with: pip install ansible-lint" && exit 1)
	@ansible-lint playbooks/

# Display effective variables
vars-test:
	@echo "Displaying effective variables..."
	@ansible-playbook -i inventory/hosts playbooks/site.yml \
		-e "ansible_host=$(VPS_IP)" \
		--tags debug --skip-tags never 2>/dev/null || true
	@echo "To see all variables, run:"
	@echo "  ansible-inventory -i inventory/hosts --list"

# ============================================
# Server Operations
# ============================================

# Display SSH connection details
ssh-info:
	@echo "==========================================="
	@echo "SSH Connection Details"
	@echo "==========================================="
	@echo ""
	@echo "VPS IP: $(VPS_IP)"
	@echo "Domain: $(DOMAIN)"
	@echo ""
	@echo "Connect as root:"
	@echo "  ssh -i ~/.ssh/leviosa root@$(VPS_IP)"
	@echo ""
	@echo "==========================================="

# Check service status
status:
	@echo "Checking production service status..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "cd /opt/$(APP_NAME) && docker compose ps"

# Check staging service status
staging-status:
	@echo "Checking staging service status..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "cd /opt/$(APP_NAME)-staging && docker compose ps"

# View application logs
logs:
	@echo "Streaming production logs (Ctrl+C to exit)..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "cd /opt/$(APP_NAME) && docker compose logs -f"

# View staging logs
staging-logs:
	@echo "Streaming staging logs (Ctrl+C to exit)..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "cd /opt/$(APP_NAME)-staging && docker compose logs -f"

# Run health check
health-check:
	@echo "Running health check..."
	@ssh -i ~/.ssh/leviosa $(DEPLOY_USER)@$(VPS_IP) "/opt/$(APP_NAME)/scripts/health-check.sh"

# ============================================
# Dependencies
# ============================================

# Install Ansible and required collections
install-deps:
	@echo "Installing Ansible dependencies..."
	@pip install ansible
	@ansible-galaxy collection install community.docker community.general
	@echo "Done! Run 'make check-clean' to verify."

# ============================================
# Terraform Integration
# ============================================

# Show Terraform outputs
terraform-outputs:
	@echo "==========================================="
	@echo "Terraform Outputs"
	@echo "==========================================="
	@echo ""
	@cd ../terraform && terraform output
