NAME=$(shell basename "$(PWD)")
PKGS:=$(shell go list ./... | grep -v /vendor | grep -v /fakes)

.PHONY: help
help:
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[34m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy
tidy: ## Remove unused modules
	go mod tidy

.PHONY: install
install: ## Install local dependencies
	go get ./...

.PHONY: vet
vet: ## Vet the code
	go vet ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status $(PKGS)

.PHONY: build
build: ## Build the application
	go build .

.PHONY: test
test: ## Run unit tests
	go test -bench $(PKGS)
	go tool cover -func=coverage.out

.PHONY: test-cov
test-cov: test ## Run unit tests with coverage
	go tool cover -html=coverage.out

.PHONY: all ## Run everything
all: install vet build test
