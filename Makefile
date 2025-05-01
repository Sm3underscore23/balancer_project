include .env

LOCAL_BIN := "$(CURDIR)"/bin
LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := "host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.2

migration-create:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create create_tables sql

migration-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v

migration-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

migration-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

build-balancer:
	docker buildx build --platform linux/amd64 -t balancer:v0.1 -f ./balancer.Dockerfile .

build-backend-test:
	docker buildx build --platform linux/amd64 -t backend-test:v0.1 -f ./backend-test.Dockerfile .

db-up:
	docker compose -f db.docker-compose.yaml up

db-down:
	docker compose -f db.docker-compose.yaml down -v

fullsetup-up:
	docker compose -f fullsetup.docker-compose.yaml up

fullsetup-down:
	docker compose -f fullsetup.docker-compose.yaml down -v