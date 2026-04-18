# Leviosa Ansible Makefile
# Usage: make <target>
#
# Workspace-aware deployment for staging and production on the same VPS

.PHONY: help help-deploy help-ops \
        setup setup-staging setup-production \
        deploy deploy-staging deploy-production \
        restart restart-staging restart-production \
        backup check-clean lint \
        ssh-info status staging-status production-status \
        logs staging-logs production-logs \
        health-check install-deps vars-test

# Default target
help:
	@echo "Leviosa Ansible Commands"
	@echo "============================"
	@echo ""
	@echo "Setup (first time only):"
	@echo "  make setup-staging      - Complete VPS setup for staging"
	@echo "  make setup-production   - Complete VPS setup for production"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy-staging     - Deploy/update staging environment"
	@echo "  make deploy-production  - Deploy/update production environment"
	@echo ""
	@echo "Quick Deploy (restart only):"
	@echo "  make restart-staging    - Pull latest image and restart staging"
	@echo "  make restart-production - Pull latest image and restart production"
	@echo ""
	@echo "Operations:"
	@echo "  make backup             - Configure database backups"
	@echo "  make status             - Check production service status"
	@echo "  make staging-status     - Check staging service status"
	@echo "  make logs               - View production logs"
	@echo "  make staging-logs       - View staging logs"
	@echo ""
	@echo "Verification:"
	@echo "  make check-clean        - Check for syntax errors"
	@echo "  make lint               - Run ansible-lint"
	@echo "  make vars-test          - Display effective variables"
	@echo ""
	@echo "Server:"
	@echo "  make ssh-info           - Display SSH connection details"
	@echo "  make health-check       - Run health check on server"
	@echo ""
	@echo "Dependencies:"
	@echo "  make install-deps       - Install Ansible dependencies"
	@echo ""
	@echo "For more commands: make help-deploy or make help-ops"

# Get VPS IP from Terraform output (uses default workspace)
VPS_IP=$(shell cd ../terraform && AWS_PROFILE=terraform-leviosa terraform workspace select default >/dev/null 2>&1 && AWS_PROFILE=terraform-leviosa terraform output -raw server_ipv4_address 2>/dev/null || echo "")
DOMAIN=$(shell cd ../terraform && AWS_PROFILE=terraform-leviosa terraform workspace select default >/dev/null 2>&1 && AWS_PROFILE=terraform-leviosa terraform output -raw domain_name 2>/dev/null || echo "")

# Get SSH public key
SSH_KEY=$(shell cat ~/.ssh/leviosa.pub 2>/dev/null || cat ~/.ssh/id_rsa.pub 2>/dev/null || echo "")

# ============================================
# Setup Commands (first time only)
# ============================================

# Complete VPS setup for staging (run once on fresh server)
setup-staging:
	@echo "==========================================="
	@echo "Leviosa Staging VPS Setup"
	@echo "==========================================="
	@echo ""
	@if [ -z "$(VPS_IP)" ]; then \
		echo "Error: Could not get VPS IP from Terraform."; \
		echo "Run 'cd ../terraform && make apply-staging' first"; \
		exit 1; \
	fi
	@echo "VPS IP: $(VPS_IP)"
	@echo "Domain: $(DOMAIN)"
	@echo "SSH Key: $(SSH_KEY)"
	@echo ""
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo ""
	@echo "Running staging setup playbook..."
	ansible-playbook -i inventory/hosts.yml playbooks/site.yml \
		-e "ansible_host=$(VPS_IP)" \
		-e "app_domain=$(DOMAIN)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/leviosa_staging.yml" \
		-e "deploy_ssh_public_key=$(SSH_KEY)"

# Complete VPS setup for production (run once on fresh server)
setup-production:
	@echo "==========================================="
	@echo "Leviosa Production VPS Setup"
	@echo "==========================================="
	@echo ""
	@if [ -z "$(VPS_IP)" ]; then \
		echo "Error: Could not get VPS IP from Terraform."; \
		echo "Run 'cd ../terraform && make apply-production' first"; \
		exit 1; \
	fi
	@echo "VPS IP: $(VPS_IP)"
	@echo "Domain: $(DOMAIN)"
	@echo "SSH Key: $(SSH_KEY)"
	@echo ""
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo ""
	@echo "Running production setup playbook..."
	ansible-playbook -i inventory/hosts.yml playbooks/site.yml \
		-e "ansible_host=$(VPS_IP)" \
		-e "app_domain=$(DOMAIN)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/leviosa_production.yml" \
		-e "deploy_ssh_public_key=$(SSH_KEY)"

# ============================================
# Deployment Commands
# ============================================

# Deploy production application
deploy-production:
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
		-e "ansible_host=$(VPS_IP)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/leviosa_production.yml"

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
	@echo "Deploying staging to: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts.yml playbooks/deploy-staging.yml \
		-e "ansible_host=$(VPS_IP)" \
		-e "@group_vars/all.yml" \
		-e "@group_vars/leviosa_production.yml" \
		-e "@group_vars/leviosa_staging.yml"

# Quick restart: pull latest app image and restart only the app container
restart:
	@echo "Pulling latest prod image and restarting..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa && docker compose pull app && docker compose up -d app"

# Quick restart for staging
restart-staging:
	@echo "Pulling latest staging image and restarting..."
	@ssh -i ~/.ssh/leviosa leviosa@$(VPS_IP) "cd /opt/leviosa-staging && docker compose pull app && docker compose up -d app"

# Configure database backups (not yet configured — add S3 backup bucket to Terraform first)
backup:
	@echo "Backup not yet configured. Add a backup S3 bucket to Terraform first."
	@exit 1

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
	@which pipx > /dev/null 2>&1 || sudo apt install -y pipx
	@pipx install ansible
	@pipx inject ansible passlib
	@pipx ensurepath
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
