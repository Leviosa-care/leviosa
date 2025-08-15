# Environment file (create a .env file in your project root if you haven't) ENV_FILE ?= .env

# Redis configuration
REDIS_IMAGE = redis:alpine
REDIS_CONTAINER = myapp-redis
REDIS_PORT = 6379

POSTGRES_IMAGE = postgres:17.5-alpine3.21
POSTGRES_CONTAINER = myapp-postgres
POSTGRES_PORT = 5432

RABBIT_IMAGE = rabbitmq:3-management-alpine
RABBIT_CONTAINER = myapp-rabbit
RABBIT_PORT1=5672
RABBIT_PORT2=15672
RABBIT_USER=guest
RABBIT_PASSWORD=guest

.DEFAULT_GOAL:= run

.PHONY: run
run: deps
	@# APP_ENV=production air -c .air.toml
	APP_ENV=production \
	BUCKETNAME=leviosa-assets \
	AWS_REGION=eu-central-1 \
	AWS_ACCESS_KEY_ID=AKIA6D6JBL4GFGD47U4N \
	AWS_SECRET_ACCESS_KEY=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 \
	GMAIL_EMAIL=leviosacare@gmail.com \
	GMAIL_PASSWORD="eurh dnow ulys zzty" \
	STRIPE_SECRET_KEY="sk_test_51OprHjHwHXlEm0ohOgDJOaArlojFAGnNz3xULv6hgG4OCKwHw6dynWchGxjzQDZC0v8s1lXBvVEAxlbjuWq3Ho9f00ex88AzJ3" \
	REDIS_ADDR=127.0.0.1:6379 \
	REDIS_DB=0 \
	REDIS_PASSWORD=secret \
	LOGGING_SALT=20eaddaabb68c5db7bee3232552138b3 \
	VAULT_ADDR=https://leviosa-cluster-public-vault-9da92754.27b5d5f2.z1.hashicorp.cloud:8200 \
	VAULT_NAMESPACE=admin/leviosa \
	VAULT_ROLE_ID=49f35e23-4c96-a413-e4e6-dd31271d7299 \
	VAULT_SECRET_ID=645c80cb-d8b7-6c5c-db1a-e71210be6e1b \
	RABBITMQ_HOST=localhost \
	RABBITMQ_PORT=5672 \
	RABBITMQ_USER=guest \
	RABBITMQ_PASSWORD=guest \
	POSTGRES_PASSWORD=password \
	POSTGRES_USER=user \
	POSTGRES_DB=leviosa \
	POSTGRES_PORT=5432 \
	POSTGRES_HOST=localhost \
	air -c .air.toml;

.PHONY: deps
deps: start-redis start-postgres start-rabbit
	echo "Redis, Postgres and RabbitMQ started"


.PHONY: start-redis
start-redis:
	@if [ -z "$$(docker ps -q -f name=$(REDIS_CONTAINER))" ]; then \
		if [ -z "$$(docker ps -aq -f name=$(REDIS_CONTAINER))" ]; then \
			docker run -d \
				--name $(REDIS_CONTAINER) \
				-p "127.0.0.1:$(REDIS_PORT)":$(REDIS_PORT) \
				--env-file $(ENV_FILE) \
				$(REDIS_IMAGE) redis-server --requirepass secret; \
			echo "Redis container created and started"; \
		else \
			docker start $(REDIS_CONTAINER); \
			echo "Redis container started"; \
		fi \
	else \
		echo "Redis container already running"; \
	fi

.PHONY: stop-redis
stop-redis:
	@docker stop $(REDIS_CONTAINER) 2>/dev/null || true
	echo "Redis container stopped"

.PHONY: rm-redis
rm-redis:
	@docker stop $(REDIS_CONTAINER) 2>/dev/null || true
	@docker rm $(REDIS_CONTAINER) 2>/dev/null || true
	echo "Redis container stopped"

.PHONY: start-postgres
start-postgres:
	@if [ -z "$$(docker ps -q -f name=$(POSTGRES_CONTAINER))" ]; then \
		if [ -z "$$(docker ps -aq -f name=$(POSTGRES_CONTAINER))" ]; then \
			docker run -d \
				--name $(POSTGRES_CONTAINER) \
				-p $(POSTGRES_PORT):$(POSTGRES_PORT) \
				-v pg_data:/var/lib/postgresql/data \
				-e POSTGRES_PASSWORD=password \
				-e POSTGRES_USER=user \
				-e POSTGRES_DB=development_leviosa \
				$(POSTGRES_IMAGE); \
			echo "Postgres container created and started"; \
		else \
			docker start $(POSTGRES_CONTAINER); \
			echo "Postgres container started"; \
		fi \
	else \
		echo "Postgres container already running"; \
	fi

