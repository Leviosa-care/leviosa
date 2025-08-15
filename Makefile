# vars
.PHONY: help dev build start test test-unit test-integration
.DEFAULT_GOAL := help

help:
	@./help.sh "$(MAKEFILE_LIST)"
	@echo "\n => Run 'help-front' or 'help-back' for more details on the frontend or backend part of the application"
	@echo " => You can also switch directory to run the frontend or backend part of the program. \n"

help-front:
	@$(MAKE) -C frontend help --no-print-directory

help-back:
	@$(MAKE) -C backend help --no-print-directory

dev:## Run application in developement mode
	@echo "Running the application in developement mode..."
	@$(MAKE) -C frontend $@ --no-print-directory
	@echo "Frontend application is running in dev mode"
	@$(MAKE) -C backend $@ --no-print-directory
	@echo "Backend application is running in dev mode"

build:## Build application and run it locally

start:## Start application in production mode in docker contianer using docker-compose file

test: test-unit test-integration## Run all test for application

test-unit:## Run unit test for application

test-integration:## Run integration test for application

up:
	docker compose -f local.compose.yaml up

down:
	docker compose -f local.compose.yaml down

