include default.env

# GLOBAL ENVIRONMENT VARIABLES #################################################
APP_NAME = posts-api

# GENERAL PURPOSE TARGETS ######################################################
.PHONY: help
help: ## Display available commands.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Style and cleanup.

.PHONY: lint
lint: ## Run all code formatters, linters and analysers. Must always be run before GIT commit or push commands!
	# Format source code.
	goimports -w .
	# Report source code irregularities.
	go vet ./...
	# Report cyclomatic complexity of functions over 15.
	gocyclo -over 15 -ignore "_test|testdata/" .
	# Report source code errors, bugs and smells.
	golangci-lint run --config scripts/linter/.golangci.yaml
	# Prune unused packages and dependencies from the go.mod file.
	go mod tidy -v
	# Remove artefact files left behind by other operations.
	go clean ./...

# Test targets.

TIMEOUT = 20
ALL_PACKAGES = $(shell go list ./... | grep -v -E "vendor")
ALL_PACKAGES_EXCL_INTEGRATION=$(shell go list ./... | grep -v -E "vendor" | grep -v -E "integration")
INTEGRATION_PACKAGES=$(shell go list ./... | grep -v -E "vendor" | grep "integration")
.PHONY: test
test: ## Run tests on a complied project
	go test -timeout $(TIMEOUT)s -v $(ALL_PACKAGES)

.PHONY: test-unit
test-unit: format ## Run unit tests on a complied project (excluding integration tests)
	go test -timeout $(TIMEOUT)s -v $(ALL_PACKAGES_EXCL_INTEGRATION)

.PHONY: test-integration
test-integration: format ## Run integration tests on a complied project
	go test -timeout $(TIMEOUT)s -v $(INTEGRATION_PACKAGES)

.PHONY: test-with-coverage
test-with-coverage: ## Run all unit tests with coverage
	go test -coverprofile=coverage_unit.txt -timeout $(TIMEOUT)s -covermode=atomic $(ALL_PACKAGES)

# Build targets.

.PHONY: build
build:
	./scripts/build.sh

.PHONY: build-darwin
build-darwin:
	go mod tidy
	GOOS=darwin GOARCH=arm64 go build -a -o bin/gd-api cmd/main.go

.PHONY: run
run: ## Run application with default arugments
	go run cmd/main.go --config default.env

.PHONY: clean
clean: ## Clean repository
	@test ! -e bin || rm -r bin

# Docker targets.

VERSION = local
APP_BIN = bin/gd-api
REPOSITORY = gd-api
.PHONY: build-docker-local
build-docker-local: ## Build local docker image
	docker build \
		--tag $(REPOSITORY):$(VERSION) \
		--build-arg	APP_BIN=$(APP_BIN) \
		-f deploy/local/Dockerfile .

.PHONY: build-docker
build-docker: ## Build docker image
	docker build \
		--tag $(REPOSITORY):$(VERSION) \
		--build-arg	APP_BIN=$(APP_BIN) \
		-f deploy/Dockerfile .

.PHONY: docker-compose-up-infra-local
docker-compose-up-infra-local: ## Run infra services locally using docker compose
	docker-compose -f deploy/local/docker-compose.yaml up -d db-setup db

.PHONY: docker-compose-up-local
docker-compose-up-local: ## Run all services locally using docker compose
	docker-compose -f deploy/local/docker-compose.yaml up -d --build

.PHONY: docker-compose-down-local
docker-compose-down-local: ## Run infra services locally using docker compose
	docker-compose -f deploy/local/docker-compose.yaml down