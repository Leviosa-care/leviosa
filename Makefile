# Leviosa VPS Deployment Makefile
# Usage: make <target>

# AWS profile for Terraform operations
export AWS_PROFILE := terraform-leviosa

# Docker image configuration
# These must match the images defined in infra/ansible/group_vars/leviosa_*.yml
DOCKER_USER ?= henga
DOCKER_IMAGE_FRONTEND_PROD ?= $(DOCKER_USER)/leviosa:frontend-latest
DOCKER_IMAGE_BACKEND_PROD ?= $(DOCKER_USER)/leviosa:backend-latest
DOCKER_IMAGE_FRONTEND_STAGING ?= $(DOCKER_USER)/leviosa:frontend-staging
DOCKER_IMAGE_BACKEND_STAGING ?= $(DOCKER_USER)/leviosa:backend-staging

# VPS connection - resolved from Terraform output
VPS_IP ?= $(shell cd infra/terraform && AWS_PROFILE=terraform-leviosa terraform output -raw server_ipv4_address 2>/dev/null || echo "")
VPS_USER ?= root
VPS_SSH = ssh -i ~/.ssh/leviosa $(VPS_USER)@$(VPS_IP)

# App directories on the VPS
PROD_DIR = /opt/leviosa
STAGING_DIR = /opt/leviosa-staging

# Build args
SESSION_COOKIE_NAME ?= leviosa_session
CLIENT_IP_HEADER ?= x-client-ip
BACKEND_PORT ?= 3500

.PHONY: help \
		local-dev local-up local-down \
		image-build-front image-build-back image-build \
		image-build-staging-front image-build-staging-back image-build-staging \
		image-push-front image-push-back image-push \
		image-push-staging-front image-push-staging-back image-push-staging \
		image-release-front image-release-back image-release \
		image-release-staging-front image-release-staging-back image-release-staging \
		deploy deploy-staging \
		prod-start prod-stop prod-restart prod-logs prod-status prod-ssh \
		prod-update-front prod-update-back \
		staging-start staging-stop staging-restart staging-logs staging-status \
		staging-update-front staging-update-back \
		infra-init infra-plan infra-apply infra-output infra-destroy \
		ansible-setup ansible-deploy ansible-deploy-staging \
		ansible-restart ansible-restart-staging \
		ansible-logs ansible-logs-staging ansible-status ansible-status-staging

# Default target
help:
	@echo "Leviosa VPS Deployment Commands"
	@echo "================================="
	@echo ""
	@echo "Local Development:"
	@echo "  make local-dev        - Start frontend and backend in dev mode"
	@echo "  make local-up         - Start local docker compose"
	@echo "  make local-down       - Stop local docker compose"
	@echo ""
	@echo "Docker Images (build locally):"
	@echo "  make image-build      - Build frontend and backend images (prod)"
	@echo "  make image-build-front       - Build production frontend only"
	@echo "  make image-build-back        - Build production backend only"
	@echo "  make image-build-staging     - Build frontend and backend images (staging)"
	@echo "  make image-build-staging-front    - Build staging frontend only"
	@echo "  make image-build-staging-back     - Build staging backend only"
	@echo "  make image-push       - Push images to Docker Hub (prod)"
	@echo "  make image-push-staging    - Push staging images to Docker Hub"
	@echo "  make image-release    - Build and push images in one step (prod)"
	@echo "  make image-release-staging - Build and push staging images"
	@echo ""
	@echo "Quick Deploy (SSH-based):"
	@echo "  make deploy           - Pull and restart production on VPS"
	@echo "  make deploy-staging   - Pull and restart staging on VPS"
	@echo ""
	@echo "Production Operations:"
	@echo "  make prod-start       - Start production stack on VPS"
	@echo "  make prod-stop        - Stop production stack on VPS"
	@echo "  make prod-restart     - Restart production stack on VPS"
	@echo "  make prod-logs        - View production logs"
	@echo "  make prod-status      - Check production container status"
	@echo "  make prod-ssh         - SSH into the VPS"
	@echo "  make prod-update-front - Pull latest frontend image and restart container"
	@echo "  make prod-update-back  - Pull latest backend image and restart container"
	@echo ""
	@echo "Staging Operations:"
	@echo "  make staging-start    - Start staging stack on VPS"
	@echo "  make staging-stop     - Stop staging stack on VPS"
	@echo "  make staging-restart  - Restart staging stack on VPS"
	@echo "  make staging-logs     - View staging logs"
	@echo "  make staging-status   - Check staging container status"
	@echo "  make staging-update-front - Pull latest staging frontend image and restart container"
	@echo "  make staging-update-back  - Pull latest staging backend image and restart container"
	@echo ""
	@echo "Infrastructure (Terraform):"
	@echo "  make infra-init       - Initialize Terraform"
	@echo "  make infra-plan       - Show Terraform plan"
	@echo "  make infra-apply      - Apply Terraform changes"
	@echo "  make infra-output     - Show Terraform outputs"
	@echo "  make infra-destroy    - Destroy all infrastructure"
	@echo ""
	@echo "Configuration (Ansible):"
	@echo "  make ansible-setup    - Complete VPS setup (first time only)"
	@echo "  make ansible-deploy   - Deploy/update production via Ansible"
	@echo "  make ansible-deploy-staging - Deploy staging via Ansible"
	@echo "  make ansible-restart  - Quick restart production"
	@echo "  make ansible-restart-staging - Quick restart staging"
	@echo ""
	@echo "VPS: $(VPS_USER)@$(VPS_IP)"
	@echo "Prod frontend: $(DOCKER_IMAGE_FRONTEND_PROD)"
	@echo "Prod backend:  $(DOCKER_IMAGE_BACKEND_PROD)"
	@echo "Staging frontend: $(DOCKER_IMAGE_FRONTEND_STAGING)"
	@echo "Staging backend:  $(DOCKER_IMAGE_BACKEND_STAGING)"

