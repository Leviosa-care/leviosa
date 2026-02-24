# Leviosa Ansible Makefile
# Usage: make <target>

.PHONY: help setup deploy backup check-clean lint \
        ssh-info status logs health-check \
        install-deps vars-test

# Default target
help:
	@echo "Leviosa Ansible Commands"
	@echo "============================"
	@echo ""
	@echo "Setup:"
	@echo "  make setup           - Complete VPS setup (first time only)"
	@echo "  make deploy          - Deploy/update application"
	@echo "  make backup          - Configure database backups"
	@echo ""
	@echo "Verification:"
	@echo "  make check-clean     - Check for syntax errors"
	@echo "  make lint            - Run ansible-lint"
	@echo "  make vars-test       - Display effective variables"
	@echo ""
	@echo "Server Operations:"
	@echo "  make ssh-info        - Display SSH connection details"
	@echo "  make status          - Check service status"
	@echo "  make logs            - View application logs"
	@echo "  make health-check    - Run health check on server"
	@echo ""
	@echo "Dependencies:"
	@echo "  make install-deps    - Install Ansible dependencies"

# Get VPS IP from inventory (first webserver)
VPS_IP=$(shell grep -A1 '\[webserver\]' inventory/hosts | tail -1 | awk '{print $$1}')
DOMAIN=$(shell grep -E '^domain' /etc/hosts 2>/dev/null || echo 'leviosa.care')

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
	ansible-playbook -i inventory/hosts playbook.yml --ask-vault-pass

# Deploy production application
deploy:
	@echo "==========================================="
	@echo "Leviosa Deployment"
	@echo "==========================================="
	@echo ""
	@echo "Deploying to: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts deploy.yml --ask-vault-pass

# Configure database backups
backup:
	@echo "==========================================="
	@echo "Database Backup Configuration"
	@echo "==========================================="
	@echo ""
	@echo "VPS IP: $(VPS_IP)"
	@echo ""
	ansible-playbook -i inventory/hosts backup.yml --ask-vault-pass \
		-e "backup_s3_bucket=$${BACKUP_S3_BUCKET:-staging-leviosa-backups}" \
		-e "backup_s3_region=$${BACKUP_S3_REGION:-eu-central-1}"

# ============================================
# Verification Commands
# ============================================

# Check playbooks for syntax errors
check-clean:
	@echo "Checking playbooks for syntax errors..."
	@ansible-playbook playbook.yml --syntax-check
	@ansible-playbook deploy.yml --syntax-check
	@ansible-playbook backup.yml --syntax-check
	@echo "All playbooks are clean!"

# Run ansible-lint (requires installation)
lint:
	@echo "Running ansible-lint..."
	@which ansible-lint > /dev/null || (echo "Install with: pip install ansible-lint" && exit 1)
	@ansible-lint playbook.yml deploy.yml backup.yml

# Display effective variables
vars-test:
	@echo "Displaying effective variables..."
	@ansible-playbook -i inventory/hosts playbook.yml \
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
	@echo "Checking service status..."
	@ssh -i ~/.ssh/leviosa root@$(VPS_IP) "docker compose ps"

# View application logs
logs:
	@echo "Streaming logs (Ctrl+C to exit)..."
	@ssh -i ~/.ssh/leviosa root@$(VPS_IP) "docker compose logs -f"

# Run health check
health-check:
	@echo "Running health check..."
	@ssh -i ~/.ssh/leviosa root@$(VPS_IP) "curl -f http://localhost:3000/health || exit 1"

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
