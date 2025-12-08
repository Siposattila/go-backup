APP_NAME := go-backup
BUILD_DIR := build

DOCKER_TEST ?= small

.PHONY: all tidy watch watch_server watch_client build clean proto buffer docker_test_build docker_test docker_test_stop

all: build

tidy:
	@go fmt ./...
	@go mod tidy -v

test:
	@echo "Testing..."
	@go test ./... -v

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)

	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	@go build -o $(BUILD_DIR)/$(APP_NAME)

	@echo "Building done!"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleaning up done!"

buffer:
	@sudo sysctl -w net.core.rmem_max=7500000
	@sudo sysctl -w net.core.wmem_max=7500000

watch:
	@if ! [ -x "$(command -v air)" ]; then \
		echo "Air is not installed." >&2; \
		go install github.com/air-verse/air@latest; \
	fi \

	@echo "Watching..."

watch_server: watch
	@air -c .air_server.toml

watch_client: watch
	@air -c .air_client.toml

proto:
	@protoc --go_out=. ./protos/*.proto

docker_test_build: buffer
	@./docker/data/generate_data.sh
	@mkdir -p ./docker/qlog
	@docker compose -f ./docker/build.docker-compose.yaml build

docker_test:
	@mkdir -p ./docker/backup
	@docker swarm init
	@docker stack deploy -c ./docker/test.$(DOCKER_TEST).docker-compose.yaml go-backup-docker-test --detach=false

docker_test_stop:
	@docker stack rm go-backup-docker-test
	@docker swarm leave --force
	@sudo rm -r ./docker/backup
