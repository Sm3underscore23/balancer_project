include .env

LOCAL_BIN := "$(CURDIR)"/bin
LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := "host=$(PG_HOST) port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.2

check-goose:
	@if [ ! -f $(LOCAL_BIN)/goose ]; then \
		echo "goose not found, installing..."; \
		$(MAKE) install-goose; \
	fi

migration-create:
	$(MAKE) check-goose
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create create_tables sql

migration-status:
	$(MAKE) check-goose
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v

migration-up:
	$(MAKE) check-goose
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

migration-down:
	$(MAKE) check-goose
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

build-balancer:
	docker buildx build --platform linux/amd64 -t balancer:v0.1 -f ./balancer.Dockerfile .

build-backend-test:
	docker buildx build --platform linux/amd64 -t backend-test:v0.1 -f ./backend-test.Dockerfile .

local-up:
	docker compose -f local.docker-compose.yaml up

local-down:
	docker compose -f local.docker-compose.yaml down -v

local-run:
	go run cmd/main.go -config-path config/local/config.yaml -local

fullsetup-up:
	docker compose -f fullsetup.docker-compose.yaml up

fullsetup-down:
	docker compose -f fullsetup.docker-compose.yaml down -v