include .env

LOCAL_PROTO_SRC_DIR := "$(CURDIR)"/api/$(PROTO_FILE_VERSION)
LOCAL_PROTO_DST_DIR  = "$(CURDIR)"/pkg/$(PROTO_FILE_VERSION)

LOCAL_BIN := "$(CURDIR)"/bin
LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := "host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"

proto-generator:
	protoc -I="./api/$(PROTO_FILE_VERSION)" --go_out=paths=source_relative:"./pkg/$(PROTO_FILE_VERSION)" ./api/$(PROTO_FILE_VERSION)/user.proto
	protoc -I="./api/$(PROTO_FILE_VERSION)" --go-grpc_out=paths=source_relative:"./pkg/$(PROTO_FILE_VERSION)" ./api/$(PROTO_FILE_VERSION)/user.proto

install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.2

local-migration-create:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create create_tables sql

local-migration-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

app-build-and-run:
	sudo docker buildx build --platform linux/amd64 -t balancer:v0.1 .
	sudo docker run -d --rm -p 8080:8080 auth-service:v0.1

fast-start:
	sudo docker compose up -d
	sudo docker compose ps
	until sudo docker compose ps | grep "healthy"; do sleep 1; done
	make local-migration-up
	go run cmd/main.go

force-stop:
	sudo docker stop db
	sudo docker rm db
	sudo docker ps -a