# ===========================================
# Local Development
# ===========================================

local-dev:
	@echo "Starting frontend and backend in dev mode..."
	@cd frontend && pnpm run dev &
	@cd backend && docker compose up

local-up:
	@echo "Starting local docker compose..."
	docker compose -f local.compose.yaml up

local-down:
	@echo "Stopping local docker compose..."
	docker compose -f local.compose.yaml down

# ===========================================
# Docker Image Commands (run on local machine)
# ===========================================

image-build-front:
	@echo "Building production frontend image..."
	docker build \
		-t $(DOCKER_IMAGE_FRONTEND_PROD) \
		-f frontend/Dockerfile \
		--build-arg SESSION_COOKIE_NAME=$(SESSION_COOKIE_NAME) \
		--build-arg CLIENT_IP_HEADER=$(CLIENT_IP_HEADER) \
		--build-arg BACKEND_PORT=$(BACKEND_PORT) \
		./frontend
	@echo "Frontend image built: $(DOCKER_IMAGE_FRONTEND_PROD)"

image-build-back:
	@echo "Building production backend image..."
	docker build \
		-t $(DOCKER_IMAGE_BACKEND_PROD) \
		-f backend/Dockerfile \
		./backend
	@echo "Backend image built: $(DOCKER_IMAGE_BACKEND_PROD)"

image-build: image-build-front image-build-back

image-push-front:
	@echo "Pushing frontend image to Docker Hub..."
	docker push $(DOCKER_IMAGE_FRONTEND_PROD)
	@echo "Frontend image pushed: $(DOCKER_IMAGE_FRONTEND_PROD)"

image-push-back:
	@echo "Pushing backend image to Docker Hub..."
	docker push $(DOCKER_IMAGE_BACKEND_PROD)
	@echo "Backend image pushed: $(DOCKER_IMAGE_BACKEND_PROD)"

image-push: image-push-front image-push-back

image-release: image-build image-push
	@echo "Production images released."
	@echo "Run 'make ansible-deploy' or 'make ansible-restart' to deploy."

# Production release commands (build + push)
image-release-front: image-build-front image-push-front
	@echo "Production frontend image released: $(DOCKER_IMAGE_FRONTEND_PROD)"

image-release-back: image-build-back image-push-back
	@echo "Production backend image released: $(DOCKER_IMAGE_BACKEND_PROD)"

# Staging build commands
image-build-staging-front:
	@echo "Building staging frontend image..."
	docker build \
		-t $(DOCKER_IMAGE_FRONTEND_STAGING) \
		-f frontend/Dockerfile \
		--build-arg SESSION_COOKIE_NAME=$(SESSION_COOKIE_NAME) \
		--build-arg CLIENT_IP_HEADER=$(CLIENT_IP_HEADER) \
		--build-arg BACKEND_PORT=$(BACKEND_PORT) \
		./frontend
	@echo "Staging frontend image built: $(DOCKER_IMAGE_FRONTEND_STAGING)"

image-build-staging-back:
	@echo "Building staging backend image..."
	docker build \
		-t $(DOCKER_IMAGE_BACKEND_STAGING) \
		-f backend/Dockerfile \
		./backend
	@echo "Staging backend image built: $(DOCKER_IMAGE_BACKEND_STAGING)"

image-build-staging: image-build-staging-front image-build-staging-back

image-push-staging-front:
	@echo "Pushing staging frontend image to Docker Hub..."
	docker push $(DOCKER_IMAGE_FRONTEND_STAGING)
	@echo "Staging frontend image pushed: $(DOCKER_IMAGE_FRONTEND_STAGING)"

