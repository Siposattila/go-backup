APP_NAME := go-backup
BUILD_DIR := build

.PHONY: all tidy watch build clean

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

	# @GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .
	# @GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .
	# @GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .
	@go build -o $(BUILD_DIR)/$(APP_NAME)

	@echo "Building done!"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleaning up done!"

linux_max_buffer_size:
	@sysctl -w net.core.rmem_max=7500000
	@sysctl -w net.core.wmem_max=7500000

watch:
	@if ! [ -x "$(command -v air)" ]; then \
		echo "Air is not installed." >&2; \
		go install github.com/air-verse/air@latest; \
	fi \

	@echo "Watching..."
	@air
