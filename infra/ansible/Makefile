# Leviosa Ansible Makefile
# Usage: make <target>

.PHONY: help setup deploy deploy-staging restart restart-staging backup check-clean lint \
        ssh-info status staging-status logs staging-logs health-check \
        install-deps vars-test

# Default target
help:
	@echo "Leviosa Ansible Commands"
	@echo "============================"
	@echo ""
	@echo "Setup:"
	@echo "  make setup           - Complete VPS setup (first time only)"
	@echo "  make deploy          - Deploy/update production"
	@echo "  make deploy-staging  - Deploy/update staging (+ update Caddy)"
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

# Get VPS IP from Terraform output
VPS_IP=$(shell cd ../terraform && terraform output -raw server_ipv4_address 2>/dev/null || echo "")
DOMAIN=$(shell cd ../terraform && terraform output -raw domain_name 2>/dev/null || echo "")
BACKUP_BUCKET=$(shell cd ../terraform && terraform output -raw backup_bucket_name 2>/dev/null || echo "")

# Get SSH public key
SSH_KEY=$(shell cat ~/.ssh/leviosa.pub 2>/dev/null || cat ~/.ssh/id_rsa.pub 2>/dev/null || echo "")

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
		echo "Error: Could not get VPS IP from Terraform."; \
		echo "Run 'cd ../terraform && terraform output server_ipv4_address'"; \
		exit 1; \
	fi
	@echo "VPS IP: $(VPS_IP)"
	@echo "Domain: $(DOMAIN)"
	@echo "SSH Key: $(SSH_KEY)"
	@echo ""
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo ""
	@echo "Running setup playbook..."
	ansible-playbook -i inventory/hosts.yml playbooks/site.yml \
		--ask-vault-pass \
		-e "ansible_host=$(VPS_IP)" \
		-e "deploy_ssh_public_key=$(SSH_KEY)" \
		-e "app_domain=$(DOMAIN)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/secrets.yml"

# Deploy production application
deploy:
	@echo "==========================================="
	@echo "Leviosa Production Deployment"
	@echo "==========================================="
	@echo ""
	@if [ -z "$(VPS_IP)" ]; then \
		echo "Error: Could not get VPS IP from Terraform."; \
		exit 1; \
	fi
	@echo "Deploying to: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts.yml playbooks/deploy.yml \
		--ask-vault-pass \
		-e "ansible_host=$(VPS_IP)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/secrets.yml"

# Deploy staging application (also updates prod Caddy for staging domains)
deploy-staging:
	@echo "==========================================="
	@echo "Leviosa Staging Deployment"
	@echo "==========================================="
	@echo ""
	@if [ -z "$(VPS_IP)" ]; then \
		echo "Error: Could not get VPS IP from Terraform."; \
		exit 1; \
	fi
	@if [ ! -f "group_vars/leviosa_staging.vault.yml" ]; then \
		echo "Error: group_vars/leviosa_staging.vault.yml not found."; \
		echo "Copy group_vars/leviosa_staging.vault.example.yml and encrypt it:"; \
		echo "  cp group_vars/leviosa_staging.vault.example.yml group_vars/leviosa_staging.vault.yml"; \
		echo "  ansible-vault encrypt group_vars/leviosa_staging.vault.yml"; \
		exit 1; \
	fi
	@echo "Deploying staging to: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts.yml playbooks/deploy-staging.yml \
		--ask-vault-pass \
		-e "ansible_host=$(VPS_IP)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/secrets.yml" \
		-e "@group_vars/leviosa_staging.yml" \
		-e "@group_vars/leviosa_staging.vault.yml"

# Quick restart: pull latest app image and restart only the app container
restart:
	@echo "Pulling latest prod image and restarting..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa && docker compose pull app && docker compose up -d app"

# Quick restart for staging
restart-staging:
	@echo "Pulling latest staging image and restarting..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa-staging && docker compose pull app && docker compose up -d app"

# Configure database backups
backup:
	@echo "==========================================="
	@echo "Database Backup Configuration"
	@echo "==========================================="
	@echo ""
	@if [ -z "$(VPS_IP)" ]; then \
		echo "Error: Could not get VPS IP from Terraform."; \
		exit 1; \
	fi
	@if [ -z "$(BACKUP_BUCKET)" ]; then \
		echo "Error: Could not get backup bucket from Terraform."; \
		exit 1; \
	fi
	@echo "VPS IP: $(VPS_IP)"
	@echo "Backup Bucket: $(BACKUP_BUCKET)"
	@echo ""
	@read -p "Enter AWS Access Key ID: " aws_key; \
	read -sp "Enter AWS Secret Access Key: " aws_secret; \
	echo ""; \
	ansible-playbook -i inventory/hosts.yml playbooks/backup.yml \
		-e "ansible_host=$(VPS_IP)" \
		-e "backup_s3_bucket=$(BACKUP_BUCKET)" \
		-e "backup_aws_access_key=$$aws_key" \
		-e "backup_aws_secret_key=$$aws_secret"

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
	@ansible-playbook -i inventory/hosts.yml playbooks/site.yml \
		-e "ansible_host=$(VPS_IP)" \
		--tags debug --skip-tags never \
		-e "ansible_debug=true" 2>/dev/null || true
	@echo "To see all variables, run:"
	@echo "  ansible-inventory -i inventory/hosts.yml --list"

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
	@echo "Connect as deploy user:"
	@echo "  ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP)"
	@echo ""
	@echo "==========================================="

# Check production service status
status:
	@echo "Checking production service status..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa && docker compose ps"

# Check staging service status
staging-status:
	@echo "Checking staging service status..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa-staging && docker compose ps"

# View production logs
logs:
	@echo "Streaming production logs (Ctrl+C to exit)..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa && docker compose logs -f"

# View staging logs
staging-logs:
	@echo "Streaming staging logs (Ctrl+C to exit)..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa-staging && docker compose logs -f"

# Run health check
health-check:
	@echo "Running health check..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "/opt/leviosa/scripts/health-check.sh"

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