image-push-staging-back:
	@echo "Pushing staging backend image to Docker Hub..."
	docker push $(DOCKER_IMAGE_BACKEND_STAGING)
	@echo "Staging backend image pushed: $(DOCKER_IMAGE_BACKEND_STAGING)"

image-push-staging: image-push-staging-front image-push-staging-back

# Staging release commands (build + push)
image-release-staging-front: image-build-staging-front image-push-staging-front
	@echo "Staging frontend image released: $(DOCKER_IMAGE_FRONTEND_STAGING)"

image-release-staging-back: image-build-staging-back image-push-staging-back
	@echo "Staging backend image released: $(DOCKER_IMAGE_BACKEND_STAGING)"

image-release-staging: image-release-staging-front image-release-staging-back
	@echo "Staging images released."
	@echo "Run 'make ansible-deploy-staging' or 'make ansible-restart-staging' to deploy."

# ===========================================
# Quick Deploy (SSH-based)
# ===========================================

deploy:
	@echo "Deploying production (quick)..."
	@$(VPS_SSH) "cd $(PROD_DIR) && docker compose pull && docker compose up -d"
	@echo "Production deployed."

deploy-staging:
	@echo "Deploying staging (quick)..."
	@$(VPS_SSH) "cd $(STAGING_DIR) && docker compose pull && docker compose up -d"
	@echo "Staging deployed."

# ===========================================
# Production Commands (VPS)
# ===========================================

prod-start:
	@echo "Starting production on VPS..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose up -d"

prod-stop:
	@echo "Stopping production on VPS..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose down"

prod-restart:
	@echo "Restarting production on VPS..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose down && docker compose up -d"

prod-logs:
	@echo "Streaming production logs (Ctrl+C to exit)..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose logs -f"

prod-status:
	@echo "Checking production container status..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose ps"

prod-ssh:
	$(VPS_SSH)

prod-update-front:
	@echo "Pulling latest frontend image and restarting container..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose pull frontend && docker compose up -d --no-deps frontend"

prod-update-back:
	@echo "Pulling latest backend image and restarting container..."
	$(VPS_SSH) "cd $(PROD_DIR) && docker compose pull backend && docker compose up -d --no-deps backend"

# ===========================================
# Staging Commands (VPS)
# ===========================================

staging-start:
	@echo "Starting staging on VPS..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose up -d"

staging-stop:
	@echo "Stopping staging on VPS..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose down"

staging-restart:
	@echo "Restarting staging on VPS..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose down && docker compose up -d"

staging-logs:
	@echo "Streaming staging logs (Ctrl+C to exit)..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose logs -f"

staging-status:
	@echo "Checking staging container status..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose ps"

staging-update-front:
	@echo "Pulling latest staging frontend image and restarting container..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose pull frontend && docker compose up -d --no-deps frontend"

staging-update-back:
	@echo "Pulling latest staging backend image and restarting container..."
	$(VPS_SSH) "cd $(STAGING_DIR) && docker compose pull backend && docker compose up -d --no-deps backend"

# ===========================================
# Infrastructure Commands (Terraform)
# ===========================================

infra-init:
	@echo "Initializing Terraform..."
	cd infra/terraform && terraform init

infra-plan:
	@echo "Showing Terraform plan..."
	cd infra/terraform && terraform plan

infra-apply:
	@echo "Applying Terraform changes..."
	cd infra/terraform && terraform apply

infra-output:
	@echo "Terraform outputs:"
	cd infra/terraform && terraform output

infra-destroy:
	@echo "Destroying all Terraform resources..."
	cd infra/terraform && terraform destroy

# ===========================================
# Configuration Commands (Ansible)
# ===========================================

ansible-setup:
	@echo "Running complete VPS setup..."
	cd infra/ansible && $(MAKE) setup

ansible-deploy:
	@echo "Deploying production via Ansible..."
	cd infra/ansible && $(MAKE) deploy-production

ansible-deploy-staging:
	@echo "Deploying staging via Ansible..."
	cd infra/ansible && $(MAKE) deploy-staging

ansible-restart:
	@echo "Quick restart production..."
	cd infra/ansible && $(MAKE) restart

ansible-restart-staging:
	@echo "Quick restart staging..."
	cd infra/ansible && $(MAKE) restart-staging

ansible-logs:
	@echo "Streaming production logs..."
	cd infra/ansible && $(MAKE) logs

ansible-logs-staging:
	@echo "Streaming staging logs..."
	cd infra/ansible && $(MAKE) staging-logs

ansible-status:
	@echo "Checking production status..."
	cd infra/ansible && $(MAKE) status

ansible-status-staging:
	@echo "Checking staging status..."
	cd infra/ansible && $(MAKE) staging-status
