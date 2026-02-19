# Variables
APP_NAME=galao

# Commands
GOCMD := go
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOTEST := $(GOCMD) test
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
LINT := golangci-lint run
GOFMT := gofmt
VULN := golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: build
build: ## Run build
	@echo "Running build"
	xcodebuild \
		-project renderer/GalaoRenderer/GalaoRenderer.xcodeproj \
		-scheme GalaoRenderer \
		-configuration Release \
		CONFIGURATION_BUILD_DIR=$(PWD)/renderer

.PHONY: clean
clean: ## Run clean
	@echo "Running clean"
	xcodebuild clean \
		-project renderer/GalaoRenderer/GalaoRenderer.xcodeproj \
		-scheme GalaoRenderer

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: lint
lint: ## Run linter
	@echo "Running lint..."
	$(LINT)

.PHONY: fmt
fmt: ## Run formatter
	@echo "Formatting code..."
	$(GOFMT) -s -l -e .

.PHONY: audit
audit: ## Audit code
	@echo "Running audit..."
	$(GOMOD) tidy
	$(GOMOD) verify
	$(GOVET) ./...
	$(GORUN) $(VULN) ./...

.PHONY: help
help: ## Show available commands
	@grep -E '^[a-zA-Z_/.-]+:.*?##' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = "##"}; {printf "%-20s %s\n", $$1, $$2}'