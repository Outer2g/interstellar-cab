PROJECT := alblogs-to-kinesis-integration

.PHONY: run
test: ## Run unit tests
	@go run cmd/main.go

.PHONY: test
test: ## Run unit tests
	@go test -v -mod=readonly -race ./pkg/... ./cmd/...

.PHONY: build
build:  ## Build the binary and drop it under bin/
	@mkdir -p ./bin
	@go build -ldflags "-s -w" -v -mod=readonly -o bin/lambda cmd/lambda/main.go

.PHONY: dist
dist: build  ## Zip the binary and drop it under bin/
	@zip bin/alblogs-to-kinesis.zip bin/lambda

.PHONY: clean
clean: ## Remove build related files
	@rm -rf ./bin

.PHONY: lint
lint: ## Run linters
	@golangci-lint run --disable-all

.PHONY: coverage
coverage:  ## Run unit tests with coverage
	@go test -cover -covermode=count -coverprofile=cover.out ./pkg/... ./cmd/...
	@go tool cover -html=cover.out

.PHONY: integration
integration:  ## Run all tests (requires docker.up)
	@./scripts/run-tests.sh

.PHONY: docker.up
docker.up: ## Creates docker containers needed to run the integration tests
	@docker-compose -p $(PROJECT) -f docker-compose.yml -f docker-compose.integration.yml up -d localstack

.PHONY: docker.down
docker.down: ## Remove all docker containers
	@docker-compose -p $(PROJECT) -f docker-compose.yml -f docker-compose.integration.yml down --remove-orphans

.PHONY: help
help:
	@grep -E '^[a-z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help


