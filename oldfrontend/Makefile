# vars
IMAGE_NAME := "henga/leviosa"## The name of the docker image for the frontend part of the project

.PHONY: help run test test-unit test-integration

.DEFAULT_GOAL := help

help:
	@../help.sh "$(MAKEFILE_LIST)"

dev:## Run the application in development environment
	@echo "Adding environment variables inside .env.development"
	@if ! [ -f ".env.example" ]; then \
		echo "API_URL=http://localhost:3500\nORIGIN=http://localhost:3000" >> .env.example; \
	fi
	@pnpm run dev

build:## Run the application in production environment and run it locally
	@echo "Building application..."
	pnpm run build
	@echo "Application built."
	@echo "Running image..."
	@node build

prod:## Run the application in production environment in docker containers
	@echo "Building image..."
	docker build -t $(IMAGE_NAME):frontend . 
	@echo "Image built."
	@echo "Running image..."
	@docker run -d -p 3000:3000 $(IMAGE_NAME):frontend

test: test-unit test-integration## Run all tests for the frontend part of application

test-unit: ## Run unit tests for the frontend part of the applicatioon
	@pnpm run test:unit

test-integration:## Run the intergation tests for the frontend part of the applicatioon
	@pnpm run test:integration