.PHONY: stop-postgres
stop-postgres:
	@docker stop $(POSTGRES_CONTAINER) 2>/dev/null || true
	echo "Postgres container stopped"

.PHONY: rm-postgres
rm-postgres:
	@docker stop $(POSTGRES_CONTAINER) 2>/dev/null || true
	@docker rm $(POSTGRES_CONTAINER) 2>/dev/null || true
	echo "Postgres container stopped and removed"

.PHONY: start-rabbit
start-rabbit:
	@if [ -z "$$(docker ps -q -f name=$(RABBIT_CONTAINER))" ]; then \
		if [ -z "$$(docker ps -aq -f name=$(RABBIT_CONTAINER))" ]; then \
			docker run -d \
				--name $(RABBIT_CONTAINER) \
				-p $(RABBIT_PORT1):$(RABBIT_PORT1) \
				-p $(RABBIT_PORT2):$(RABBIT_PORT2) \
				-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
				-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
				$(RABBIT_IMAGE); \
			echo "RabbitMQ container created and started"; \
		else \
			docker start $(RABBIT_CONTAINER); \
			echo "RabbitMQ container started"; \
		fi \
	else \
		echo "RabbitMQ container already running"; \
	fi

.PHONY: stop-rabbit
stop-rabbit:
	@docker stop $(RABBIT_CONTAINER) 2>/dev/null || true
	echo "rabbit container stopped"

.PHONY: rm-rabbit
rm-rabbit:
	@docker stop $(RABBIT_CONTAINER) 2>/dev/null || true
	@docker rm $(RABBIT_CONTAINER) 2>/dev/null || true
	echo "rabbit container removed"

.PHONY: stop
stop: stop-redis stop-postgres stop-rabbit

.PHONY: clean
clean: rm-redis rm-postgres rm-rabbit

.PHONY: up
up:
	APP_ENV=production docker compose -f compose.yaml up

.PHONY: down
down:
	APP_ENV=production docker compose -f compose.yaml down

.PHONY: docker
docker:
	docker run \
		-p 3500:3500 \
		-e APP_ENV=production \
		-e BUCKETNAME=leviosa-assets \
		-e AWS_REGION=eu-central-1 \
		-e AWS_ACCESS_KEY_ID=AKIA6D6JBL4GFGD47U4N \
		-e AWS_SECRET_ACCESS_KEY=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 \
		-e GMAIL_EMAIL=leviosacare@gmail.com \
		-e GMAIL_PASSWORD="eurh dnow ulys zzty" \
		-e STRIPE_SECRET_KEY="sk_test_51OprHjHwHXlEm0ohOgDJOaArlojFAGnNz3xULv6hgG4OCKwHw6dynWchGxjzQDZC0v8s1lXBvVEAxlbjuWq3Ho9f00ex88AzJ3" \
		-e REDIS_ADDR=127.0.0.1:6379 \
		-e REDIS_DB=0 \
		-e REDIS_PASSWORD=secret \
		-e LOGGING_SALT=20eaddaabb68c5db7bee3232552138b3 \
		-e VAULT_ADDR=https://leviosa-cluster-public-vault-9da92754.27b5d5f2.z1.hashicorp.cloud:8200 \
		-e VAULT_NAMESPACE=admin/leviosa \
		-e VAULT_ROLE_ID=49f35e23-4c96-a413-e4e6-dd31271d7299 \
		-e VAULT_SECRET_ID=645c80cb-d8b7-6c5c-db1a-e71210be6e1b \
		-e RABBITMQ_HOST=localhost \
		-e RABBITMQ_PORT=5672 \
		-e RABBITMQ_USER=guest \
		-e RABBITMQ_PASSWORD=guest \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_USER=user \
		-e POSTGRES_DB=leviosa \
		-e POSTGRES_PORT=5432 \
		-e POSTGRES_HOST=localhost \
		leviosa-backend:1;
