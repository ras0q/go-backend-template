SHELL    := /bin/bash
APP_NAME := $(shell basename $(PWD))
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GO_TEST_FLAGS := -v -cover -race -shuffle=on

.PHONY: help
help: ## Display this help
	@echo "Commands:"
	@grep -E '^[a-z-]+:.*## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: dev
dev: ## Run the development server with live-reload
	docker compose up --build

.PHONY: mod
mod: ## Download the dependencies
	go mod download

build: $(GO_FILES) ## Compile the binary
	go build -o $(APP_NAME)

.PHONY: test
test: test-unit test-integration ## Run all the tests

.PHONY: test-unit
test-unit: ## Run the unit tests
	go test $(GO_TEST_FLAGS) . ./internal/...

.PHONY: test-integration
test-integration: ## Run the integration tests
	go test $(GO_TEST_FLAGS) ./integration/...

.PHONY: lint
lint: ## Run the linter
	golangci-lint run --timeout=5m --fix ./...
